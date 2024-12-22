package entities

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type OrderRequest struct {
	UserId       string            `json:"user_id"`
	TotalPrice   float64           `json:"total_price"`
	Status       int               `json:"status"`
	OrderDetails []OrderDetailItem `json:"order_details"`
}

type Order struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderDetailItem struct {
	ProductId  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

// Convert order details to database-compatible array
func (v OrderDetailItem) Value() (driver.Value, error) {
    return []byte(fmt.Sprintf("(%s,%d,%f)", v.ProductId, v.Quantity, v.UnitPrice)), nil
}