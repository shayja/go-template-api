package entities

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Order represents an order entity.
type Order struct {
	// The UUID of a product
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	Id  string    `json:"id" example:"6204037c-30e6-408b-8aaa-dd8219860b4b" minLength:"36"`
	// The user that created the order
	// example: 451fa817-41f4-40cf-8dc2-c9f22aa98a4f
	// required: true
	UserId string    `json:"user_id" example:"451fa817-41f4-40cf-8dc2-c9f22aa98a4f" minLength:"36"`
	// The total price of the order
	// example: 100.00
	// required: true
	TotalPrice float64  `json:"total_price" example:"100.00" format:"float64"`
	// The status of the order (1=created/pending, 2=processing, 3=completed, 4=cancelled)
	// example: 2
	// required: true
	Status int  `json:"status" example:"1" format:"int32" minimum:"1"`
	// The date and time the order was created
	// example: 2024-07-01T12:00:00Z
	CreatedAt  time.Time `json:"created_at" example:"2024-07-01T12:00:00Z" minLength:"20"`
	// The date and time the order was last updated
	// example: 2025-01-01T12:00:00Z
	UpdatedAt  time.Time `json:"updated_at"`
}

// OrderDetail represents an order line item entity.
type OrderDetail struct {
	// The UUID of an order detail (line item)
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	Id         string  `json:"id"`
	// The UUID of the related order.
	// example: 451fa817-41f4-40cf-8dc2-c9f22aa98a4f
	// required: true
	OrderId    string  `json:"order_id" example:"451fa817-41f4-40cf-8dc2-c9f22aa98a4f" minLength:"36"`
	// The UUID of the related product
	// example: 063d0ff7-e17e-4957-8d92-a988caeda8a1
	// required: true
	ProductId  string  `json:"product_id" example:"063d0ff7-e17e-4957-8d92-a988caeda8a1" minLength:"36"`
	// The quantity of the product
	// example: 2
	// required: true
	Quantity   int `json:"quantity" example:"1" format:"int32" minimum:"1"`
	// The unit price of the product
	// example: 50.00
	// required: true
	UnitPrice  float64 `json:"unit_price" example:"50.00" format:"float64"`
	// The date and time the order detail was created
	// example: 2024-07-01T12:00:00Z
	// required: true
	TotalPrice float64 `json:"total_price" example:"55.00" format:"float64"`
	// The date and time the order detail was created
	// example: 2025-01-01T12:00:00Z
	CreatedAt  time.Time `json:"created_at"`
	// The date and time the order detail was last updated
	// example: 2025-01-01T12:00:00Z
	UpdatedAt  time.Time `json:"updated_at" example:"2024-07-01T12:00:00Z" minLength:"20"`
}

// OrderRequest represents a request to create an order.
type OrderRequest struct {
	// The user that creates the order
	// example: 451fa817-41f4-40cf-8dc2-c9f22aa98a4f
	// required: true
	UserId string `json:"user_id" example:"063d0ff7-e17e-4957-8d92-a988caeda8a1" minLength:"36"`
	// The total price of the order
	// example: 100.00
	// required: true
	TotalPrice   float64 `json:"total_price" example:"100.00" format:"float64"`
	// The status of the order (1=created/pending, 2=processing, 3=completed, 4=cancelled)
	// example: 1
	// required: true
	Status int  `json:"status" example:"1" format:"int32" minimum:"1"`
	// Array of the order line items.
	// example: [{"product_id":"6204037c-30e6-408b-8aaa-dd8219860b4b","quantity":2,"unit_price":50.00}]
	// required: true
	OrderDetails []OrderDetail `json:"order_details"`
}

// Convert order details to database-compatible array
func (v OrderDetail) Value() (driver.Value, error) {
    return []byte(fmt.Sprintf("(%s,%d,%f)", v.ProductId, v.Quantity, v.UnitPrice)), nil
}