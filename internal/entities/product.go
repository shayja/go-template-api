// internal/entities/product.go
package entities

import (
	"time"
)

// Product model info
//
//  @description    Product model contains the information about the product
//  @description    to be added to the store
type Product struct {
	Id 			string	`json:"id" validate:"required"`
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	ImageURL    string	`json:"image" validate:"required"`
	Price       float64	`json:"price" validate:"required"`
	Sku         string	`json:"sku" validate:"required"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}

type ProductRequest struct {
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	ImageURL    string	`json:"image" validate:"required"`
	Price       float64	`json:"price" validate:"required"`
	Sku         string	`json:"sku" validate:"required"`
}


type ProductUri struct {
	Id string `uri:"id" binding:"required"`
}

type ProductPriceRequest struct {
	Price       float64	`json:"price" validate:"required"`
}

type ProductImageRequest struct {
	ImageURL	string `json:"image" validate:"required"`
}
