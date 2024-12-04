package repository

import (
	"database/sql"
	"html"
	"log"
	"strings"
	"time"

	"github.com/shayja/go-template-api/model"
	"github.com/shayja/go-template-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &UserRepository{Db: db}
}

// Get a single item by id
func (m *UserRepository) GetById(id string) (model.User, error) {
	SQL := `SELECT * FROM get_user($1)`
	query, err := m.Db.Query(SQL, id)
	if err != nil {
		log.Fatal(err)
		return model.User{}, err
	}

	var user model.User
	if query != nil {
		for query.Next() {
			err := query.Scan(&user.Id, &user.Username, &user.Password, &user.Mobile, &user.Name, &user.Email,  &user.UpdatedAt, &user.CreatedAt)
			if err != nil {
				log.Fatal(err)
				return model.User{}, err
			}
		}
	}
	return user, nil
}

// GetByUsername implements UserRepositoryInterface.
func (m *UserRepository) GetByUsername(username string) (model.User, error) {
	SQL := `SELECT * FROM get_user_by_username($1)`
	query, err := m.Db.Query(SQL, username)
	if err != nil {
		log.Fatal(err)
		return model.User{}, err
	}

	var user model.User
	if query != nil {
		for query.Next() {
			err := query.Scan(&user.Id, &user.Username, &user.Password, &user.Mobile, &user.Name, &user.Email,  &user.UpdatedAt, &user.CreatedAt)
			if err != nil {
				log.Fatal(err)
				return model.User{}, err
			}
		}
	}
	return user, nil
}


// Create user implements UserRepositoryInterface.
func (m *UserRepository) Create(user model.User) (string, error) {
	err := OnBeforeSave(&user)
	if err != nil {
		return "", err
	}

	var newId string
	db_err := m.Db.QueryRow("CALL users_insert($1, $2, $3, $4, $5, $6, $7)", user.Username, user.Password, user.Mobile, user.Name, user.Email, time.Now(), user.Id).Scan(&newId)
	
	if db_err != nil {
		log.Fatal(db_err)
		return newId, db_err
	}

	log.Printf("user %s created successfully (new id is %s)\n", user.Name, newId)

	// return the id of the new row
	return newId, nil
}

func OnBeforeSave(user *model.User) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Id = utils.CreateNewUUID().String()
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}


func (m *UserRepository) ValidatePassword(user model.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
