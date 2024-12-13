package usecases_test

import (
	"errors"
	"testing"

	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository mocks the ProductRepository interface
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll(page int) ([]*entities.Product, error) {
	args := m.Called(page)
	if products, ok := args.Get(0).([]*entities.Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) GetById(id string) (*entities.Product, error) {
	args := m.Called(id)
	if product, ok := args.Get(0).(*entities.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) Create(product *entities.ProductRequest) (string, error) {
	args := m.Called(product)
	return args.String(0), args.Error(1)
}

func (m *MockProductRepository) Update(id string, product *entities.ProductRequest) (*entities.Product, error) {
	args := m.Called(id, product)
	if updatedProduct, ok := args.Get(0).(*entities.Product); ok {
		return updatedProduct, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) UpdatePrice(id string, product *entities.ProductPriceRequest) (*entities.Product, error) {
	args := m.Called(id, product)
	if updatedProduct, ok := args.Get(0).(*entities.Product); ok {
		return updatedProduct, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) UpdateImage(id string, product *entities.ProductImageRequest) (*entities.Product, error) {
	args := m.Called(id, product)
	if updatedProduct, ok := args.Get(0).(*entities.Product); ok {
		return updatedProduct, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) Delete(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}


func TestProductInteractor_GetAll(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	products := []*entities.Product{
		{Id: "1", Name: "Product1"},
		{Id: "2", Name: "Product2"},
	}

	repo.On("GetAll", 1).Return(products, nil)

	result, err := interactor.GetAll(1)

	assert.NoError(t, err)
	assert.Equal(t, products, result)
	repo.AssertExpectations(t)
}

func TestProductInteractor_GetById(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	product := &entities.Product{Id: "1", Name: "Product1"}
	repo.On("GetById", "1").Return(product, nil)

	result, err := interactor.GetById("1")

	assert.NoError(t, err)
	assert.Equal(t, product, result)
	repo.AssertExpectations(t)
}

func TestProductInteractor_Create(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	productRequest := &entities.ProductRequest{Name: "New Product"}
	repo.On("Create", productRequest).Return("1", nil)

	result, err := interactor.Create(productRequest)

	assert.NoError(t, err)
	assert.Equal(t, "1", result)
	repo.AssertExpectations(t)
}

func TestProductInteractor_Update(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	productRequest := &entities.ProductRequest{Name: "Updated Product"}
	updatedProduct := &entities.Product{Id: "1", Name: "Updated Product"}
	repo.On("Update", "1", productRequest).Return(updatedProduct, nil)

	result, err := interactor.Update("1", productRequest)

	assert.NoError(t, err)
	assert.Equal(t, updatedProduct, result)
	repo.AssertExpectations(t)
}

func TestProductInteractor_UpdatePrice(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	priceRequest := &entities.ProductPriceRequest{Price: 99.99}
	updatedProduct := &entities.Product{Id: "1", Price: 99.99}
	repo.On("UpdatePrice", "1", priceRequest).Return(updatedProduct, nil)

	result, err := interactor.UpdatePrice("1", priceRequest)

	assert.NoError(t, err)
	assert.Equal(t, updatedProduct, result)
	repo.AssertExpectations(t)
}

func TestProductInteractor_UpdateImage(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	imageRequest := &entities.ProductImageRequest{ImageURL: "http://example.com/image.jpg"}
	updatedProduct := &entities.Product{Id: "1", ImageURL: "http://example.com/image.jpg"}
	repo.On("UpdateImage", "1", imageRequest).Return(updatedProduct, nil)

	result, err := interactor.UpdateImage("1", imageRequest)

	assert.NoError(t, err)
	assert.Equal(t, updatedProduct, result)
	repo.AssertExpectations(t)
}

func TestProductInteractor_Delete(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	repo.On("Delete", "1").Return(true, nil)

	result, err := interactor.Delete("1")

	assert.NoError(t, err)
	assert.True(t, result)
	repo.AssertExpectations(t)
}

func TestProductInteractor_GetAll_Error(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	repo.On("GetAll", 1).Return(nil, errors.New("database error"))

	result, err := interactor.GetAll(1)

	assert.Error(t, err)           // Assert that an error is returned
	assert.Nil(t, result)          // Assert that the result is nil
	repo.AssertExpectations(t)     // Verify that all mocked methods were called
}

func TestProductInteractor_Delete_Error(t *testing.T) {
	repo := new(MockProductRepository)
	interactor := &usecases.ProductInteractor{ProductRepository: repo}

	repo.On("Delete", "1").Return(false, errors.New("delete error"))

	result, err := interactor.Delete("1")

	assert.Error(t, err)
	assert.False(t, result)
	repo.AssertExpectations(t)
}
