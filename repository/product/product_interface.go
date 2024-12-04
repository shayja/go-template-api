package repository

import (
	"github.com/shayja/go-template-api/model"
)

type ProductRepositoryInterface interface {
	GetAll(page int)([]model.Product, error)
	GetSingle(id string) (model.Product, error)
	Create(model.ProductRequest) (string, error)
	Update(string, model.ProductRequest) (model.Product, error)
	UpdatePrice(id string, post model.ProductPriceRequest) (model.Product, error)
	UpdateImage(id string, post model.ProductImageRequest) (model.Product, error) 
	Delete(string) bool
}
