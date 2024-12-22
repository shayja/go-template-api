package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/errors"
	"github.com/shayja/go-template-api/internal/utils"
)

type OrderRepository struct {
	Db *sql.DB
}

const PAGE_SIZE = 20

// Get all orders
func (r *OrderRepository) GetAll(page int, user_id string) ([]*entities.Order, error) {
	offset := PAGE_SIZE * (page - 1)
	SQL := `SELECT * FROM get_user_orders($1, $2)`
	query, err := r.Db.Query(SQL, offset, PAGE_SIZE)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer query.Close()

	var orders []*entities.Order
	for query.Next() {
		order := &entities.Order{}
		err := query.Scan(&order.Id, &order.UserId, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			fmt.Print(err)
			return nil, errors.ErrDatabase
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// Get order by ID
func (r *OrderRepository) GetById(id string) (*entities.Order, error) {
	SQL := `SELECT * FROM get_order($1)`
	query, err := r.Db.Query(SQL, id)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer query.Close()

	order := &entities.Order{}
	if query.Next() {
		err := query.Scan(&order.Id, &order.UserId, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			fmt.Print(err)
			return nil, errors.ErrDatabase
		}
	}
	return order, nil
}

// Create a new order
func (r *OrderRepository) Create(orderRequest *entities.OrderRequest) (string, error) {
	newId := utils.CreateNewUUID().String()
	_, err := r.Db.Exec(
		`CALL orders_insert($1, $2, $3, $4::order_detail_type[], $5)`,
		orderRequest.UserId,
		orderRequest.TotalPrice,
		orderRequest.Status,
		pq.Array(orderRequest.OrderDetails),
		&newId,
	)

	if err != nil {
		fmt.Print(err)
		return "", errors.ErrDatabase
	}

	fmt.Printf("Order %s created successfully\n", newId)
	return newId, nil
}

// Update order status
func (r *OrderRepository) UpdateStatus(id string, status int) (*entities.Order, error) {
	_, err := r.Db.Exec("CALL orders_update_status($1, $2)", id, status)
	if err != nil {
		fmt.Print(err)
		return nil, errors.ErrDatabase
	}
	return r.GetById(id)
}