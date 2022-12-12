package handlers

import "github.com/GerardoHP/inventory-service/data"

var (
	repository data.ProductRepository = data.NewSqlRepository() // data.NewFirestoreRepository() // data.NewJsonRepository()
)
