// usecases/user_usecase.go
package usecases

import (
	"github.com/shayja/go-template-api/pkg/entities"
)

type ProductRepository interface {
	GetAll(page int) ([]*entities.Product, error)
	GetById(id string) (*entities.Product, error)
	Create(product *entities.ProductRequest) (string, error)
	Update(id string, product *entities.ProductRequest) (*entities.Product, error)
	UpdatePrice(id string, product *entities.ProductPriceRequest) (*entities.Product, error)
	UpdateImage(id string, product *entities.ProductImageRequest) (*entities.Product, error)
	Delete(id string) bool
}
	
type ProductInteractor struct {
    ProductRepository ProductRepository
}

func (uc *ProductInteractor) GetAll(page int) ([]*entities.Product, error) {
	return uc.ProductRepository.GetAll(page)
}

func (uc *ProductInteractor) GetById(id string) (*entities.Product, error) {
	return uc.ProductRepository.GetById(id)
}

func (uc *ProductInteractor) Create(product *entities.ProductRequest) (string, error) {
	return uc.ProductRepository.Create(product)
}

func (uc *ProductInteractor) Update(id string, product *entities.ProductRequest) (*entities.Product, error) {
	return uc.ProductRepository.Update(id, product)
}

func (uc *ProductInteractor) UpdatePrice(id string, product *entities.ProductPriceRequest) (*entities.Product, error) {
	return uc.ProductRepository.UpdatePrice(id, product)
}

func (uc *ProductInteractor) UpdateImage(id string, product *entities.ProductImageRequest) (*entities.Product, error) {
	return uc.ProductRepository.UpdateImage(id, product)
}

func (uc *ProductInteractor) Delete(id string) bool {
	return uc.ProductRepository.Delete(id)
}

