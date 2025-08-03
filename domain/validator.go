package domain

type Validator interface {
	Validate(data any) error //it also prints exact error message
}
