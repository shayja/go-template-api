// adapters/repositories/user_interface.go
package repositories

import "github.com/shayja/go-template-api/pkg/entities"

type UserRepositoryInterface interface {
	GetById(id string) (entities.User, error)
	GetByUsername(username string) (entities.User, error)
	Create(user entities.User) (string, error)
	ValidatePassword(user entities.User, password string) (error)
}
