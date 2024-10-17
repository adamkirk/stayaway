package venues

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type CreateCommand struct {
	OrganisationID *string  `validate:"required"`
	Name           *string  `validate:"required,min=3"`
	Slug           *string  `validate:"required,min=3,slug"`
	Type           *string  `validate:"required,venuetype"`
	AddressLine1   *string  `validate:"required,min=1"`
	AddressLine2   *string  `validate:"omitnil,min=1"`
	Municipality   *string  `validate:"required,min=1"`
	PostCode       *string  `validate:"required,postcode"`
	Lat            *float64 `validate:"required,min=0"`
	Long           *float64 `validate:"required,min=0"`
}

func (svc *Service) Create(cmd CreateCommand) (*Venue, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	venueBySlug, err := svc.repo.BySlugAndOrganisation(*cmd.Slug, *cmd.OrganisationID)

	if venueBySlug != nil {
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

	v := &Venue{
		ID:             svc.idGen.Generate(),
		OrganisationID: *cmd.OrganisationID,
		Name:           *cmd.Name,
		Slug:           *cmd.Slug,
		Type:           Type(*cmd.Type),
		Address: &Address{
			Line1:        *cmd.AddressLine1,
			Line2:        cmd.AddressLine2,
			Municipality: *cmd.Municipality,
			PostCode:     *cmd.PostCode,
			Coordinates: &Coordinates{
				Lat:  *cmd.Lat,
				Long: *cmd.Long,
			},
		},
	}

	return svc.repo.Save(v)
}
