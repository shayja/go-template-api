package repositories_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	repositories "github.com/shayja/go-template-api/internal/adapters/repositories/product"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repositories.ProductRepository{Db: db}

	mockRows := sqlmock.NewRows([]string{"id", "name", "description", "image", "price", "sku", "updated_at", "created_at"}).
		AddRow("1", "Product1", "Description1", "image1.jpg", 10.5, "SKU1", time.Now(), time.Now()).
		AddRow("2", "Product2", "Description2", "image2.jpg", 20.5, "SKU2", time.Now(), time.Now())

	mock.ExpectQuery("SELECT \\* FROM get_products\\(\\$1, \\$2\\)").
		WithArgs(0, 20).
		WillReturnRows(mockRows)

	products, err := repo.GetAll(1)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Product1", products[0].Name)
	assert.Equal(t, "Product2", products[1].Name)
}

func TestGetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repositories.ProductRepository{Db: db}

	mock.ExpectQuery("SELECT \\* FROM get_products\\(\\$1, \\$2\\)").
		WithArgs(0, 20).
		WillReturnError(errors.New("query error"))

	products, err := repo.GetAll(1)
	assert.Error(t, err)
	assert.Nil(t, products)
}

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repositories.ProductRepository{Db: db}

	mock.ExpectQuery("CALL products_insert\\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\)").
		WithArgs("Product1", "Description1", 10.5, "image1.jpg", "SKU1", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	id, err := repo.Create(&entities.ProductRequest{
		Name:        "Product1",
		Description: "Description1",
		Price:       10.5,
		ImageURL:    "image1.jpg",
		Sku:         "SKU1",
	})
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
}

func TestCreate_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repositories.ProductRepository{Db: db}

	mock.ExpectQuery("CALL products_insert\\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\)").
		WithArgs("Product1", "Description1", 10.5, "image1.jpg", "SKU1", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))

	id, err := repo.Create(&entities.ProductRequest{
		Name:        "Product1",
		Description: "Description1",
		Price:       10.5,
		ImageURL:    "image1.jpg",
		Sku:         "SKU1",
	})
	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repositories.ProductRepository{Db: db}

	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.Delete("1")
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestDelete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repositories.ProductRepository{Db: db}

	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs("1").
		WillReturnError(errors.New("delete error"))

	result, err := repo.Delete("1")
	assert.Error(t, err)
	assert.False(t, result)
}
