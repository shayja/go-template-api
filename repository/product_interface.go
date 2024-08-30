package repository

import (
	"github.com/shayja/go-template-api/model"
)


type ProductRepositoryInterface interface {
	GetAll(page int)([]model.Product, error)
	GetSingle(string) (model.Product, error)
	Create(model.ValidateProduct) (string, error)
	Update(string, model.ValidateProduct) (model.Product, error)
	UpdatePrice(id string, post model.ValidateProductPrice) (model.Product, error) 
	Delete(string) bool
}
