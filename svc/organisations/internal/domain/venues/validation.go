package venues

import (
	"fmt"
	"strings"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationExtension struct{}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{
		{
			Rule: "venuetype",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"venuetype must be one of: '%s'",
					strings.Join(AllTypes(), "', '"),
				)
				return ut.Add("venuetype", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("venuetype")

				return t
			},
		},
		{
			Rule: "venues_sortfield",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"must be one of: '%s', '%s'",
					string(SortByName),
					string(SortBySlug),
				)
				return ut.Add("venues_sortfield", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("venues_sortfield")

				return t
			},
		},
	}
}

func (ve *ValidationExtension) Rules() []validation.CustomRule {
	return []validation.CustomRule{
		{
			Rule: "venuetype",
			Handler: func(fl validator.FieldLevel) bool {
				val := fl.Field().String()

				return Type(val).IsValid()
			},
		},
		{
			Rule: "venues_sortfield",
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
