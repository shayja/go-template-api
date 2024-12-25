// internal/entities/common.go
package entities

// IdRequest represents a request to get any entity by its ID.
type IdRequest struct {
	Id string `uri:"id" binding:"required"`
}