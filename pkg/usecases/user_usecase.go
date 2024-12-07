// usecases/user_usecase.go
package usecases

import (
	"github.com/shayja/go-template-api/pkg/entities"
)

type UserRepository interface {
	GetUserById(id string) (*entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	ValidatePassword(user *entities.User, password string) error
	CreateUser(user *entities.User) (*entities.User, error)
}

type UserInteractor struct {
    UserRepository UserRepository
}

func (uc *UserInteractor) GetUserById(id string) (*entities.User, error) {
	return uc.UserRepository.GetUserById(id)
}

func (uc *UserInteractor) GetByUsername(username string) (*entities.User, error) {
	return uc.UserRepository.GetByUsername(username)
}

func (uc *UserInteractor) ValidatePassword(user *entities.User, password string) error {
	return uc.UserRepository.ValidatePassword(user, password)
}


func (uc *UserInteractor) RegisterUser(name string, email string, password string) (*entities.User, error) {
    user := &entities.User{
        Name:     name,
        Email:    email,
        Password: password,
    }
    return uc.UserRepository.CreateUser(user)
}

