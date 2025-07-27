package domain

type Validator interface {
	Validate(data any) error
}
