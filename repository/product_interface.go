package repository

import "github.com/shayja/go-template-api/model"


type ProductRepositoryInterface interface {
	GetAll(page int)([]model.Product, error)
	GetSingle(uint) (model.Product, error)
	Create(model.ValidateProduct) (int, error)
	Update(uint, model.ValidateProduct) (model.Product, error)
	UpdatePrice(id uint, post model.ValidateProductPrice) (model.Product, error) 
	Delete(uint) bool
}
