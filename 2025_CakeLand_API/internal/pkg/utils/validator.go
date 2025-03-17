package utils

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

var (
	passwordRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	emailRegexp    = regexp.MustCompile(`^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,4}$`)
	nameRegexp     = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

// ValidateEmail Функция валидации почты
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "email обязателен")
	} else if !emailRegexp.MatchString(email) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	return nil
}

// ValidatePassword Функция для проверки валидности пароля
func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return status.Error(codes.InvalidArgument, "password обязателен")
	} else if len(password) < 8 {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	} else if !passwordRegexp.MatchString(password) {
		return status.Error(codes.InvalidArgument, "password must contain at least one letter and one number")
	}

	return nil
}

// ValidateName Функция валидации имени пользователя
func (v *Validator) ValidateName(name string) error {
	if len(name) < 2 || len(name) > 50 {
		return errors.New("name must be between 2 and 50 characters long")
	} else if !nameRegexp.MatchString(name) {
		return errors.New("name can only contain letters, spaces, and '-'")
	}

	return nil
}
