package data

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/GerardoHP/inventory-service/models"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	credentialsFile string = "service-account-file.json"
	collectionName  string = "products"
)

var (
	collection *firestore.CollectionRef
	products   = struct {
		sync.RWMutex
		m map[int]models.Product
	}{m: make(map[int]models.Product)}
)

type FirestoreRepository struct {
}

var _ ProductRepository = (*FirestoreRepository)(nil)

func NewFirestoreRepository() *FirestoreRepository {
	return &FirestoreRepository{}
}

func (r FirestoreRepository) GetProducts() ([]models.Product, error) {
	products.Lock()
	defer products.Unlock()
	prods := make([]models.Product, 0)
	for _, value := range products.m {
		prods = append(prods, value)
	}

	return prods, nil
}

func (r FirestoreRepository) AddProduct(p models.Product) error {
	m := p.ToMap()
	_, _, err := collection.Add(context.Background(), m)
	if err != nil {
		fmt.Println(err)
	} else {
		products.Lock()
		defer products.Unlock()
		products.m[p.ProductID] = p
	}

	return nil
}

func (r FirestoreRepository) GetNextID() (id int, err error) {
	products.Lock()
	defer products.Unlock()
	productIds := []int{}
	for i := range products.m {
		productIds = append(productIds, i)
	}

	sort.Ints(productIds)
	id = productIds[len(productIds)-1] + 1
	return
}

func (r FirestoreRepository) GetProductById(id int) (*models.Product, error) {
	products.Lock()
	defer products.Unlock()
	if product, ok := products.m[id]; ok {
		return &product, nil
	}

	return nil, nil
}

func (r FirestoreRepository) UpdateProduct(p models.Product) error {
	products.Lock()
	products.m[p.ProductID] = p
	products.Unlock()
	_, refId := getProduct(p.ProductID)
	docRef := collection.Doc(refId)
	_, err := docRef.Set(context.Background(), p.ToMap())
	if err != nil {
		log.Printf("an error has ocurred: %s \n", err)
	}

	return nil
}

func (r FirestoreRepository) DeleteProduct(id int) error {
	products.Lock()
	delete(products.m, id)
	products.Unlock()
	_, refId := getProduct(id)
	_, err := collection.Doc(refId).Delete(context.Background())
	if err != nil {
		log.Printf("an error has ocurred: %s \n", err)
	}

	return nil
}

func (FirestoreRepository) GetTopProducts(top int) ([]models.Product, error) {
	return nil, nil
}

func (FirestoreRepository) SearchForProductData(models.ProductReportFilter) ([]models.Product, error) {
	return nil, nil
}

func init() {
	opt := option.WithCredentialsFile(credentialsFile)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	collection = client.Collection(collectionName)
	iter := collection.Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatalf("Failed to iterate: %v \n", err)
		}

		product := models.ProductFromMap(doc.Data())
		products.m[product.ProductID] = *product
	}
}

func getProduct(productId int) (*models.Product, string) {
	iter := collection.Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatalf("Failed to iterate: %v \n", err)
		}

		d := doc.Data()["productId"]
		if number, ok := d.(int64); ok && productId == int(number) {
			return models.ProductFromMap(doc.Data()), doc.Ref.ID
		}
	}

	return nil, ""
}
