package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/GerardoHP/inventory-service/models"
)

type Product struct {
}

var (
	separator string = "products/"
)

func NewProductHandler() *Product {
	return &Product{}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodHead:
	case http.MethodPost:
	case http.MethodPatch:
	case http.MethodConnect:
	case http.MethodOptions:
	case http.MethodTrace:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	product, err := getProduct(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPut {
		var updatedProduct models.Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = repository.UpdateProduct(updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updatedProduct
	}

	if r.Method == http.MethodDelete {
		repository.DeleteProduct(product.ProductID)
		product = nil
	}

	productJson, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(productJson)
}

func getProduct(w http.ResponseWriter, r *http.Request) (*models.Product, error) {
	urlPatSegments := strings.Split(r.URL.Path, separator)
	index := len(urlPatSegments) - 1
	productID, err := strconv.Atoi(urlPatSegments[index])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return nil, err
	}

	product, _ := repository.GetProductById(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return nil, err
	}

	return product, nil
}
