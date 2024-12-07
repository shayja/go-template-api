// usecases/user_interface.go
package usecases

import "github.com/shayja/go-template-api/pkg/entities"

type UserServiceInterface interface {
	GetById(id string) (entities.User, error)
}
