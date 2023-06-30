package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"id"`
	Name        *string            `json:"name" validate:"required,min=3,max=255"`
	Description *string            `json:"description" validate:"required,min=3,max=255"`
	Price       *string            `json:"price" validate:"required"`
	Quantity    int                `json:"qty" validate:"required,min=1"`
	OwnerID     string             `json:"ownerId"`
	OwnerName   string             `json:"ownerName"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	ProductID   string             `json:"productId"`
}

type PurchaseProduct struct {
	ProductId       string    `json:"product_id" validate:"required"`
	ProductName     string    `json:"product_name"`
	Quantity        int       `json:"qty" validate:"required"`
	SellerId        string    `json:"seller_id"`
	SellerName      string    `json:"seller_name"`
	BuyerId         string    `json:"buyer_id"`
	BuyerName       string    `json:"buyer_name"`
	TransactionDate time.Time `json:"transaction_date"`
}

func (p Product) Validate() error {
	if err := validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Price, validation.Required),
		validation.Field(&p.Quantity, validation.Required),
		validation.Field(&p.Description, validation.Required),
	); err != nil {
		return err
	}

	return nil
}

type ProductList struct {
	Data   []*Product `json:"data"`
	Count  int64      `json:"count"`
	Source Source     `json:"source"`
}

type Source string

const (
	DBData    Source = "database"
	CacheData Source = "cache_memory"
)

type ProductFilter struct {
	Limit int64 `json:"limit"`
}
