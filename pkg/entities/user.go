// entities/user.go
package entities

import "time"

type User struct {
	Id        string    `json:"id"`
	Name 	  string    `json:"name"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"passhash" validate:"required"`
	Mobile    string    `json:"mobile"`
	Email  	  string    `json:"email" binding:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	Name 	  string    `json:"name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Mobile    string    `json:"mobile" validate:"required"`
	Email  	  string    `json:"email" binding:"email"`
}