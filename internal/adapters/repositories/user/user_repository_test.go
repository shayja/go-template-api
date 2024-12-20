package repositories_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	repositories "github.com/shayja/go-template-api/internal/adapters/repositories/user"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupMock() (*sql.DB, sqlmock.Sqlmock, *repositories.UserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("Error creating mock database:", err)
	}
	repo := &repositories.UserRepository{Db: db}
	return db, mock, repo
}

func TestGetUserById_Success(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "first_name", "last_name", "email", "otp_types", "verified", "verified_at", "updated_at", "created_at"}).
		AddRow("1", "testuser", "passwordHash", "1234567890", "John", "Doe", "john.doe@example.com", "{1,2,3}", true, time.Now(), time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM get_user\(\$1\)`).
		WithArgs("1").
		WillReturnRows(mockRows)

	user, err := repo.GetUserById("1")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "1", user.Id)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, []int{1, 2, 3}, user.OtpTypes)
}

func TestGetUserById_NotFound(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "first_name", "last_name", "email", "otp_types", "verified", "verified_at", "updated_at", "created_at"})

	mock.ExpectQuery(`SELECT \* FROM get_user\(\$1\)`).
		WithArgs("99").
		WillReturnRows(mockRows)

	user, err := repo.GetUserById("99")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user with id 99 not found")
}

func TestGetUserByUsername_Success(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "first_name", "last_name", "email", "otp_types", "verified", "verified_at", "updated_at", "created_at"}).
		AddRow("1", "testuser", "passwordHash", "1234567890", "John", "Doe", "john.doe@example.com", "{1,2,3}", true, time.Now(), time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM get_user_by_username\(\$1\)`).
		WithArgs("testuser").
		WillReturnRows(mockRows)

	user, err := repo.GetUserByUsername("testuser")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "first_name", "last_name", "email", "otp_types", "verified", "verified_at", "updated_at", "created_at"})

	mock.ExpectQuery(`SELECT \* FROM get_user_by_username\(\$1\)`).
		WithArgs("unknown").
		WillReturnRows(mockRows)

	user, err := repo.GetUserByUsername("unknown")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user with username unknown not found")
}

func TestSaveOTP_Success(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	mock.ExpectQuery(`CALL otpcodes_insert\(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs("userId", "1234567890", "otp123", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("newOtpId"))

	otp := &entities.OTP{
		UserId:     "userId",
		Mobile:     "1234567890",
		OTP:        "otp123",
		Expiration: time.Now().Add(5 * time.Minute),
	}

	err := repo.SaveOTP(otp)

	assert.NoError(t, err)
}

func TestValidateOTP_Success(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	mock.ExpectQuery(`SELECT expiration FROM otpcodes WHERE mobile = \$1 AND otp = \$2`).
		WithArgs("1234567890", "otp123").
		WillReturnRows(sqlmock.NewRows([]string{"expiration"}).AddRow(time.Now().Add(5 * time.Minute)))

	valid, err := repo.ValidateOTP("1234567890", "otp123")

	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestValidateOTP_Expired(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	mock.ExpectQuery(`SELECT expiration FROM otpcodes WHERE mobile = \$1 AND otp = \$2`).
		WithArgs("1234567890", "otp123").
		WillReturnRows(sqlmock.NewRows([]string{"expiration"}).AddRow(time.Now().Add(-5 * time.Minute)))

	valid, err := repo.ValidateOTP("1234567890", "otp123")

	assert.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "OTP expired")
}


func TestGetUserByMobile_Success(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	// Mock the database query
	rows := sqlmock.NewRows([]string{"id", "username", "password", "mobile", "first_name", "last_name", "email", "otpTypes", "verified", "verified_at", "updated_at", "created_at"}).
		AddRow("1", "testuser", "hashedpassword", "123456789", "Test", "User", "test@example.com", "{1,2,3}", true, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM get_user_by_mobile\\(\\$1\\)").
		WithArgs("123456789").
		WillReturnRows(rows)

	// Test the GetUserByMobile method
	user, err := repo.GetUserByMobile("123456789")
	assert.NoError(t, err)
	assert.Equal(t, "1", user.Id)
	assert.Equal(t, "123456789", user.Mobile)
}

func TestCreateUser_Success(t *testing.T) {
	db, mock, repo := setupMock()
	defer db.Close()

	// Create a test user
	user := &entities.User{
		Username: "newuser",
		Password: "password123",
		Mobile:   "987654321",
		FirstName:     "New",
		LastName:     "User here",
		Email:    "newuser@example.com",
		CreatedAt: time.Now(),
	}

	// Mock the database insert
	mock.ExpectQuery(`CALL users_insert\(\$1, \$2, \$3, \$4, \$5, \$6, \$7 \$8\)`).
		WithArgs(
			user.Username,
			sqlmock.AnyArg(), // Allow any password hash
			user.Mobile,
			user.FirstName,
			user.LastName,
			user.Email,
			user.CreatedAt,
			sqlmock.AnyArg(), // Allow any ID (it will be a UUID)
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	// Test the CreateUser method
	createdUser, err := repo.CreateUser(user)

	if err != nil {
		fmt.Println("Error users_insert:", err)
	}

	//assert.NoError(t, err)
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
	_, _, repo := setupMock()
 
	// Mock bcrypt hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Create a user with the hashed password
	password := string(hashedPassword)

	// Test valid password
	err := repo.ValidatePassword(password, "password123")
	assert.NoError(t, err)

	// Test invalid password
	err = repo.ValidatePassword(password, "wrongpassword")
	assert.Error(t, err)
}

func TestValidatePassword_Error(t *testing.T) {
	_, _, repo := setupMock()
    password := "password"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    userPassword := string(hash)
    err := repo.ValidatePassword(userPassword, "wrongpassword")
    assert.Error(t, err)
}


func TestOnBeforeSave_Success(t *testing.T) {
	_, _, repo := setupMock()
	// Create a user instance
	user := &entities.User{
		Username: " testuser ",
		Password: "password123",
		Mobile:   "987654321",
		FirstName:     "User",
		LastName: "Name",
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