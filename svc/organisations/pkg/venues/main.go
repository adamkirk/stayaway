package venues

type Validator interface {
	Validate(any) error
}