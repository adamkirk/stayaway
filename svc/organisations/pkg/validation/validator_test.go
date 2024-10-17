package validation_test

import (
	"testing"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pointTo[T any](value T) *T {
	return &value
}

type subject struct {
	RequiredString *string `validate:"required"`
}

type subjectWithCustomRule struct {
	RequiredString *string `validate:"required,enforceblah"`
}

type extensionProvider struct {
	translations []validation.Translation
	rules []validation.CustomRule
}

func (ep *extensionProvider) Translations() []validation.Translation {
	return ep.translations
}

func (ep *extensionProvider) Rules() []validation.CustomRule {
	return ep.rules
}

func TestValidationError(t *testing.T) {
	err := validation.ValidationError{}

	assert.Equal(t, "invalid data", err.Error())
}

func TestValidator(t *testing.T) {
	tests := []struct{
		name string
		in any
		expectNilError bool
		genericErrorType error
		validationError validation.ValidationError
		extensions []validation.Extension
	}{
		{
			name: "test basic validation - struct name is removed from key",
			in: subject{},
			expectNilError: false,
			validationError: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "RequiredString",
						Errors: []string{
							// Default validation message from library
							// Annoying but need to deal with this if it changes in the library
							// Really just ensuring that the library is being used
							"RequiredString is a required field",
						},
					},
				},
			},
		},
		{
			name: "valid data",
			in: subject{
				RequiredString: pointTo("blah"),
			},
			expectNilError: true,
		},
		{
			name: "error from validator",
			in: 1234, // can't pass int to validate struct
			expectNilError: false,
			genericErrorType: &validator.InvalidValidationError{},
		},
		{
			name: "test custom rule is available",
			in: subjectWithCustomRule{
				RequiredString: pointTo("meh"),
			},
			expectNilError: false,
			validationError: validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key: "RequiredString",
						Errors: []string{
							// Default validation message from library
							// Annoying but need to deal with this if it changes in the library
							// Really just ensuring that the library is being used
							"should be 'blah'",
						},
					},
				},
			},
			extensions: []validation.Extension{
				&extensionProvider{
					rules: []validation.CustomRule{
						{
							Rule: "enforceblah",
							Handler: func(fl validator.FieldLevel) bool {
								return fl.Field().String() == "blah"
							},
						},
					},
					translations: []validation.Translation{
						{
							Rule: "enforceblah",
							RegisterFunc: func(trans ut.Translator) error {
								return trans.Add("enforceblah", "should be 'blah'", true)
							},
							TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
								t, _ := ut.T("enforceblah")
				
								return t
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			v := validation.NewValidator(test.extensions...)

			err := v.Validate(test.in)

			if test.expectNilError {
				assert.Nil(tt, err)
				return
			}

			if test.genericErrorType != nil {
				assert.IsType(tt, test.genericErrorType, err)
				return
			}

			require.IsType(tt, validation.ValidationError{}, err)

			validationErr := err.(validation.ValidationError)

			assert.ElementsMatch(tt, test.validationError.Errs, validationErr.Errs)
		})
	}
}