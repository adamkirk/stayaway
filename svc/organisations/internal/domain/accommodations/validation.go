package accommodations

import (
	"fmt"
	"strings"

	"github.com/adamkirk-stayaway/organisations/internal/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationExtension struct{}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{
		{
			Rule: "accommodationtype",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"accommodationtype must be one of: '%s'",
					strings.Join(AllTypes(), "', '"),
				)
				return ut.Add("accommodationtype", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("accommodationtype")

				return t
			},
		},
	}
}

func (ve *ValidationExtension) Rules() []validation.CustomRule {
	return []validation.CustomRule{
		{
			Rule: "accommodationtype",
			Handler: func(fl validator.FieldLevel) bool {
				val := fl.Field().String()

				return Type(val).IsValid()
			},
		},
	}
}

func NewValidationExtension() *ValidationExtension {
	return &ValidationExtension{}
}
