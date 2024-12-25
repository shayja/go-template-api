// usecases/order_usecase.go
package usecases

import (
	"github.com/shayja/go-template-api/internal/entities"
)

type OrderRepository interface {
	GetAllOrders(page int, userId string) ([]*entities.Order, error)
	GetById(id string) (*entities.Order, error)
	Create(orderRequest *entities.OrderRequest) (string, error)
	UpdateStatus(id string, status int) (*entities.Order, error)
}

type OrderUsecase struct {
	OrderRepo OrderRepository
}

func (uc *OrderUsecase) GetOrders(page int, userId string) ([]*entities.Order, error) {
	return uc.OrderRepo.GetAllOrders(page, userId)
}

func (uc *OrderUsecase) GetById(id string) (*entities.Order, error) {
	return uc.OrderRepo.GetById(id)
}

func (uc *OrderUsecase) Create(orderRequest *entities.OrderRequest) (string, error) {
	return uc.OrderRepo.Create(orderRequest)
}

func (uc *OrderUsecase) UpdateStatus(id string, status int) (*entities.Order, error) {
	return uc.OrderRepo.UpdateStatus(id, status)
}