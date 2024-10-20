package accommodations

import (
	"fmt"

	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationExtension struct{}

func (ve *ValidationExtension) StructValidators() []validation.StructValidator {
	return []validation.StructValidator{
		{
			Struct: CreateCommand{},
			Validator: func(sl validator.StructLevel) {
				cmd := sl.Current().Interface().(CreateCommand)

				// Because the min is optional, we can't just rely on the getfield
				// rule as it will fail when min is nil and max isn't.
				// TODO: Add a translation for this rule
				if cmd.MinOccupancy != nil && cmd.MaxOccupancy != nil && *cmd.MaxOccupancy < *cmd.MinOccupancy {
					sl.ReportError(cmd.MaxOccupancy, "MaxOccupancy", "MaxOccupancy", "max_occupancy_less_than_min", "")
				}

				if cmd.VenueTemplateID != nil {
					return 
				}

				// So if the venue template id is not set, then all of these fields must be
				// TODO: Add a translation for this rule
				if cmd.Description == nil || cmd.MinOccupancy == nil || cmd.MaxOccupancy == nil || cmd.Name == nil || cmd.Type == nil {
					sl.ReportError(cmd.Name, "Name", "Name", "templateid_or_config", "")
					sl.ReportError(cmd.Description, "Description", "Description", "templateid_or_config", "")
					sl.ReportError(cmd.MinOccupancy, "MinOccupancy", "MinOccupancy", "templateid_or_config", "")
					sl.ReportError(cmd.MaxOccupancy, "MaxOccupancy", "MaxOccupancy", "templateid_or_config", "")
					sl.ReportError(cmd.Type, "Type", "Type", "templateid_or_config", "")
					sl.ReportError(cmd.VenueTemplateID, "VenueTemplateID", "VenueTemplateID", "templateid_or_config", "")
				}
			},
		},
		{
			Struct: UpdateCommand{},
			Validator: func(sl validator.StructLevel) {
				cmd := sl.Current().Interface().(UpdateCommand)

				anyConfigFieldIsNull := cmd.NullifyDescription || cmd.NullifyMinOccupancy || cmd.NullifyName || cmd.NullifyType || cmd.NullifyMaxOccupancy

				// If the template is being nullified, then all of these values 
				// must be set. If the template is omitted or set to a value, then
				// none of these values need to be present.
				// There is a later validation in the update method that checks
				// the merged configs is valid before persisting.
				if anyConfigFieldIsNull && cmd.NullifyVenueTemplateID {
					sl.ReportError(cmd.VenueTemplateID, "VenueTemplateID", "VenueTemplateID", "templateid_or_config", "")
				}
			},
		},
	}
}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{
		{
			Rule: "venueaccommodation_sortfield",
			RegisterFunc: func(ut ut.Translator) error {
				msg := fmt.Sprintf(
					"must be one of: '%s'",
					string(SortByReference),
				)
				return ut.Add("venueaccommodation_sortfield", msg, true)
			},
			TranslateFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("venueaccommodation_sortfield")

				return t
			},
		},
	}
}

func (ve *ValidationExtension) Rules() []validation.CustomRule {
	return []validation.CustomRule{
		{
			Rule: "venueaccommodation_sortfield",
			Handler: func(fl validator.FieldLevel) bool {
				val := fl.Field().String()

				return val == string(SortByReference)
			},
		},
	}
}

func NewValidationExtension() *ValidationExtension {
	return &ValidationExtension{}
}
