package validation_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type exampleDtoNested struct {
	SubField *string `validate:"required"`
}

type exampleDto struct {
	ShortField *string `validate:"required"`
	OtherField *string `validate:"required"`
	Nested exampleDtoNested 
}

type mapTarget struct {
	ShortFieldAlt *string `json:"field_a" validationmap:"ShortField"`
	OtherFieldAlt *string `json:"field_b" validationmap:"OtherField"`
	NestedSubAlt *string `json:"sub" validationmap:"Nested.SubField"`
}

type targetSubField struct {
	Short *string `json:"short" validationmap:"ShortField"`
	Other *string `json:"other" validationmap:"OtherField"`
}

type targetWithNesting struct {
	Sub targetSubField `json:"sub"`
	Top *string `json:"top" validationmap:"Nested.SubField"`
}

type mapTargetWithDifferentTag struct {
	ShortFieldAlt *string `query:"query_field_a" json:"json_field_a" validationmap:"ShortField"`
	OtherFieldAlt *string `random:"random_field_b" json:"json_field_b" validationmap:"OtherField"`
	NestedSubAlt *string `validationmap:"Nested.SubField"` // has no tag
}

// Just exemplary of how you might customise the behaviour
func tagFinder(f reflect.StructField) string {
	if tag := f.Tag.Get("query"); tag != "" {
		return strings.Split(tag, ",")[0]
	}

	if tag := f.Tag.Get("random"); tag != "" {
		// Sensible default
		return strings.Split(tag, ",")[0]
	}

	if tag := f.Tag.Get("json"); tag != "" {
		// Sensible default
		return strings.Split(tag, ",")[0]
	}

	return f.Name
}

func TestCanMapValidatorResults(t *testing.T) {
	extensions := []validation.Extension{
		&extensionProvider{
			translations: []validation.Translation{
				{
					Rule: "required",
					RegisterFunc: func(trans ut.Translator) error {
						return trans.Add("required", "{0} is required", true)
					},
					TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
						
						t, _ := ut.T("required", fe.StructNamespace())
		
						return t
					},
				},
			},
		},
	}

	tests := []struct{
		name string
		in any
		target any
		expect validation.ValidationError
		opts []validation.ValidationMapperOpt
	}{
		{
			// Errors are mapped back to json tags in the target based on the 
			// validationmap tags.
			// exampleDto.ShortField -> field_a
			// exampleDto.OtherField -> field_b
			// exampleDto.Nested.SubField -> sub
			name: "test violations are mapped to json tags of target",
			target: mapTarget{},
			in: exampleDto{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "field_a",
						Errors: []string{
							"exampleDto.ShortField is required",
						},
					},
					{
						Key: "field_b",
						Errors: []string{
							"exampleDto.OtherField is required",
						},
					},
					{
						Key: "sub",
						Errors: []string{
							"exampleDto.Nested.SubField is required",
						},
					},
				},
			},
		},
		{
			// Reverses the nesting of the dto so that the fields from nested are
			// placed at top level of the target
			// exampleDto.ShortField -> sub.short
			// exampleDto.OtherField -> sub.other
			// exampleDto.Nested.SubField -> top
			name: "test violations are mapped to nested fields in target",
			target: targetWithNesting{},
			in: exampleDto{},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "sub.short",
						Errors: []string{
							"exampleDto.ShortField is required",
						},
					},
					{
						Key: "sub.other",
						Errors: []string{
							"exampleDto.OtherField is required",
						},
					},
					{
						Key: "top",
						Errors: []string{
							"exampleDto.Nested.SubField is required",
						},
					},
				},
			},
		},
		{
			// Reverses the nesting of the dto so that the fields from nested are
			// placed at top level of the target
			// exampleDto.ShortField -> query_field_a
			// exampleDto.OtherField -> random_field_b
			// exampleDto.Nested.SubField -> NestedSubAlt
			name: "test violations with custom tag finder",
			target: mapTargetWithDifferentTag{},
			in: exampleDto{},
			opts: []validation.ValidationMapperOpt{
				validation.WithTagFinder(tagFinder),
			},
			expect: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "query_field_a",
						Errors: []string{
							"exampleDto.ShortField is required",
						},
					},
					{
						Key: "random_field_b",
						Errors: []string{
							"exampleDto.OtherField is required",
						},
					},
					{
						Key: "NestedSubAlt",
						Errors: []string{
							"exampleDto.Nested.SubField is required",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			v := validation.NewValidator(extensions...)
			vm := validation.NewValidationMapper(test.opts...)

			err := v.Validate(test.in)

			require.IsType(tt, validation.ValidationError{}, err)

			validationErr := err.(validation.ValidationError)

			validationErr = vm.Map(validationErr, test.target)

			assert.ElementsMatch(tt, test.expect.Errs, validationErr.Errs)
		})
	}
}