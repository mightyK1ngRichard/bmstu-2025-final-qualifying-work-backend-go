package errs

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidUUIDFormat = errors.New("invalid UUID format")
)
