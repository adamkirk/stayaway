package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
)

type ListVenueTemplatesCommand struct {
	OrganisationID string               `validate:"required"`
	VenueID        string               `validate:"required"`
	NamePrefix     *string              `validate:"omitnil,min=3"`
	OrderDirection common.SortDirection `validate:"required"`
	OrderBy        SortBy               `validate:"required"`
	Page           int                  `validate:"required,min=1"`
	PerPage        int                  `validate:"required"`
}

func NewListVenueTemplatesCommand() ListVenueTemplatesCommand {
	return ListVenueTemplatesCommand{
		OrderDirection: common.SortAsc,
		OrderBy:        SortByName,
		Page:           1,
		PerPage:        50,
	}
}

func (svc *VenueTemplatesService) List(cmd ListVenueTemplatesCommand) (VenueTemplates, common.PaginationResult, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, common.PaginationResult{}, err
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
			return nil, common.PaginationResult{}, err
		}

		return nil, common.PaginationResult{}, err
	}

	p := PaginationFilter{
		OrderBy:  cmd.OrderBy,
		OrderDir: cmd.OrderDirection,
		Page:     cmd.Page,
		PerPage:  cmd.PerPage,
	}

	s := SearchFilter{
		NamePrefix: cmd.NamePrefix,
		VenueID: []string{
			cmd.VenueID,
		},
	}

	return svc.repo.Paginate(p, s)
}
