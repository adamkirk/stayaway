package templates

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type UpdateVenueTemplateCommand struct {
	OrganisationID      string  `validate:"required"`
	VenueID             string  `validate:"required"`
	ID                  string  `validate:"required"`
	Name                *string `validate:"omitnil,min=3"`
	Type                *string `validate:"omitnil,accommodationtype"`
	MinOccupancy        *int    `validate:"omitnil,min=1"`
	MaxOccupancy        *int    `validate:"omitnil,min=1"`
	Description         *string `validate:"omitnil,min=10"`
}

func (svc *Service) Update(cmd UpdateVenueTemplateCommand) (*VenueTemplate, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	// This is largely to make sure the venue and organisation exist
	// Then we include the venue id in the get query to ensure the template
	// belongs to the given venue.
	// Feel like generally there is a better pattern for this rather than
	// keeping the full hierarchy of ids around, but this is simple enough
	// for now.
	// Applies to other areas...
	_, err = svc.venuesRepo.Get(cmd.VenueID, cmd.OrganisationID)

	if err != nil {
		if _, ok := err.(common.ErrNotFound); ok {
			return nil, err
		}

		return nil, err
	}

	vt, err := svc.repo.Get(cmd.ID, cmd.VenueID)

	if err != nil {
		return nil, err
	}

	if cmd.Name != nil {
		byName, err := svc.repo.ByNameAndVenue(*cmd.Name, cmd.VenueID)

		if byName != nil && byName.ID != vt.ID {
			return nil, validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key:    "Name",
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

		vt.Name = *cmd.Name
	}

	if cmd.Type != nil {
		vt.Type = common.AccommodationConfigType(*cmd.Type)
	}

	if cmd.MaxOccupancy != nil {
		vt.MaxOccupancy = *cmd.MaxOccupancy
	}

	if cmd.MinOccupancy != nil {
		vt.MinOccupancy = *cmd.MinOccupancy
	}

	if cmd.Description != nil {
		vt.Description = *cmd.Description
	}

	return svc.repo.Save(vt)
}
