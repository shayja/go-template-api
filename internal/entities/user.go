// internal/entities/user.go
package entities

import "time"

type User struct {
	Id        	string    `json:"id"`
	FirstName 	string    `json:"first_name"`
	LastName 	string    `json:"last_name"`
	Username  	string    `json:"username" validate:"required"`
	Password  	string    `json:"password" validate:"required"`
	Mobile    	string    `json:"mobile"`
	Email  	  	string    `json:"email" binding:"email"`
	OtpTypes   	[]int     `json:"otp_types"`
	Verified   	bool 	  `json:"verified"`
	VerifiedAt	*time.Time `json:"verified_at"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}

type UserRequest struct {
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Mobile    string    `json:"mobile" validate:"required"`
	Email  	  string    `json:"email" binding:"email"`
}