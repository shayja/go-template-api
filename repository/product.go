package repository

import (
	"database/sql"
	"log"

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
	if query != nil {
		for query.Next() {
			var (
				id		uint
				name	string
				description	string
				image	string
				price	float64
				sku		string
			)
			err := query.Scan(&id, &name, &description, &image, &price, &sku)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			product := model.Product{Id: id, Name: name, Description: description, Image: image, Price: price, Sku: sku}
			products = append(products, product)
		}
	}
	return products, nil
}

// Get a single item by id
func (m *ProductRepository) GetSingle(id uint) (*model.Product, error) {
	query, err := m.Db.Query("SELECT id, name, description, image, price, sku FROM products WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var product model.Product
	if query != nil {
		for query.Next() {
			var (
				id      uint
				name	string
				description	string
				image	string
				price	float64
				sku		string
			)
			err := query.Scan(&id, &name, &description, &image, &price, &sku)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			product = model.Product{Id: id, Name: name, Description: description, Image: image, Price: price, Sku: sku}
		}
	}
	return &product, nil
}

// Create implements ProductRepositoryInterface
func (m *ProductRepository) Create(post model.ValidateProduct) (int, error) {
	stmt, err := m.Db.Prepare(
		`INSERT INTO products (name, description, image, price, sku)
		 VALUES ($1,$2,$3,$4,$5) 
		 RETURNING id;`)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(post.Name, post.Description, post.Image, post.Price, post.Sku).Scan(&id)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	// return the id of the new row
	return id, nil
}

// Update product item
func (m *ProductRepository) Update(id uint, post model.ValidateProduct) (*model.Product, error) {
	_, err := m.Db.Exec("UPDATE products SET name = $1, description = $2, price = $3, image = $4, sku = $5 WHERE id = $6", post.Name, post.Description, post.Price, post.Image, post.Sku, id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return m.GetSingle(id)
}

// Delete product by id
func (m *ProductRepository) Delete(id uint) bool {
	_, err := m.Db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

