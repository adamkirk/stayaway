package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/templates"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type UpdateCommand struct {
	ID         *string `validate:"required"`
	VenueID         *string `validate:"required"`

	VenueTemplateID *string
	NullifyVenueTemplateID bool

	Reference            *string `validate:"omitnil,min=3"`

	Name *string `validate:"omitnil,min=3" validationmap:"Name"`
	NullifyName bool

	Type            *string `validate:"omitnil,accommodationtype" validationmap:"Type"`
	NullifyType bool

	MinOccupancy    *int    `validate:"omitnil,min=1" validationmap:"MinOccupancy"`
	NullifyMinOccupancy bool

	// 100 seems an appropriate max
	MaxOccupancy *int    `validate:"omitnil" validationmap:"MaxOccupancy"`
	NullifyMaxOccupancy bool

	Description  *string `validate:"omitnil,min=10" validationmap:"MinOccupancy"`
	NullifyDescription bool
}

func (svc *Service) Update(cmd UpdateCommand) (*Accommodation, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	acc, err := svc.repo.Get(*cmd.ID, *cmd.VenueID)

	if err != nil {
		return nil, err
	}

	if cmd.Reference != nil {
		accByName, err := svc.repo.ByReferenceAndVenueID(*cmd.Reference, *cmd.VenueID)
	
		if accByName != nil && accByName.ID != acc.ID {
			return nil, validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key:    "Reference",
						Errors: []string{"must be unique"},
					},
				},
			}
		} else if _, ok := err.(common.ErrNotFound); !ok {
			return nil, err
		}
	}

	var t *templates.VenueTemplate

	if cmd.NullifyVenueTemplateID {
		acc.VenueTemplateID = nil
	} else if cmd.VenueTemplateID != nil {
		t, err = svc.templatesRepo.Get(*cmd.VenueTemplateID, *cmd.VenueID)

		if err != nil {
			if _, ok := err.(common.ErrNotFound); !ok {
				return nil, err
			}
			return nil, err
		}

		acc.VenueTemplateID = &t.ID
	} else {
		// Need to set the template to use in decoration later
		t, err = svc.templatesRepo.Get(*acc.VenueTemplateID, *cmd.VenueID)

		if err != nil {
			return nil, err
		}
	}

	// The struct validator should ensure we can't nullify the template and any
	// of these properties. If the template was already nil, then the props will
	// still be nullified here, but the later validation of the config will ensure
	// we can't persist a bad config.
	if cmd.NullifyDescription {
		acc.Overrides.Description = nil
	} else if cmd.Description != nil {
		acc.Overrides.Description = cmd.Description
	}

	if cmd.NullifyName {
		acc.Overrides.Name = nil
	} else if cmd.Name != nil {
		acc.Overrides.Name = cmd.Name
	}

	if cmd.NullifyMinOccupancy {
		acc.Overrides.MinOccupancy = nil
	} else if cmd.MinOccupancy != nil {
		acc.Overrides.MinOccupancy = cmd.MinOccupancy
	}

	if cmd.NullifyMaxOccupancy {
		acc.Overrides.MaxOccupancy = nil
	} else if cmd.MaxOccupancy != nil {
		acc.Overrides.MaxOccupancy = cmd.MaxOccupancy
	}

	if cmd.NullifyType {
		acc.Overrides.Type = nil
	} else if cmd.Type != nil {
		acc.Overrides.Type = (*common.AccommodationConfigType)(cmd.Type)
	}

	if err := svc.decorateWithTemplateConfig(t, acc); err != nil {
		if _, ok := err.(ErrCannotUseOverridesForConfig); ok {
			return nil, validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key:    "VenueTemplateID",
						Errors: []string{"must be set when any other config field is not"},
					},
				},
			}
		}

		return nil, err
	}

	err = svc.validator.Validate(acc.Config)

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			err = err.PrefixAll("incompatible with template: ")
			return nil, svc.validationMapper.Map(err, UpdateCommand{})
		}
		return nil, err
	}
	

	acc, err = svc.repo.Save(acc)

	if err != nil {
		return nil, err
	}

	return acc, nil
}
