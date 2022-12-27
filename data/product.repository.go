package data

import "github.com/GerardoHP/inventory-service/models"

type ProductRepository interface {
	GetProducts() ([]models.Product, error)
	AddProduct(p models.Product) error
	GetNextID() (id int, err error)
	GetProductById(id int) (*models.Product, error)
	UpdateProduct(p models.Product) error
	DeleteProduct(id int) error
}
