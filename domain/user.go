package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserRole          = "USER"
	AdminRole         = "ADMIN"
	MaxFullNameLength = 100
	MinFullNameLength = 3
)

type User struct {
	Id              uuid.UUID `json:"id" validate:"required"`
	FullName        string    `json:"fullName" validate:"required"`
	Role            string    `json:"role" validate:"required"`
	Password        string    `json:"password" validate:"required"`
	IsEmailVerified bool      `json:"isEmailVerified" validate:"required"`
	CreatedAt       time.Time `json:"createdAt" validate:"required"`
	Email           string    `json:"email" validate:"required"`
}

func NewUser(fullName, password, email string) (*User, error) {

	if err := ValidateFullName(fullName); err != nil {
		return nil, err
	}

	if IsPasswordTooShort(password) {
		return nil, ErrPasswordTooShort
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, ErrInternalServer
	}

	return &User{
		Id:       uuid.New(),
		FullName: fullName,
		Role:     UserRole,
		Password: hashedPassword,
		Email:    email,
	}, nil
}

func (u *User) ValidatePassword(password string) error {

	if IsPasswordTooShort(password) {
		return ErrPasswordTooShort
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}

		return ErrInternalServer
	}
	return nil
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return ErrInternalServer
	}
	u.Password = hashedPassword
	return nil
}

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrInternalServer
	}
	return string(hashedPassword), nil
}

func ValidateFullName(fullName string) error {

	if len(fullName) == 0 {
		return ErrEmptyFullName
	}

	if len(fullName) < MinFullNameLength {
		return ErrTooShortFullName
	}

	return nil
}

func IsPasswordTooShort(password string) bool {
	return len(password) < 8
}
