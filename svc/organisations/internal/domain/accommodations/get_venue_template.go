package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
)

type GetVenueTemplateCommand struct {
	OrganisationID string `validate:"required"`
	VenueID        string `validate:"required"`
	ID             string `validate:"required"`
}

func (svc *VenueTemplatesService) Get(cmd GetVenueTemplateCommand) (*VenueTemplate, error) {
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

	return svc.repo.Get(cmd.ID, cmd.VenueID)
}
