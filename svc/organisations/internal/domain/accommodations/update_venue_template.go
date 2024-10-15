package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
)

type UpdateVenueTemplateCommand struct {
	OrganisationID      string  `validate:"required"`
	VenueID             string  `validate:"required"`
	ID                  string  `validate:"required"`
	Name                *string `validate:"omitnil,min=3"`
	Type                *string `validate:"omitnil,accommodationtype"`
	MinOccupancy        *int    `validate:"omitnil,min=1"`
	MaxOccupancy        *int    `validate:"omitnil"`
	NullifyMaxOccupancy bool
	Description         *string `validate:"omitnil,min=10"`
}

func (svc *VenueTemplatesService) Update(cmd UpdateVenueTemplateCommand) (*VenueTemplate, error) {
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
		vt.Type = Type(*cmd.Type)
	}

	if cmd.NullifyMaxOccupancy {
		vt.MaxOccupancy = nil
	} else if cmd.MaxOccupancy != nil && cmd.MinOccupancy == nil && *cmd.MaxOccupancy < vt.MinOccupancy {
		return nil, validation.ValidationError{
			Errs: []validation.FieldError{
				{
					Key:    "MaxOccupancy",
					Errors: []string{"must be greater than min occupancy"},
				},
			},
		}
	} else if cmd.MaxOccupancy != nil {
		// Covers two scenarios:
		// 1. The max and min occupancy are updated, so the validation will have
		// ensured that the max is greater than the min
		// 2. max is updated but min isn't, and max is greater than min
		vt.MaxOccupancy = cmd.MaxOccupancy
	}

	// Note this happens after max occupancy is updated, so we don't need to check
	// anything other than this being nil. If it's not nil we set it...I think
	// Need to write some tests, cause i'm sure i've missed an edge case...
	if cmd.MinOccupancy != nil && vt.MaxOccupancy != nil && *cmd.MinOccupancy > *vt.MaxOccupancy {
		return nil, validation.ValidationError{
			Errs: []validation.FieldError{
				{
					Key:    "MinOccupancy",
					Errors: []string{"must be less than max occupancy"},
				},
			},
		}
	} else if cmd.MinOccupancy != nil {
		vt.MinOccupancy = *cmd.MinOccupancy
	}

	if cmd.Description != nil {
		vt.Description = *cmd.Description
	}

	return svc.repo.Save(vt)
}
