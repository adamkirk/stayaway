package venues

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
)

type CreateHandlerRepo interface {
	Save(org *Venue) (*Venue, error)
	BySlugAndOrganisation(slug string, orgId string) (*Venue, error)
}

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

type CreateHandler struct {
	validator Validator
	repo      CreateHandlerRepo
}

func (h *CreateHandler) Handle(cmd CreateCommand) (*Venue, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	venueBySlug, err := h.repo.BySlugAndOrganisation(*cmd.Slug, *cmd.OrganisationID)

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

	return h.repo.Save(v)
}

func NewCreateHandler(validator Validator, repo CreateHandlerRepo) *CreateHandler {
	return &CreateHandler{
		validator: validator,
		repo:      repo,
	}
}
