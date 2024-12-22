// internal/entities/common.go
package entities

type IdRequest struct {
	Id string `uri:"id" binding:"required"`
}