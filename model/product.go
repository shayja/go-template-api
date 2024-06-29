package model

type Product struct {
	Id 			uint `json:"id"`
	Name        string	`json:"name" validate:"required"`
	Description string	`json:"description" validate:"required"`
	Image       string	`json:"image" validate:"required"`
	Price       float64	`json:"price" validate:"required"`
	Sku         string	`json:"sku" validate:"required"`
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
	ID uint `uri:"id" binding:"required,number"`
}