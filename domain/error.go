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
	ErrInvalidCurrency      = errors.New("currency mismatch")
	ErrInvalidAmount        = errors.New("invalid donation amount")
	ErrInvalidDonationType  = errors.New("invalid donation type")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrInvalidFormat        = errors.New("invalid data format")
	ErrPasswordTooShort     = errors.New("password must be at least 8 characters long")
	ErrInvalidRequest       = errors.New("geçersiz istek")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrInvalidSession       = errors.New("invalid or expired session")
	ErrInvalidCredentials   = errors.New("invalid credentials")

	ErrForbidden              = errors.New("forbidden resource")
	ErrInsufficientPermission = errors.New("insufficient permissions")

	ErrBucketNotFound   = errors.New("bucket not found")
	ErrNoRows           = errors.New("no rows in result set")
	ErrEmailNotFound    = errors.New("email not found")
	ErrDonationNotFound = errors.New("donation not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrResourceNotFound = errors.New("resource not found")

	ErrNotEnoughLanguages     = errors.New("Tüm dillerde veri gönderilmesi gerekiyor")
	ErrNotEnoughCurrencies    = errors.New("Tüm dövizlerde veri gönderilmesi gerekiyor")
	ErrEmailAlreadyVerified   = errors.New("email already verified")
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrDuplicatedName         = errors.New("duplicated name")
	ErrDuplicateEntity        = errors.New("entity already exists")
	ErrConcurrentModification = errors.New("resource was modified by another request")

	ErrDatabaseOperation     = errors.New("database operation failed")
	ErrInternalServer        = errors.New("internal server error")
	ErrExternalServiceFailed = errors.New("external service failed")
)
