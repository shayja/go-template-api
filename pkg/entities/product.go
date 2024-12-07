package entities

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
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}

type ProductRequest struct {
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	Image       string	`json:"image" validate:"required"`
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
	Image       string `json:"image" validate:"required"`
}
