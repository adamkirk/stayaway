package model

import "fmt"

type ErrNotFound struct {
	ResourceName string
	ID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("resource of type '%s' with ID '%s' not found", e.ResourceName, e.ID)
}
