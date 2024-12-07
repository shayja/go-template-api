// adapters/repositories/product_interface.go
package repositories

import "github.com/shayja/go-template-api/pkg/entities"

type ProductRepositoryInterface interface {
	GetAll(page int)([]entities.Product, error)
	GetSingle(id string) (entities.Product, error)
	Create(product *entities.ProductRequest) (string, error)
	Update(id string, product *entities.ProductRequest) (entities.Product, error)
	UpdatePrice(id string, post *entities.ProductPriceRequest) (entities.Product, error)
	UpdateImage(id string, post *entities.ProductImageRequest) (entities.Product, error)
	Delete(string) bool
}
