package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GerardoHP/inventory-service/models"
)

const ProductsBasePath = "products"

type Products struct {
}

func NewProductsHandler() *Products {
	return &Products{}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		products, err := repository.GetProducts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		productsJson, err := json.Marshal(products)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
		return
	case http.MethodPost:
		var product models.Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(bodyBytes, &product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if product.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p, _ := repository.GetNextID()
		product.ProductID = p + 1
		repository.AddProduct(product)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		productJson, _ := json.Marshal(product)
		w.Write(productJson)
		return
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}
