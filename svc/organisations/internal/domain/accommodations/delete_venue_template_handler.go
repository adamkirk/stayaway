package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
)

type DeleteVenueTemplateHandlerRepo interface {
	Get(id string, venueId string) (*VenueTemplate, error)
	Delete(*VenueTemplate) (error)
}

type DeleteVenueTemplateHandlerVenuesRepo interface {
	Get(id string, orgId string) (*venues.Venue, error)
}

type DeleteVenueTemplateCommand struct {
	OrganisationID string `validate:"required"`
	VenueID string `validate:"required"`
	ID string `validate:"required"`
}

type DeleteVenueTemplateHandler struct {
	validator Validator
	repo DeleteVenueTemplateHandlerRepo
	venuesRepo DeleteVenueTemplateHandlerVenuesRepo
}

func (h *DeleteVenueTemplateHandler) Handle(cmd DeleteVenueTemplateCommand) (error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return err
	}

	// This is largely to make sure the venue and organisation exist
	// Then we include the venue id in the get query to ensure the template
	// belongs to the given venue.
	// Feel like generally there is a better pattern for this rather than
	// keeping the full hierarchy of ids around, but this is simple enough 
	// for now.
	// Applies to other areas...
	_, err = h.venuesRepo.Get(cmd.VenueID, cmd.OrganisationID)

	if err != nil {
		if _, ok := err.(common.ErrNotFound); ok {
			return err
		}

		return err
	}

	vt, err := h.repo.Get(cmd.ID, cmd.VenueID)

	if err != nil {
		return err
	}

	return h.repo.Delete(vt)
}

func NewDeleteVenueTemplateHandler(
	validator Validator, 
	repo DeleteVenueTemplateHandlerRepo,
	venuesRepo DeleteVenueTemplateHandlerVenuesRepo,
) *DeleteVenueTemplateHandler {
	return &DeleteVenueTemplateHandler{
		validator: validator,
		repo: repo,
		venuesRepo: venuesRepo,
	}
}