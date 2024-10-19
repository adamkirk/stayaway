package organisations

import (
	"fmt"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationExtension struct{}

func (ve *ValidationExtension) StructValidators() []validation.StructValidator {
	return []validation.StructValidator{}
}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{
		{
			Rule: "organisations_sortfield",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"must be one of: '%s', '%s'",
					string(SortByName),
					string(SortBySlug),
				)
				return ut.Add("organisations_sortfield", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("organisations_sortfield")

				return t
			},
		},
	}
}

func (ve *ValidationExtension) Rules() []validation.CustomRule {
	return []validation.CustomRule{
		{
			Rule: "organisations_sortfield",
			Handler: func(fl validator.FieldLevel) bool {
				val := fl.Field().String()

				return val == string(SortByName) || val == string(SortBySlug)
			},
		},
	}
}

func NewValidationExtension() *ValidationExtension {
	return &ValidationExtension{}
}
