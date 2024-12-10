package errors

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrProductNotFound = errors.New("product not found")
)
