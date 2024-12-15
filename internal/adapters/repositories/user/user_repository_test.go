// adapters/repositories/user_repository_test.go
package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupTest() (*sql.DB, sqlmock.Sqlmock, *UserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("Error creating mock database:", err)
	}
	userRepo := &UserRepository{
		Db: db,
	}
	return db, mock, userRepo
}

func TestGetUserById_Success(t *testing.T) {
	db, mock, repo := setupTest()
	defer db.Close()

	// Mock the database query
	rows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "name", "email", "updated_at", "created_at"}).
		AddRow("1", "testuser", "hashedpassword", "1234567890", "Test User", "test@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM get_user\\(\\$1\\)").
		WithArgs("1").
		WillReturnRows(rows)

	// Test the GetUserById method
	user, err := repo.GetUserById("1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetUserById_Error(t *testing.T) {
	db, mock, repo := setupTest()
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM get_user\\(\\$1\\)").
		WithArgs("1").
		WillReturnError(errors.New("query error"))

	user, err := repo.GetUserById("1")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserByUsername_Success(t *testing.T) {
	db, mock, repo := setupTest()
	defer db.Close()

	// Mock the database query
	rows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "name", "email", "updated_at", "created_at"}).
		AddRow("1", "testuser", "hashedpassword", "123456789", "Test User", "test@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM get_user_by_username\\(\\$1\\)").
		WithArgs("testuser").
		WillReturnRows(rows)

	// Test the GetUserByUsername method
	user, err := repo.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetUserByMobile_Success(t *testing.T) {
	db, mock, repo := setupTest()
	defer db.Close()

	// Mock the database query
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	mock.ExpectQuery("SELECT \\* FROM get_user_by_mobile\\(\\$1\\)").
		WithArgs("123456789").
		WillReturnRows(rows)

	// Test the GetUserByMobile method
	userId, err := repo.GetUserByMobile("123456789")
	assert.NoError(t, err)
	assert.Equal(t, "1", userId)
}

func TestCreateUser_Success(t *testing.T) {
	db, mock, repo := setupTest()
	defer db.Close()

	// Create a test user
	user := &entities.User{
		Username: "newuser",
		Password: "password123",
		Mobile:   "987654321",
		Name:     "New User",
		Email:    "newuser@example.com",
		CreatedAt: time.Now(),
	}

	// Mock the database insert
	mock.ExpectQuery("CALL users_insert\\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\)").
		WithArgs(
			user.Username,
			sqlmock.AnyArg(), // Allow any password hash
			user.Mobile,
			user.Name,
			user.Email,
			user.CreatedAt,
			sqlmock.AnyArg(), // Allow any ID (it will be a UUID)
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	// Test the CreateUser method
	createdUser, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	// Ensure the password is hashed
	assert.NotEqual(t, "password123", createdUser.Password)
	assert.True(t, IsValidUUID(createdUser.Id)) // Validate the ID is a UUID
}

func IsValidUUID(uuid string) bool {
    r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
    return r.MatchString(uuid)
}

func TestValidatePassword_Success(t *testing.T) {
	_, _, repo := setupTest()
 
	// Mock bcrypt hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Create a user with the hashed password
	user := &entities.User{
		Password: string(hashedPassword),
	}

	// Test valid password
	err := repo.ValidatePassword(user, "password123")
	assert.NoError(t, err)

	// Test invalid password
	err = repo.ValidatePassword(user, "wrongpassword")
	assert.Error(t, err)
}

func TestValidatePassword_Error(t *testing.T) {
	_, _, repo := setupTest()
    password := "password"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    user := &entities.User{Password: string(hash)}
    err := repo.ValidatePassword(user, "wrongpassword")
    assert.Error(t, err)
}


func TestOnBeforeSave_Success(t *testing.T) {
	_, _, repo := setupTest()
	// Create a user instance
	user := &entities.User{
		Username: " testuser ",
		Password: "password123",
		Mobile:   "987654321",
		Name:     "User Name",
		Email:    "user@example.com",
	}

	// Mock the OnBeforeSave method
	err := repo.OnBeforeSave(user)
	assert.NoError(t, err)

	// Check that the password is hashed
	assert.NotEqual(t, "password123", user.Password)
	assert.NotEmpty(t, user.Id)
	assert.Equal(t, "testuser", user.Username)
	assert.True(t, strings.HasPrefix(user.Password, "$2a$"))
}