package venues

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type UpdateCommand struct {
	ID                  *string `validate:"required"`
	OrganisationID      *string `validate:"required"`
	Name                *string `validate:"omitnil,min=3"`
	Slug                *string `validate:"omitnil,min=3,slug"`
	Type                *string `validate:"omitnil,venuetype"`
	AddressLine1        *string `validate:"omitnil,min=1"`
	AddressLine2        *string `validate:"omitnil,min=1"`
	NullifyAddressLine2 bool
	Municipality        *string  `validate:"omitnil,min=1"`
	PostCode            *string  `validate:"omitnil,postcode"`
	Lat                 *float64 `validate:"omitnil,min=0"`
	Long                *float64 `validate:"omitnil,min=0"`
}

func (svc *Service) Update(cmd UpdateCommand) (*Venue, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	v, err := svc.repo.Get(*cmd.ID, *cmd.OrganisationID)

	if err != nil {
		return nil, err
	}

	if cmd.Slug != nil {
		venueBySlug, err := svc.repo.BySlugAndOrganisation(*cmd.Slug, *cmd.OrganisationID)

		if venueBySlug != nil && venueBySlug.ID != v.ID {
			return nil, validation.ValidationError{
				Errs: []validation.FieldError{
					{
						Key:    "Slug",
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

		v.Slug = *cmd.Slug
	}

	if cmd.Name != nil {
		v.Name = *cmd.Name
	}

	if cmd.Type != nil {
		v.Type = Type(*cmd.Slug)
	}

	if cmd.AddressLine1 != nil {
		v.Address.Line1 = *cmd.AddressLine1
	}

	if cmd.NullifyAddressLine2 {
		v.Address.Line2 = nil
	} else if cmd.AddressLine2 != nil {
		v.Address.Line2 = cmd.AddressLine2
	}

	if cmd.Municipality != nil {
		v.Address.Municipality = *cmd.Municipality
	}

	if cmd.PostCode != nil {
		v.Address.PostCode = *cmd.PostCode
	}

	if cmd.Lat != nil {
		v.Address.Coordinates.Lat = *cmd.Lat
	}

	if cmd.Long != nil {
		v.Address.Coordinates.Long = *cmd.Long
	}

	return svc.repo.Save(v)
}
