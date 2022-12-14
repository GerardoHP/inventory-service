package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GerardoHP/inventory-service/handlers"
	"github.com/GerardoHP/inventory-service/middleware"
)

const (
	apiBasePath = "/api"
)

func SetupRoutes(apiBasePath string) {
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, handlers.ProductsBasePath), middleware.CorsHandler(handlers.NewProductsHandler()))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, handlers.ProductsBasePath), middleware.CorsHandler(handlers.NewProductHandler()))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, handlers.ReceiptPath), middleware.CorsHandler(handlers.NewReceiptsHandler()))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, handlers.ReceiptPath), middleware.CorsHandler(handlers.NewReceiptHandler()))
}

func main() {
	SetupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
