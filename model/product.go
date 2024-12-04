package model

import (
	"time"
)

//Product defines a structure for an item in product catalog
type Product struct {
	Id 			string	`json:"id" validate:"required"`
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	Image       string	`json:"image" validate:"required"`
	Price       float64	`json:"price" validate:"required"`
	Sku         string	`json:"sku" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
ProductRequest represents a product in the shopping list.
*[Name]* is the name of the product.
*[Description]* is the description of the product.
? in a real world app, the image/s will stored in a separated table/
*[Image]* is the image of the product, a base64 string.
TODO: store price in product prices table includes all price history, currency etc.
*[Price]* is the price of the product.
*[Sku]* is the catalog number of the product.
*/
type ProductRequest struct {
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	Image       string	`json:"image" validate:"required"`
	Price       float64	`json:"price" validate:"required"`
	Sku         string	`json:"sku" validate:"required"`
}

type ProductUri struct {
	ID string `uri:"id" binding:"required"`
}

type ProductPriceRequest struct {
	Price       float64	`json:"price" validate:"required"`
}

type ProductImageRequest struct {
	Image       string `json:"image" validate:"required"`
}
