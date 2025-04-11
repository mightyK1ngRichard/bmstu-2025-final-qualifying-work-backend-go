package errs

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrInvalidUUIDFormat      = errors.New("invalid UUID format")
	ErrUnexpectedSignInMethod = errors.New("unexpected signing method")
	ErrInvalidTokenOrClaims   = errors.New("invalid token or claims")
	ErrParsingToken           = errors.New("error parsing token")
)
