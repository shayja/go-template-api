package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/shayja/go-template-api/model"
)

type ProductRepository struct {
	Db *sql.DB
}

const GET_PRODUCTS_QUERY = `CALL get_products($1,$2)`
const PAGE_SIZE = 20

func NewProductRepository(db *sql.DB) ProductRepositoryInterface {
	return &ProductRepository{Db: db}
}

// Get all product items
func (m *ProductRepository) GetAll(page int)([]model.Product, error) {

	offset := PAGE_SIZE * (page - 1)
	// Optional: Use query text or call a DB Function.
	//SQL := `SELECT id, name, description, image, price, sku FROM "products" ORDER BY "id" LIMIT $2 OFFSET $1`
	SQL := `SELECT * FROM get_products($1, $2)`
	query, err := m.Db.Query(SQL, offset, PAGE_SIZE)

    if err != nil {
		log.Fatal(err)
		return nil, err
	}
    defer query.Close()
	
	var products []model.Product
	var product model.Product
	if query != nil {
		for query.Next() {
			err := query.Scan(&product.Id, &product.Name, &product.Description, &product.Image, &product.Price, &product.Sku, &product.CreateDate)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			products = append(products, product)
		}
	}
	return products, nil
}

// Get a single item by id
func (m *ProductRepository) GetSingle(id string) (model.Product, error) {
	//query, err := m.Db.Query("SELECT id, name, description, image, price, sku FROM products WHERE id = $1", id)
	SQL := `SELECT * FROM get_product($1)`
	query, err := m.Db.Query(SQL, id)
	if err != nil {
		log.Fatal(err)
		return model.Product{}, err
	}
	var product model.Product
	if query != nil {
		for query.Next() {
			err := query.Scan(&product.Id, &product.Name, &product.Description, &product.Image, &product.Price, &product.Sku, &product.CreateDate)
			if err != nil {
				log.Fatal(err)
				return model.Product{}, err
			}
		}
	}
	return product, nil
}

// Create implements ProductRepositoryInterface
func (m *ProductRepository) Create(post model.ValidateProduct) (string, error) {
	
	var newId string = ""
	err := m.Db.QueryRow("CALL products_insert($1, $2, $3, $4, $5, $6, $7)", post.Name, post.Description, post.Price, post.Image, post.Sku, time.Now(), newId).Scan(&newId)
	if err != nil {
		log.Fatal(err)
		return newId, err
	}

	log.Printf("Product %s created successfully (new id is %s)\n", post.Name, newId)

	// return the id of the new row
	return newId, nil
}

// Update product item
func (m *ProductRepository) Update(id string, post model.ValidateProduct) (model.Product, error) {

	_, err := m.Db.Exec("CALL products_update($1, $2, $3, $4, $5, $6)", id, post.Name, post.Description, post.Price, post.Image, post.Sku)
	if err != nil {
		log.Fatal(err)
		return model.Product{}, err
	}

	return m.GetSingle(id)
}

// Update product item price
func (m *ProductRepository) UpdatePrice(id string, post model.ValidateProductPrice) (model.Product, error) {
	_, err := m.Db.Exec("CALL products_update_price($1, $2)", id, post.Price)
	if err != nil {
		log.Fatal(err)
		return model.Product{}, err
	}
	return m.GetSingle(id)
}

// Update product image
func (m *ProductRepository) UpdateImage(id string, post model.ValidateProductImage) (model.Product, error) {
	_, err := m.Db.Exec("CALL products_update_image($1, $2)", id, post.Image)
	if err != nil {
		log.Fatal(err)
		return model.Product{}, err
	}
	return m.GetSingle(id)
}

// Delete product by id
func (m *ProductRepository) Delete(id string) bool {
	_, err := m.Db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

