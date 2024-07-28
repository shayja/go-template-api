package repository

import "github.com/shayja/go-template-api/model"


type ProductRepositoryInterface interface {
	GetAll(page int)([]model.Product, error)
	GetSingle(int64) (model.Product, error)
	Create(model.ValidateProduct) (int64, error)
	Update(int64, model.ValidateProduct) (model.Product, error)
	UpdatePrice(id int64, post model.ValidateProductPrice) (model.Product, error) 
	Delete(int64) bool
}
