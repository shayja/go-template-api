// usecases/user_usecase.go
package usecases

import (
	"database/sql"

	repositories "github.com/shayja/go-template-api/pkg/adapters/repositories/user"
	"github.com/shayja/go-template-api/pkg/entities"
)

type UserService struct {
	Db *sql.DB
}

func CreateUserService(db *sql.DB) UserServiceInterface {
	return &UserService{Db: db}
}

func (m *UserService) GetById(id string) (entities.User, error) {
	DB := m.Db
	repositories := repositories.NewUserRepository(DB)
	res, err := repositories.GetById(id)

    if err != nil {
        return entities.User{}, err
    }
	return res, nil
}

