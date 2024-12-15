// adapters/repositories/product_repository.go
package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/errors"
	"github.com/shayja/go-template-api/internal/utils"
)


type ProductRepository struct {
	Db *sql.DB
}

const PAGE_SIZE = 20

// Get all product items
func (m *ProductRepository) GetAll(page int) ([]*entities.Product, error) {

	offset := PAGE_SIZE * (page - 1)
	// Optional: Use query text or call a DB Function.
	//SQL := `SELECT id, name, description, image, price, sku FROM "products" ORDER BY "id" LIMIT $2 OFFSET $1`
	SQL := `SELECT * FROM get_products($1, $2)`
	query, err := m.Db.Query(SQL, offset, PAGE_SIZE)

    if err != nil {
		fmt.Print(err)
		return nil, err
	}
    defer query.Close()
	
	var products []*entities.Product
	
	if query != nil {
		for query.Next() {
			product := &entities.Product{}
			err := query.Scan(&product.Id, &product.Name, &product.Description, &product.ImageURL, &product.Price, &product.Sku, &product.UpdatedAt, &product.CreatedAt)
			if err != nil {
				fmt.Print(err)
				return nil, errors.ErrDatabase
			}
			products = append(products, product)
		}
	}
	return products, nil
}

// Get a single item by id
func (m *ProductRepository) GetById(id string) (*entities.Product, error) {
	//query, err := m.Db.Query("SELECT id, name, description, image, price, sku FROM products WHERE id = $1", id)
	SQL := `SELECT * FROM get_product($1)`
	query, err := m.Db.Query(SQL, id)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	product := &entities.Product{}
	if query != nil {
		for query.Next() {
			err := query.Scan(&product.Id, &product.Name, &product.Description, &product.ImageURL, &product.Price, &product.Sku, &product.UpdatedAt, &product.CreatedAt)
			if err != nil {
				fmt.Print(err)
				return nil, errors.ErrDatabase
			}
		}
	}
	return product, nil
}

// Create implements ProductRepositoryInterface
func (m *ProductRepository) Create(product *entities.ProductRequest) (string, error) {
	
	newId := utils.CreateNewUUID().String()
	err := m.Db.QueryRow("CALL products_insert($1, $2, $3, $4, $5, $6, $7)", product.Name, product.Description, product.Price, product.ImageURL, product.Sku, time.Now(), newId).Scan(&newId)
	if err != nil {
		fmt.Print(err)
		return "", errors.ErrDatabase
	}

	fmt.Printf("Product %s created successfully (new id is %s)\n", product.Name, newId)

	// return the id of the new row
	return newId, nil
}

// Update product item
func (m *ProductRepository) Update(id string, product *entities.ProductRequest) (*entities.Product, error) {

	_, err := m.Db.Exec("CALL products_update($1, $2, $3, $4, $5, $6)", id, product.Name, product.Description, product.Price, product.ImageURL, product.Sku)
	if err != nil {
		fmt.Print(err)
		return nil, errors.ErrDatabase
	}

	return m.GetById(id)
}

// Update product item price
func (m *ProductRepository) UpdatePrice(id string, product *entities.ProductPriceRequest) (*entities.Product, error) {
	
	res, err := m.Db.Exec("CALL products_update_price($1, $2)", id, product.Price)
	if err != nil {
		fmt.Print(err, res)
		return nil, errors.ErrDatabase
	}
	return m.GetById(id)
}

// Update product image
func (m *ProductRepository) UpdateImage(id string, product *entities.ProductImageRequest) (*entities.Product, error) {
	_, err := m.Db.Exec("CALL products_update_image($1, $2)", id, product.ImageURL)
	if err != nil {
		fmt.Print(err)
		return nil, errors.ErrDatabase
	}
	return m.GetById(id)
}

// Delete product by id
func (m *ProductRepository) Delete(id string) (bool, error) {
	_, err := m.Db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		fmt.Print(err)
		return false, errors.ErrDatabase
	}
	return true, nil
}

