// adapters/repositories/user_repository.go
package repositories

import (
	"database/sql"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Db *sql.DB
}

// Get a single item by id
func (m *UserRepository) GetUserById(id string) (*entities.User, error){
	SQL := `SELECT * FROM get_user($1)`
	query, err := m.Db.Query(SQL, id)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	user := &entities.User{}
	if query != nil {
		for query.Next() {
			err := query.Scan(&user.Id, &user.Username, &user.Password, &user.Mobile, &user.Name, &user.Email,  &user.UpdatedAt, &user.CreatedAt)
			if err != nil {
				fmt.Print(err)
				return nil, err
			}
		}
	}
	return user, nil
}

func (m *UserRepository) GetUserByUsername(username string) (*entities.User, error) {
	SQL := `SELECT * FROM get_user_by_username($1)`
	query, err := m.Db.Query(SQL, username)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	user := &entities.User{}
	if query != nil {
		for query.Next() {
			err := query.Scan(&user.Id, &user.Username, &user.Password, &user.Mobile, &user.Name, &user.Email,  &user.UpdatedAt, &user.CreatedAt)
			if err != nil {
				fmt.Print(err)
				return nil, err
			}
		}
	}
	return user, nil
}

func (m *UserRepository) GetUserByMobile(mobile string) (string, error) {
	SQL := `SELECT * FROM get_user_by_mobile($1)`
	query, err := m.Db.Query(SQL, mobile)
	userId := ""
	if err != nil {
		fmt.Print(err)
		return userId, err
	}

	if query != nil {
		for query.Next() {
			err := query.Scan(&userId)
			if err != nil {
				fmt.Print(err)
				return userId, err
			}
		}
	}
	return userId, nil
}

func (m *UserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	err := OnBeforeSave(user)
	if err != nil {
		return nil, err
	}

	var lastInsertId string
	db_err := m.Db.QueryRow("CALL users_insert($1, $2, $3, $4, $5, $6, $7)", user.Username, user.Password, user.Mobile, user.Name, user.Email, time.Now(), user.Id).Scan(&lastInsertId)
	
	if db_err != nil {
		fmt.Print(db_err)
		return user, db_err
	}

	fmt.Printf("user %s created successfully (new id is %s)\n", user.Name, lastInsertId)

	// return the id of the new row
	return user, nil
}

func OnBeforeSave(user *entities.User) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Id = utils.CreateNewUUID().String()
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}


func (m *UserRepository) ValidatePassword(user *entities.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
