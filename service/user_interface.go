package service

import "github.com/shayja/go-template-api/model"

type UserServiceInterface interface {
	GetById(id string) (model.User, error)
}
