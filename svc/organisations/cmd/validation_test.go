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

type CreatePersonRequestAddress struct {
	Line1    *string `json:"line_1" validationmap:"Address.Street"`
	Postcode *string `json:"post_code" validationmap:"Address.Postcode"`
}

type CreatePersonRequest struct {
	Name    *string                     `json:"name" validationmap:"FullName"`
	Email   *string                     `json:"email" validationmap:"EmailAddress"`
	Address CreatePersonRequestAddress `json:"address" validationmap:"Address"`
}

type Address struct {
	Street   *string `validate:"required,min=4,alphanum"`
	Postcode *string `validate:"required,min=6,alphanum"`
}

type Person struct {
	FullName     *string  `validate:"required"`
	EmailAddress *string  `validate:"required"`
	Address      *Address `validate:"required"`
}

type FlattenedPersonRequestAddress struct {
	Line1    *string `json:"line_1" validationmap:"Street"`
	Postcode *string `json:"post_code" validationmap:"Postcode"`
}

type FlattenedPersonRequest struct {
	Name    *string                     `json:"name" validationmap:"FullName"`
	Email   *string                     `json:"email" validationmap:"EmailAddress"`
	Address FlattenedPersonRequestAddress `json:"address"`
}

type FlattenedPersonDTO struct {
	FullName     *string  `validate:"required"`
	EmailAddress *string  `validate:"required"`
	Street   *string `validate:"required"`
	Postcode *string `validate:"required"`
}

func ptr[T any](value T) *T {
	return &value
}

func TestValidationMapper(t *testing.T) {
	tests := []struct{
		name string
		p any
		req any
		expect validation.ValidationError
	}{
		{
			name: "all empty (address included)",
			p: Person{
				Address: &Address{},
			},
			req: CreatePersonRequest{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address.line_1", 
						Errors: []string{"is required"},
					}, 
					{
						Key: "address.post_code",
						Errors: []string{"is required"},
					}, 
					{
						Key: "email", 
						Errors: []string{"is required"},
					}, 
					{
						Key: "name",
						Errors: []string{"is required"},
					},
				},
			},
		},
		{
			name: "all empty (address not included)",
			p: Person{},
			req: CreatePersonRequest{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address",
						Errors: []string{"is required"},
					}, 
					{
						Key: "email", 
						Errors: []string{"is required"},
					}, 
					{
						Key: "name",
						Errors: []string{"is required"},
					},
				},
			},
		},
		{
			name: "address not included",
			p: Person{
				FullName: ptr("some name"),
				EmailAddress: ptr("someone@example.com"),
			},
			req: CreatePersonRequest{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address",
						Errors: []string{"is required"},
					}, 
				},
			},
		},
		{
			// Really this just highlights how the underlying library works. It 
			// will only generate one error at a time even if there are multiple
			// violations. Annoying from a user perspective, but it'll do for now.
			name: "multiple rule violations for single field",
			p: Person{
				FullName: ptr("some name"),
				EmailAddress: ptr("someone@example.com"),
				Address: &Address{
					Postcode: ptr("!fh"), // not long enough and only alphanum
					Street: ptr("!fh"), // not long enough and only alphanum
				},
			},
			req: CreatePersonRequest{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address.line_1",
						Errors: []string{"must be more than 4 characters long"},
					}, 
					{
						Key: "address.post_code",
						Errors: []string{"must be more than 6 characters long"},
					}, 
				},
			},
		},
		{
			name: "flattened_dto_to_nested_request",
			p: FlattenedPersonDTO{
				FullName: ptr("some name"),
				EmailAddress: ptr("someone@example.com"),
				Postcode: ptr("!fh"), // not long enough and only alphanum
			},
			req: FlattenedPersonRequest{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address.line_1",
						Errors: []string{"is required"},
					}, 
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func (tt *testing.T) {
			vm := api.NewValidationMapper()
			v := validation.NewValidator([]validation.Extension{})
		
			err := v.Validate(test.p)
		
			vErr, _ := err.(validation.ValidationError)
		
			mapped := vm.Map(vErr, test.req)
		
			assert.Equal(t,
				test.expect,
				mapped)
		})
	}
}
