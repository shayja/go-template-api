// usecases/order_usecase.go
package usecases

import (
	"github.com/shayja/go-template-api/internal/entities"
)

type OrderRepository interface {
	GetAll(page int, user_id string) ([]*entities.Order, error)
	GetById(id string) (*entities.Order, error)
	Create(orderRequest *entities.OrderRequest) (string, error)
	UpdateStatus(id string, status int) (*entities.Order, error)
}

type OrderInteractor struct {
	OrderRepository OrderRepository
}

func (uc *OrderInteractor) GetAll(page int, user_id string) ([]*entities.Order, error) {
	return uc.OrderRepository.GetAll(page, user_id)
}

func (uc *OrderInteractor) GetById(id string) (*entities.Order, error) {
	return uc.OrderRepository.GetById(id)
}

func (uc *OrderInteractor) Create(orderRequest *entities.OrderRequest) (string, error) {
	return uc.OrderRepository.Create(orderRequest)
}

func (uc *OrderInteractor) UpdateStatus(id string, status int) (*entities.Order, error) {
	return uc.OrderRepository.UpdateStatus(id, status)
}