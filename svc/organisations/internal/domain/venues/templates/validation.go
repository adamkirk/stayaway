package templates

import (
	"fmt"
	"strings"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
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
			Rule: "accommodationtype",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"accommodationtype must be one of: '%s'",
					strings.Join(common.AllAccommodationConfigTypes(), "', '"),
				)
				return ut.Add("accommodationtype", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("accommodationtype")

				return t
			},
		},
		{
			Rule: "accommodationtype_sortfield",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"must be one of: '%s'",
					string(SortByName),
				)
				return ut.Add("accommodationtype_sortfield", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("accommodationtype_sortfield")

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

				return common.AccommodationConfigType(val).IsValid()
			},
		},
		{
			Rule: "accommodationtype_sortfield",
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
