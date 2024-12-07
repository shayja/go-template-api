// entities/authentication_input.go
package entities

type AuthenticationInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}