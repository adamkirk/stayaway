package validation_test

import (
	"reflect"
	"testing"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/stretchr/testify/assert"
)

type CreatePersonRequestAddress struct {
	Line1    *string `json:"line_1" validationmap:"Address.Street"`
	Postcode *string `json:"post_code" validationmap:"Address.Postcode"`
}

type CreatePersonRequest struct {
	Name    *string                    `json:"name" validationmap:"FullName"`
	Email   *string                    `json:"email" validationmap:"EmailAddress"`
	Address CreatePersonRequestAddress `json:"address" validationmap:"Address"`
	NoJsonTag *string `validationmap:"Blah"`
}

func ptr[T any](value T) *T {
	return &value
}

func TestValidationMapper(t *testing.T) {
	tests := []struct {
		name      string
		in        validation.ValidationError
		opts []validation.ValidationMapperOpt
		mapTo any
		expectErr error
		expect    validation.ValidationError
	}{
		{
			name:      "default behaviour (json tags)",
			mapTo: CreatePersonRequest{},
			opts: []validation.ValidationMapperOpt{},
			in:        validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "Address.Street",
						Errors: []string{
							"street error 1",
							"street error 2",
						},
					},
					{
						Key: "EmailAddress",
						Errors: []string{
							"email error",
						},
					},
				},
			},
			expectErr: nil,
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "address.line_1",
						Errors: []string{
							"street error 1",
							"street error 2",
						},
					},
					{
						Key: "email",
						Errors: []string{
							"email error",
						},
					},
				},
			},
		},
		{
			name:      "default json tag finder defaults to field name",
			mapTo: CreatePersonRequest{},
			opts: []validation.ValidationMapperOpt{},
			in:        validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "Blah",
						Errors: []string{
							"blah error",
						},
					},
				},
			},
			expectErr: nil,
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "NoJsonTag",
						Errors: []string{
							"blah error",
						},
					},
				},
			},
		},
		{
			name:      "with a field that has no validationmap tag",
			mapTo: CreatePersonRequest{},
			opts: []validation.ValidationMapperOpt{},
			in:        validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "Meh",
						Errors: []string{
							"blah error",
						},
					},
				},
			},
			expectErr: nil,
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "Meh",
						Errors: []string{
							"blah error",
						},
					},
				},
			},
		},
		{
			name:      "with a custom tag finder",
			mapTo: CreatePersonRequest{},
			opts: []validation.ValidationMapperOpt{
				validation.WithTagFinder(func(f reflect.StructField) string {
					// Will make every field blah
					return f.Name
				}),
			},
			in:        validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "Address.Street",
						Errors: []string{
							"street error 1",
							"street error 2",
						},
					},
					{
						Key: "EmailAddress",
						Errors: []string{
							"email error",
						},
					},
				},
			},
			expectErr: nil,
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "Address.Line1",
						Errors: []string{
							"street error 1",
							"street error 2",
						},
					},
					{
						Key: "Email",
						Errors: []string{
							"email error",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			vm := validation.NewValidationMapper(test.opts...)
			
			res := vm.Map(test.in, test.mapTo)

			assert.Equal(tt, test.expect, res)
		})
	}
}

// func TestBuildValidationMap(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		in        any
// 		expectErr error
// 		expect    StructMapMeta
// 	}{
// 		{
// 			name:      "nested props",
// 			in:        CreatePersonRequest{},
// 			expectErr: nil,
// 			expect: StructMapMeta{
// 				"address":           "Address",
// 				"name":              "FullName",
// 				"email":             "EmailAddress",
// 				"address.line_1":    "Address.Street",
// 				"address.post_code": "Address.Postcode",
// 			},
// 		},

// 		{
// 			name:      "flattened target",
// 			in:        FlattenedPersonRequest{},
// 			expectErr: nil,
// 			expect: StructMapMeta{
// 				"address.line_1":    "Street",
// 				"address.post_code": "Postcode",
// 				"email":             "EmailAddress",
// 				"name":              "FullName",
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(tt *testing.T) {

// 			res := mapJsonFieldsToMapTags(reflect.TypeOf(test.in))

// 			// assert.Equal(tt, test.expectErr, err)

// 			assert.Equal(tt, test.expect, res)
// 		})
// 	}
// }
