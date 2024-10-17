package municipalities

import (
	"fmt"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationExtension struct{}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{
		{
			Rule: "municipalities_sortfield",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"must be one of: '%s'",
					string(SortByName),
				)
				return ut.Add("municipalities_sortfield", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("municipalities_sortfield")

				return t
			},
		},
	}
}

func (ve *ValidationExtension) Rules() []validation.CustomRule {
	return []validation.CustomRule{
		{
			Rule: "municipalities_sortfield",
			Handler: func(fl validator.FieldLevel) bool {
				val := fl.Field().String()

				return val == string(SortByName)
			},
		},
	}
}

func NewValidationExtension() *ValidationExtension {
	return &ValidationExtension{}
}
