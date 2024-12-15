// adapters/repositories/user_repository.go
package repositories

import (
	"database/sql"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/errors"
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

func (m *UserRepository) GetUserByMobile(mobile string) (*entities.User, error) {
	SQL := `SELECT * FROM get_user_by_mobile($1)`
	query, err := m.Db.Query(SQL, mobile)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	if query != nil {
		user := &entities.User{}
		for query.Next() {
			err := query.Scan(&user.Id, &user.Username, &user.Password, &user.Mobile, &user.Name, &user.Email, &user.UpdatedAt, &user.CreatedAt)
			if err != nil {
				fmt.Print(err)
				return nil, err
			}
			return user, nil
		}
		
	}
	return nil, nil
}

func (m *UserRepository) SaveOTP(otp *entities.OTP) error {
	newId := utils.CreateNewUUID().String()
	err := m.Db.QueryRow("CALL otpcodes_insert($1, $2, $3, $4, $5, $6)", otp.UserId, otp.Mobile, otp.OTP, otp.Expiration, time.Now(), newId).Scan(&newId)
	if err != nil {
		fmt.Print(err)
		return errors.ErrDatabase
	}

	fmt.Printf("otp %s created successfully (new id is %s)\n", otp.OTP, newId)

	return nil
}

func (m *UserRepository) ValidateOTP(mobile string, otp string) (bool, error) {
	SQL := `SELECT expiration FROM otpcodes WHERE mobile = $1 AND otp = $2`
	var expiration time.Time
	err := m.Db.QueryRow(SQL, mobile, otp).Scan(&expiration)
	if err != nil {
		fmt.Print(err)
		return false, err
	}

	if time.Now().After(expiration) {
		return false, fmt.Errorf("OTP expired")
	}

	return true, nil
}


func (m *UserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	err := m.OnBeforeSave(user)
	if err != nil {
		return nil, err
	}

	var lastInsertId string
	db_err := m.Db.QueryRow("CALL users_insert($1, $2, $3, $4, $5, $6, $7)",
		user.Username, user.Password, user.Mobile, user.Name, user.Email, user.CreatedAt, user.Id).Scan(&lastInsertId)

	if db_err != nil {
		fmt.Print(db_err)
		return user, db_err
	}

	fmt.Printf("user %s created successfully (new id is %s)\n", user.Name, lastInsertId)

	// return the id of the new row
	return user, nil
}

func (m *UserRepository) OnBeforeSave(user *entities.User) error {

	if user.CreatedAt.IsZero() {
		user.CreatedAt = GenerateTimestamp()
	}
	user.Id = GenerateUUID()
	passwordHash, err := HashPassword(user.Password)
	if err == nil && len(passwordHash) > 0 {
		user.Password = passwordHash
	}
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}


// GetOTP retrieves the OTP and expiration time for a given mobile number
func (m *UserRepository) GetOTP(mobile string) (*entities.OTP, error) {
	SQL := `SELECT 1 FROM otps WHERE mobile = $1`
	query, err := m.Db.Query(SQL, mobile)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	item := &entities.OTP{}
	if query != nil {
		for query.Next() {
			err := query.Scan(&item.Mobile, &item.OTP, &item.Expiration)
			if err != nil {
				fmt.Print(err)
				return nil, err
			}
		}
	}

	if item == nil {
		return nil, errors.ErrOTPNotFound
	}
	return item, nil
}


func (m *UserRepository) ValidatePassword(user *entities.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return string(passwordHash), nil;
}

func GenerateUUID() (string) {
	return utils.CreateNewUUID().String()
}

func GenerateTimestamp() (time.Time) {
	return time.Now()
}