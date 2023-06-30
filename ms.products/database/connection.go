package database

import (
	"context"

	"github.com/fredele20/microservice-practice/ms.products/models"
)

type DBInterface interface {
	CreateProduct(ctx context.Context, payload *models.Product) (*models.Product, error)
	GetProducts(ctx context.Context, filter models.ProductFilter) (*models.ProductList, error)
}
