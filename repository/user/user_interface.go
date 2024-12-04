package repository

import (
	"github.com/shayja/go-template-api/model"
)

type UserRepositoryInterface interface {
	GetById(id string) (model.User, error)
	GetByUsername(username string) (model.User, error)
	Create(user model.User) (string, error)
	ValidatePassword(user model.User, password string) (error)
}
