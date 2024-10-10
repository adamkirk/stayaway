package venues

import (
	"fmt"
	"strings"

	"github.com/adamkirk-stayaway/organisations/internal/model"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationExtension struct {}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{
		{
			Rule: "venuetype",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"venuetype must be one of: '%s'",
					strings.Join(model.AllVenueTypes(), "', '"),
				)
				return ut.Add("venuetype", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("venuetype")
				
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

				return model.VenueType(val).IsValid()
			},
		},
	}
}

func NewValidationExtension() *ValidationExtension {
	return &ValidationExtension{}
}