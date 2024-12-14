package repositories_test

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	repositories "github.com/shayja/go-template-api/internal/adapters/repositories/user"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func mockRepository(db *sql.DB) *repositories.UserRepository {
	return &repositories.UserRepository{
		Db:                db,
		HashPassword:      func(password string) (string, error) { hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); return string(hash), err },
		GenerateUUID:      func() string { return "mock-uuid" },
		GenerateTimestamp: func() time.Time { return time.Date(2024, 12, 14, 0, 0, 0, 0, time.UTC) },
	}
}

func TestGetUserById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := mockRepository(db)

	mockRows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "name", "email", "updated_at", "created_at"}).
		AddRow("1", "testuser", "hashedpassword", "1234567890", "Test User", "test@example.com", time.Now(), time.Now())

	mock.ExpectQuery("SELECT \\* FROM get_user\\(\\$1\\)").
		WithArgs("1").
		WillReturnRows(mockRows)

	user, err := repo.GetUserById("1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetUserById_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := mockRepository(db)

	mock.ExpectQuery("SELECT \\* FROM get_user\\(\\$1\\)").
		WithArgs("1").
		WillReturnError(errors.New("query error"))

	user, err := repo.GetUserById("1")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestOnBeforeSave(t *testing.T) {
	repo := mockRepository(nil)

	user := &entities.User{
		Username: " testuser ",
		Password: "password",
	}

	err := repo.OnBeforeSave(user)
	assert.NoError(t, err)
	assert.Equal(t, "mock-uuid", user.Id)
	assert.Equal(t, "testuser", user.Username)
	assert.True(t, strings.HasPrefix(user.Password, "$2a$"))
	assert.Equal(t, time.Date(2024, 12, 14, 0, 0, 0, 0, time.UTC), user.CreatedAt)
}
func TestValidatePassword_Success(t *testing.T) {
    password := "password"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    user := &entities.User{Password: string(hash)}

    repo := &repositories.UserRepository{}
    err := repo.ValidatePassword(user, password)
    assert.NoError(t, err)
}

func TestValidatePassword_Error(t *testing.T) {
    password := "password"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    user := &entities.User{Password: string(hash)}

    repo := &repositories.UserRepository{}
    err := repo.ValidatePassword(user, "wrongpassword")
    assert.Error(t, err)
}
