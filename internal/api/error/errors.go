package errs

import "errors"

var (
	ErrInternal   = errors.New("internal error")
	ErrValidation = errors.New("validation error")
	ErrBadRequest = errors.New("bad request")

	ErrNotFound = errors.New("task was not found")
)
