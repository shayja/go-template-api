package repositories

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shayja/go-template-api/pkg/entities"
)

func TestGetAllProduct(t *testing.T) {
	// Create a new SQL mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new ProductImpl instance with the mock database
	repo := &ProductRepository{Db: db}

	time1, err := time.Parse(time.RFC1123, "Sun, 11 Dec 2024 12:21:00 UTC")
	time2, err := time.Parse(time.RFC1123, "Sun, 12 Dec 2024 12:22:00 UTC")

	// Define the expected rows and columns
	rows := sqlmock.NewRows([]string{"id", "name", "description", "image", "price", "sku", "created_at", "updated_at"}).
		AddRow("2d248bb4-e831-44b1-8595-446d460cc511", "Product 1",  "Desc 1", "image1.gif", 10001.01, "10000001", time1, time1).
		AddRow("48dd8c7a-9ac1-4263-88e4-bb01b5e29001", "Product 2",  "Desc 2", "image2.gif", 10002.02, "10000002", time2, time2)

	// Set up the mock expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM get_products($1, $2)`)).WillReturnRows(rows)

	// Call the method you want to test
	products, err := repo.GetAll(1)
	if err != nil {
		t.Errorf("Error calling GetAllProduct: %v", err)
	}

	// Define the expected result
	expectedProducts := []entities.Product{
		{
			Id: "2d248bb4-e831-44b1-8595-446d460cc511",
			Name: "Product 1",
			Description: "Desc 1",
			Image: "image1.gif",
			Price: 10001.01,
			Sku: "10000001",
			CreatedAt: time1,
			UpdatedAt: time1,
		},
		{
			Id: "48dd8c7a-9ac1-4263-88e4-bb01b5e29001",
			Name: "Product 2",
			Description: "Desc 2",
			Image: "image2.gif",
			Price: 10002.02,
			Sku: "10000002",
			CreatedAt: time2,
			UpdatedAt: time2,
		},
	}

	// Compare the actual result to the expected result
	if len(products) != len(expectedProducts) {
		t.Errorf("Expected %d products, but got %d", len(expectedProducts), len(products))
	}

	for i, actual := range products {
		expected := expectedProducts[i]
		if actual != expected {
			t.Errorf("Mismatch in product data at index %d\nExpected: %+v\nActual:   %+v", i, expected, actual)
		}
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
