package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/templates"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type CreateCommand struct {
	VenueID         *string `validate:"required"`
	VenueTemplateID *string 
	Reference            *string `validate:"required,min=3"`
	Name *string `validate:"omitnil,min=3" validationmap:"Name"`
	Type            *string `validate:"omitnil,accommodationtype" validationmap:"Type"`
	MinOccupancy    *int    `validate:"omitnil,min=1" validationmap:"MinOccupancy"`
	// 100 seems an appropriate max
	MaxOccupancy *int    `validate:"omitnil" validationmap:"MaxOccupancy"`
	Description  *string `validate:"omitnil,min=10" validationmap:"MinOccupancy"`
}

func (svc *Service) Create(cmd CreateCommand) (*Accommodation, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}


	var template *templates.VenueTemplate

	if cmd.VenueTemplateID != nil {
		template, err = svc.templatesRepo.Get(*cmd.VenueTemplateID, *cmd.VenueID)

		if err != nil {
			return nil, err
		}
	}

	accByName, err := svc.repo.ByReferenceAndVenueID(*cmd.Reference, *cmd.VenueID)

	if accByName != nil {
		return nil, validation.ValidationError{
			Errs: []validation.FieldError{
				{
					Key:    "Reference",
					Errors: []string{"must be unique"},
				},
			},
		}
	}

	if err != nil {
		if _, ok := err.(common.ErrNotFound); !ok {
			return nil, err
		}
	}


	a := &Accommodation{
		ID:              svc.idGen.Generate(),
		VenueID:         *cmd.VenueID,
		VenueTemplateID: cmd.VenueTemplateID,
		Reference: *cmd.Reference,
		Overrides: &common.AccommodationConfigOverrides{
			Name:         cmd.Name,
			MinOccupancy: cmd.MinOccupancy,
			MaxOccupancy: cmd.MaxOccupancy,
			Description:  cmd.Description,
			Type:         (*common.AccommodationConfigType)(cmd.Type),
		},

	}

	svc.decorateWithTemplateConfig(template, a)

	err = svc.validator.Validate(a.Config)

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			err = err.PrefixAll("incompatible with template: ")
			return nil, svc.validationMapper.Map(err, CreateCommand{})
		}
		return nil, err
	}
	

	a, err = svc.repo.Save(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}
