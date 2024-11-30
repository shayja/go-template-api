package repository

import (
	"github.com/shayja/go-template-api/model"
)

type ProductRepositoryInterface interface {
	GetAll(page int)([]model.Product, error)
	GetSingle(string) (model.Product, error)
	Create(model.ProductRequest) (string, error)
	Update(string, model.ProductRequest) (model.Product, error)
	UpdatePrice(id string, post model.ProductRequestPrice) (model.Product, error)
	UpdateImage(id string, post model.ProductRequestImage) (model.Product, error) 
	Delete(string) bool
}
