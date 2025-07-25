package domain

import "errors"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

const (
	CodeBadRequest          = 400
	CodeUnauthorized        = 401
	CodeForbidden           = 403
	CodeNotFound            = 404
	CodeConflict            = 409
	CodeUnprocessableEntity = 422
	CodeInternalServerError = 500
	CodeServiceUnavailable  = 503
)

var (
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrInvalidRequest     = errors.New("geçersiz istek")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoRows            = errors.New("no rows in result set")
	ErrEmailNotFound     = errors.New("email not found")
	ErrUserNotFound      = errors.New("user not found")

	ErrEmailAlreadyVerified = errors.New("email already verified")
	ErrEmailAlreadyExists   = errors.New("email already exists")

	ErrInternalServer = errors.New("internal server error")
)
