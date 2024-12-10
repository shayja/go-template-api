package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/shayja/go-template-api/internal/entities"
)


func init() {
    // Load the .env file before running tests
    err := godotenv.Load("../../../.env")
    if err != nil {
       panic(err)
    }
}

// MockUserInteractor mocks the UserInteractor
type MockUserInteractor struct {
	mock.Mock
}

func (m *MockUserInteractor) GetUserById(id string) (*entities.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserInteractor) GetUserByUsername(username string) (*entities.User, error) {
    args := m.Called(username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserInteractor) GetUserByMobile(mobile string) (string, error) {
	args := m.Called(mobile)
    if args.Get(0) == nil {
        return "", args.Error(1)
    }
    return args.Get(0).(string), args.Error(1)
}


func (m *MockUserInteractor) ValidatePassword(user *entities.User, password string) error {
	args := m.Called(user, password)
	return args.Error(0)
}

func (m *MockUserInteractor) RegisterUser(userReq *entities.UserRequest) (*entities.User, error) {
	args := m.Called(userReq)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func TestLoginSuccess(t *testing.T) {
	mockInteractor := new(MockUserInteractor)
	controller := &UserController{UserInteractor: mockInteractor}

	user := &entities.User{
		Id:       "1",
		Username: "testuser",
		Password: "hashedpassword",
	}
	input := entities.AuthenticationInput{Username: "testuser", Password: "password"}

	mockInteractor.On("GetUserByUsername", "testuser").Return(user, nil)
	mockInteractor.On("ValidatePassword", user, "password").Return(nil)

	router := gin.Default()
	router.POST("/login", controller.Login)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockInteractor.AssertExpectations(t)
}

func TestLoginUserNotFound(t *testing.T) {
	mockInteractor := new(MockUserInteractor)
	controller := &UserController{UserInteractor: mockInteractor}

	mockInteractor.On("GetUserByUsername", "unknownuser").Return(nil, errors.New("user not found"))

	router := gin.Default()
	router.POST("/login", controller.Login)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(entities.AuthenticationInput{Username: "unknownuser", Password: "password"})
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockInteractor.AssertExpectations(t)
}

func TestRegisterUserSuccess(t *testing.T) {
	mockInteractor := new(MockUserInteractor)
	controller := &UserController{UserInteractor: mockInteractor}

	userReq := entities.UserRequest{Name: "John", Username: "john123", Email: "john@example.com", Password: "secure", Mobile: "1234567890"}
	createdUser := &entities.User{Id: "1", Username: "john123"}

	mockInteractor.On("RegisterUser", &userReq).Return(createdUser, nil)

	router := gin.Default()
	router.POST("/register", controller.RegisterUser)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(userReq)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockInteractor.AssertExpectations(t)
}

func TestRegisterUserError(t *testing.T) {
	mockInteractor := new(MockUserInteractor)
	controller := &UserController{UserInteractor: mockInteractor}

	userReq := entities.UserRequest{Name: "John", Username: "john123", Email: "john@example.com", Password: "secure", Mobile: "1234567890"}

	mockInteractor.On("RegisterUser", &userReq).Return(nil, errors.New("failed to create user"))

	router := gin.Default()
	router.POST("/register", controller.RegisterUser)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(userReq)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockInteractor.AssertExpectations(t)
}