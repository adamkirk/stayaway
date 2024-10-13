package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
)

type ListVenueTemplatesHandlerRepo interface {
	Paginate(p PaginationFilter, search SearchFilter) (VenueTemplates, common.PaginationResult, error)
}

type ListVenueTemplatesHandlerVenuesRepo interface {
	Get(id string, orgId string) (*venues.Venue, error)
}

type ListVenueTemplatesCommand struct {
	OrganisationID string `validate:"required"`
	VenueID string `validate:"required"`
	NamePrefix *string `validate:"omitnil,min=3"`
	OrderDirection common.SortDirection `validate:"required"`
	OrderBy SortBy `validate:"required"`
	Page int `validate:"required,min=1"`
	PerPage int `validate:"required"`
}

func NewListVenueTemplatesCommand() ListVenueTemplatesCommand {
	return ListVenueTemplatesCommand{
		OrderDirection: common.SortAsc,
		OrderBy: SortByName,
		Page: 1,
		PerPage: 50,
	}
}


type ListVenueTemplatesHandler struct {
	validator Validator
	repo ListVenueTemplatesHandlerRepo
	venuesRepo ListVenueTemplatesHandlerVenuesRepo
}

func (h *ListVenueTemplatesHandler) Handle(cmd ListVenueTemplatesCommand) (VenueTemplates, common.PaginationResult, error) {
	err := h.validator.Validate(cmd)

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
	_, err = h.venuesRepo.Get(cmd.VenueID, cmd.OrganisationID)

	if err != nil {
		if _, ok := err.(common.ErrNotFound); ok {
			return nil, common.PaginationResult{}, err
		}

		return nil, common.PaginationResult{}, err
	}

	p := PaginationFilter{
		OrderBy: cmd.OrderBy,
		OrderDir: cmd.OrderDirection,
		Page: cmd.Page,
		PerPage:cmd.PerPage,
	}

	s := SearchFilter{
		NamePrefix: cmd.NamePrefix,
		VenueID: []string{
			cmd.VenueID,
		},
	}

	return h.repo.Paginate(p, s)
}

func NewListVenueTemplatesHandler(
	validator Validator, 
	repo ListVenueTemplatesHandlerRepo,
	venuesRepo ListVenueTemplatesHandlerVenuesRepo,
) *ListVenueTemplatesHandler {
	return &ListVenueTemplatesHandler{
		validator: validator,
		repo: repo,
		venuesRepo: venuesRepo,
	}
}