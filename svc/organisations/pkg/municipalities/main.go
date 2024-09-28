package municipalities

type Validator interface {
	Validate(any) error
}