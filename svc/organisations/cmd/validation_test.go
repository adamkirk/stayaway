// Test the validation services altogether
// It's an integration tests so i figure this is the best place for it...
// TODO: Add more tests
package main_test

import (
	"testing"

	"github.com/adamkirk-stayaway/organisations/internal/api"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/stretchr/testify/assert"
)

type CreatePersonsRequestAddress struct {
	Line1    *string `json:"line_1" validationmap:"Address.Street"`
	Postcode *string `json:"post_code" validationmap:"Address.Postcode"`
}

type CreatePersonRequest struct {
	Name    *string                     `json:"name" validationmap:"FullName"`
	Email   *string                     `json:"email" validationmap:"EmailAddress"`
	Address CreatePersonsRequestAddress `json:"address" validationmap:"Address"`
}

type Address struct {
	Street   *string `validate:"nonnil"`
	Postcode *string `validate:"nonnil"`
}

type Person struct {
	FullName     *string  `validate:"nonnil"`
	EmailAddress *string  `validate:"nonnil"`
	Address      *Address `validate:"nonnil"`
}

func ptr[T any](value T) *T {
	return &value
}

func TestValidationMapper(t *testing.T) {
	tests := []struct{
		name string
		p Person
		expect validation.ValidationError
	}{
		{
			name: "all empty (address included)",
			p: Person{
				Address: &Address{},
			},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address.line_1", 
						Errors: []string{"zero value"},
					}, 
					{
						Key: "address.post_code",
						Errors: []string{"zero value"},
					}, 
					{
						Key: "email", 
						Errors: []string{"zero value"},
					}, 
					{
						Key: "name",
						Errors: []string{"zero value"},
					},
				},
			},
		},
		{
			name: "all empty (address not included)",
			p: Person{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address",
						Errors: []string{"zero value"},
					}, 
					{
						Key: "email", 
						Errors: []string{"zero value"},
					}, 
					{
						Key: "name",
						Errors: []string{"zero value"},
					},
				},
			},
		},
		{
			name: "all empty (address not included)",
			p: Person{
				FullName: ptr("some name"),
				EmailAddress: ptr("someone@example.com"),
			},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address",
						Errors: []string{"zero value"},
					}, 
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func (tt *testing.T) {
			vm := api.NewValidationMapper()
			v := validation.NewValidator()
		
			err := v.Validate(test.p)
		
			vErr, _ := err.(validation.ValidationError)
		
			mapped := vm.Map(vErr, CreatePersonRequest{})
		
			assert.Equal(t,
				test.expect,
				mapped)
		})
	}
}
