package organisations

type Validator interface {
	Validate(any) error
}