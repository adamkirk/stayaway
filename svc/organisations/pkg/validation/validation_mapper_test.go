// Provides a few basic tests, could do with a lot more examples, as this isn't
// particularly thorough of how much this can do.
package validation_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/adamkirk-stayaway/organisations/pkg/validation/test/mocks"
	"github.com/stretchr/testify/assert"
)

type CreatePersonRequestAddress struct {
	Line1    *string `json:"line_1" validationmap:"Address.Street"`
	Postcode *string `json:"post_code" validationmap:"Address.Postcode"`
}

type CreatePersonRequest struct {
	Name      *string                    `json:"name" query:"name_query" validationmap:"FullName"`
	Email     *string                    `json:"email" validationmap:"EmailAddress"`
	Address   CreatePersonRequestAddress `json:"address" validationmap:"Address"`
	NoJsonTag *string                    `validationmap:"Blah"`
}

func ptr[T any](value T) *T {
	return &value
}

// Just exemplary of how you might customise the behaviour
func getQueryThenJsonTagNameForField(f reflect.StructField) string {
	if queryTag := f.Tag.Get("query"); queryTag != "" {
		return strings.Split(queryTag, ",")[0]
	}

	if jsonTag := f.Tag.Get("json"); jsonTag != "" {
		// Sensible default
		return strings.Split(jsonTag, ",")[0]
	}

	return f.Name
}

func TestValidationMapper(t *testing.T) {
	tests := []struct {
		name      string
		in        validation.ValidationError
		opts      []validation.ValidationMapperOpt
		mapTo     any
		expectErr error
		expect    validation.ValidationError
	}{
		{
			name:  "default behaviour (json tags)",
			mapTo: CreatePersonRequest{},
			opts:  []validation.ValidationMapperOpt{},
			in: validation.ValidationError{
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
			name:  "default json tag finder defaults to field name",
			mapTo: CreatePersonRequest{},
			opts:  []validation.ValidationMapperOpt{},
			in: validation.ValidationError{
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
			name:  "with a field that has no validationmap tag",
			mapTo: CreatePersonRequest{},
			opts:  []validation.ValidationMapperOpt{},
			in: validation.ValidationError{
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
			name:  "with a custom tag finder",
			mapTo: CreatePersonRequest{},
			opts: []validation.ValidationMapperOpt{
				validation.WithTagFinder(func(f reflect.StructField) string {
					// Will make every field blah
					return f.Name
				}),
			},
			in: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "FullName",
						Errors: []string{
							"name error",
						},
					},
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
						Key: "Name",
						Errors: []string{
							"name error",
						},
					},
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
			name:  "test custom query tag finder",
			mapTo: CreatePersonRequest{},
			opts: []validation.ValidationMapperOpt{
				validation.WithTagFinder(getQueryThenJsonTagNameForField),
			},
			in: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "FullName",
						Errors: []string{
							"name error",
						},
					},
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
					{
						Key: "name_query",
						Errors: []string{
							"name error",
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

			// Random ordering makes things messy
			assert.ElementsMatch(tt, test.expect.Errs, res.Errs)
		})
	}
}

// Ensure it's logging anythig it should
func TestValidationMapperWithMockLogger(t *testing.T) {
	mapTo := CreatePersonRequest{}
	in := validation.ValidationError{
		Errs: []validation.FieldError{
			{
				Key: "Meh",
				Errors: []string{
					"blah error",
				},
			},
		},
	}
	expect := validation.ValidationError{
		Errs: []validation.FieldError{
			{
				Key: "Meh",
				Errors: []string{
					"blah error",
				},
			},
		},
	}

	logger := mocks.NewMockLogger(t)
	logger.EXPECT().Warn("did not find validationmap field error", "field", "Meh", "type", "CreatePersonRequest")

	vm := validation.NewValidationMapper(validation.WithLogger(logger))

	res := vm.Map(in, mapTo)

	assert.Equal(t, expect, res)
}

// Run a couple of different scenarios with same mapper instance
func TestValidationMapperMultipleRuns(t *testing.T) {
	mapTo := CreatePersonRequest{}
	in := validation.ValidationError{
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
	}
	expect := validation.ValidationError{
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
	}

	vm := validation.NewValidationMapper()

	res := vm.Map(in, mapTo)

	assert.Equal(t, expect, res)

	in = validation.ValidationError{
		Errs: []validation.FieldError{
			{
				Key: "Address.Postcode",
				Errors: []string{
					"postcode error",
				},
			},
		},
	}
	expect = validation.ValidationError{
		Errs: []validation.FieldError{
			{
				Key: "address.post_code",
				Errors: []string{
					"postcode error",
				},
			},
		},
	}

	res = vm.Map(in, mapTo)

	assert.Equal(t, expect, res)
}

// Ensure it's actually using the cache we provide
func TestValidationMapperWithMockCacher(t *testing.T) {
	mapTo := CreatePersonRequest{}
	in := validation.ValidationError{
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
	}
	expect := validation.ValidationError{
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
	}

	meta := validation.StructMapMeta{
		"address.line_1":    "Address.Street",
		"address.post_code": "Address.Postcode",
		"email":             "EmailAddress",
		"name":              "FullName",
		"NoJsonTag":         "Blah",
		"address":           "Address",
	}

	cacher := mocks.NewMockCacher(t)
	cacher.EXPECT().Get("github.com/adamkirk-stayaway/organisations/pkg/validation_test/CreatePersonRequest").Return(validation.StructMapMeta{}, false)
	cacher.EXPECT().Add("github.com/adamkirk-stayaway/organisations/pkg/validation_test/CreatePersonRequest", meta)
	cacher.EXPECT().Get("github.com/adamkirk-stayaway/organisations/pkg/validation_test/CreatePersonRequest").Return(meta, true)

	vm := validation.NewValidationMapper(validation.WithCacher(cacher))

	res := vm.Map(in, mapTo)

	assert.Equal(t, expect, res)

	res = vm.Map(in, mapTo)

	assert.Equal(t, expect, res)
}
