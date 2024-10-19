package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
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
				if cmd.Description == nil || cmd.MinOccupancy == nil || cmd.Name == nil || cmd.Type == nil {
					sl.ReportError(cmd.Name, "Name", "Name", "templateid_or_config", "")
					sl.ReportError(cmd.Description, "Description", "Description", "templateid_or_config", "")
					sl.ReportError(cmd.MinOccupancy, "MinOccupancy", "MinOccupancy", "templateid_or_config", "")
					sl.ReportError(cmd.Type, "Type", "Type", "templateid_or_config", "")
					sl.ReportError(cmd.VenueTemplateID, "VenueTemplateID", "VenueTemplateID", "templateid_or_config", "")
				}
			},
		},
	}
}

func (ve *ValidationExtension) Translations() []validation.Translation {
	return []validation.Translation{}
}

func (ve *ValidationExtension) Rules() []validation.CustomRule {
	return []validation.CustomRule{}
}

func NewValidationExtension() *ValidationExtension {
	return &ValidationExtension{}
}
