package model

import (
	"time"
)

type Product struct {
	Id 			string	`json:"id" validate:"required"`
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	Image       string	`json:"image" validate:"required"`
	Price       float64	`json:"price" validate:"required"`
	Sku         string	`json:"sku" validate:"required"`
	CreateDate  time.Time `json:"create_date"`
}

// ValidateProduct represents a product in the shopping list.
// [Name] is the name of the product.
// [Description] is the description of the product.
// [Image] is the image of the product, a base64 string.
// [Price] is the price of the product.
// [Sku] is the catalog number of the product.
type ValidateProduct struct {
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	Image       string	`json:"image" validate:"required"`
	Price       float64	`json:"price" validate:"required"`
	Sku         string	`json:"sku" validate:"required"`
}

type ProductUri struct {
	ID string `uri:"id" binding:"required"`
}

type ValidateProductPrice struct {
	Price       float64	`json:"price" validate:"required"`
}

type ValidateProductImage struct {
	Image       string `json:"image" validate:"required"`
}
