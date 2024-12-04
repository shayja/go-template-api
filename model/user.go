package model

import "time"

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"passhash" validate:"required"`
	Mobile    string    `json:"mobile"`
	Name 	  string    `json:"name"`
	Email  	  string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}