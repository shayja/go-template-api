package service

import (
	"database/sql"

	"github.com/shayja/go-template-api/model"
	repository "github.com/shayja/go-template-api/repository/user"
)

type UserService struct {
	Db *sql.DB
}

func CreateUserService(db *sql.DB) UserServiceInterface {
	return &UserService{Db: db}
}

func (m *UserService) GetById(id string) (model.User, error) {
	DB := m.Db
	repository := repository.NewUserRepository(DB)
	res, err := repository.GetById(id)

    if err != nil {
        return model.User{}, err
    }
	return res, nil
}

