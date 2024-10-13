package accommodations

import (
	"log/slog"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
)

type CreateVenueTemplateHandlerRepo interface {
	Save(org *VenueTemplate) (*VenueTemplate, error)
	ByNameAndVenue(name string, venueId string) (*VenueTemplate, error)
}

type CreateVenueTemplateHandlerVenuesRepo interface {
	Get(id string, orgId string) (*venues.Venue, error)
}

type CreateVenueTemplateCommand struct {
	OrganisationID *string `validate:"required"`
	VenueID        *string `validate:"required"`
	Name           *string `validate:"required,min=3"`
	Type           *string `validate:"required,accommodationtype"`
	MinOccupancy   *int    `validate:"required,min=1"`
	// 100 seems an appropriate max
	MaxOccupancy *int    `validate:"omitnil,gtefield=MinOccupancy"`
	Description  *string `validate:"required,min=10"`
}

type CreateVenueTemplateHandler struct {
	validator  Validator
	repo       CreateVenueTemplateHandlerRepo
	venuesRepo CreateVenueTemplateHandlerVenuesRepo
}

func (h *CreateVenueTemplateHandler) Handle(cmd CreateVenueTemplateCommand) (*VenueTemplate, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	_, err = h.venuesRepo.Get(*cmd.VenueID, *cmd.OrganisationID)

	if err != nil {
		if _, ok := err.(common.ErrNotFound); ok {
			return nil, err
		}

		return nil, err
	}

	blah, err := h.repo.ByNameAndVenue(*cmd.Name, *cmd.VenueID)

	slog.Debug("found", "vt", blah, "err", err)

	if err == nil {
		return nil, validation.ValidationError{
			Errs: []validation.FieldError{
				{
					Key:    "Name",
					Errors: []string{"must be unique"},
				},
			},
		}
	} else {
		_, ok := err.(common.ErrNotFound)

		if !ok {
			return nil, err
		}
	}

	vt := &VenueTemplate{
		VenueID: *cmd.VenueID,
		Template: Template{
			Name:         *cmd.Name,
			MinOccupancy: *cmd.MinOccupancy,
			MaxOccupancy: cmd.MaxOccupancy,
			Description:  *cmd.Description,
			Type:         Type(*cmd.Type),
		},
	}

	return h.repo.Save(vt)
}

func NewCreateVenueTemplateHandler(
	validator Validator,
	repo CreateVenueTemplateHandlerRepo,
	venuesRepo CreateVenueTemplateHandlerVenuesRepo,
) *CreateVenueTemplateHandler {
	return &CreateVenueTemplateHandler{
		validator:  validator,
		repo:       repo,
		venuesRepo: venuesRepo,
	}
}
