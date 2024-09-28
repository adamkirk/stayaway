package model

import "fmt"

type ErrNotFound struct {
	ResourceName string
	ID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("resource of type '%s' with ID '%s' not found", e.ResourceName, e.ID)
}


type  ErrInvalidSortBy struct {
	Chosen string
}

func (e ErrInvalidSortBy) Error() string {
	return fmt.Sprintf("invalid sorting field chosen: %s", e.Chosen)
}

type  ErrInvalidSortDir struct {
	Chosen string
}

func (e ErrInvalidSortDir) Error() string {
	return fmt.Sprintf("invalid sorting direction chosen: %s", e.Chosen)
}

type ErrGroup struct {
	Errors []error
}

func (e ErrGroup) Error() string {
	return fmt.Sprintf("%d errors occurred", len(e.Errors))
}

func (e ErrGroup) All() []error {
	return e.Errors
}