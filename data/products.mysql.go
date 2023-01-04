package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/GerardoHP/inventory-service/models"
)

var config *Config

type SqlRepository struct {
}

var _ ProductRepository = (*SqlRepository)(nil)
var DbConn *sql.DB

func NewSqlRepository() *SqlRepository {
	return &SqlRepository{}
}

func (SqlRepository) GetProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.QueryTimeout))
	defer cancel()
	results, err := DbConn.QueryContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products
	WHERE deleted <> 1`)
	if err != nil {
		return nil, err
	}

	defer results.Close()
	products := make([]models.Product, 0, len(productMap.m))
	for results.Next() {
		var product models.Product
		results.Scan(
			&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName,
		)

		products = append(products, product)
	}

	return products, nil
}

func (SqlRepository) AddProduct(p models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.QueryTimeout))
	defer cancel()

	query := p.GetProductInsertQuery()
	_, err := DbConn.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (SqlRepository) GetNextID() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.QueryTimeout))
	defer cancel()

	row := DbConn.QueryRowContext(ctx, `SELECT productID FROM products ORDER BY productID DESC LIMIT 1;`)
	product := &models.Product{}
	err := row.Scan(&product.ProductID)
	if err == sql.ErrNoRows {
		return -1, nil
	}

	if err != nil {
		return -1, err
	}

	return product.ProductID, nil
}

func (SqlRepository) GetProductById(id int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.QueryTimeout))
	defer cancel()

	row := DbConn.QueryRowContext(ctx, `SELECT 
	productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products WHERE deleted <> 1 AND productId = ?`, id)

	product := &models.Product{}
	err := row.Scan(
		&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (SqlRepository) UpdateProduct(p models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.QueryTimeout))
	defer cancel()

	_, err := DbConn.ExecContext(ctx, `UPDATE products SET
	manufacturer=?,
	sku=?,
	upc=?,
	pricePerUnit=?,
	quantityOnHand=?,
	productName=? 
	WHERE productId=?
	`,
		p.Manufacturer,
		p.Sku,
		p.Upc,
		p.PricePerUnit,
		p.QuantityOnHand,
		p.ProductName,
		p.ProductID)

	return err
}

func (SqlRepository) DeleteProduct(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.QueryTimeout))
	defer cancel()

	_, err := DbConn.ExecContext(ctx, `UPDATE products SET deleted=1 WHERE productId = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

func (SqlRepository) GetTopTenProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.QueryTimeout))
	defer cancel()
	results, err := DbConn.QueryContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products ORDER BY quantityOnHand DESC LIMIT 10
	WHERE deleted <> 1`)
	if err != nil {
		return nil, err
	}

	defer results.Close()
	products := make([]models.Product, 0, len(productMap.m))
	for results.Next() {
		var product models.Product
		results.Scan(
			&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName,
		)

		products = append(products, product)
	}

	return products, nil
}

func init() {
	var err error
	config = NewConfigFromFile()

	connString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", config.User, config.Pass, config.Host, config.Port, config.DBName)
	fmt.Printf(" conn string %v \n", connString)
	DbConn, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatal(err)
	}

	DbConn.SetMaxOpenConns(config.MaxOpenConnections)
	DbConn.SetMaxIdleConns(config.MaxIdleConnections)
	DbConn.SetConnMaxLifetime(time.Second * time.Duration(config.MaxLifetime))
}
