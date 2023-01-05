package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/GerardoHP/inventory-service/models"
	_ "github.com/go-sql-driver/mysql"
)

type JsonRepository struct {
}

var _ ProductRepository = (*JsonRepository)(nil)

func NewJsonRepository() *JsonRepository {
	return &JsonRepository{}
}

var (
	productMap = struct {
		sync.RWMutex
		m map[int]models.Product
	}{m: make(map[int]models.Product)}
)

func (r JsonRepository) GetProducts() ([]models.Product, error) {
	productMap.Lock()
	products := make([]models.Product, 0, len(productMap.m))
	for _, value := range productMap.m {
		products = append(products, value)
	}

	productMap.Unlock()
	return products, nil
}

func (r JsonRepository) AddProduct(p models.Product) error {
	productMap.Lock()
	defer productMap.Unlock()
	productMap.m[p.ProductID] = p

	return nil
}

func (r JsonRepository) GetNextID() (id int, err error) {
	productMap.Lock()
	defer productMap.Unlock()
	productIds := []int{}
	for i := range productMap.m {
		productIds = append(productIds, i)
	}

	sort.Ints(productIds)
	id = productIds[len(productIds)-1] + 1
	return
}

func (r JsonRepository) GetProductById(id int) (*models.Product, error) {
	productMap.Lock()
	defer productMap.Unlock()
	if product, ok := productMap.m[id]; ok {
		return &product, nil
	}

	return nil, nil
}

func (r JsonRepository) UpdateProduct(p models.Product) error {
	productMap.Lock()
	defer productMap.Unlock()
	productMap.m[p.ProductID] = p
	return nil
}

func (r JsonRepository) DeleteProduct(id int) error {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, id)

	return nil
}

func (JsonRepository) GetTopProducts(top int) ([]models.Product, error) {
	return nil, nil
}

func init() {
	fmt.Println("loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d products loaded ... \n", len(productMap.m))
}

func loadProductMap() (map[int]models.Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	productList := make([]models.Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}

	prodMap := make(map[int]models.Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}

	return prodMap, nil
}
