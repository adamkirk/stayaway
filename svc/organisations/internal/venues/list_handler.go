package venues

import "github.com/adamkirk-stayaway/organisations/internal/model"


type ListHandlerRepo interface {
	Paginate(p model.VenuePaginationFilter, search model.VenueSearchFilter) (model.Venues, model.PaginationResult, error)
}

type ListCommand struct {
	OrganisationID string `validate:"required"`
	OrderDirection model.SortDirection `validate:"required"`
	OrderBy model.VenueSortBy `validate:"required"`
	Page int `validate:"required,min=1"`
	PerPage int `validate:"required"`
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: model.SortAsc,
		OrderBy: model.VenueSortByName,
		Page: 1,
		PerPage: 50,
	}
}

type ListHandler struct {
	repo ListHandlerRepo
	validator Validator
}

func (h *ListHandler) Handle(cmd ListCommand) (model.Venues, model.PaginationResult, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, model.PaginationResult{}, err
	}

	return h.repo.Paginate(
		model.VenuePaginationFilter{
			OrderBy: cmd.OrderBy,
			OrderDir: cmd.OrderDirection,
			Page: cmd.Page,
			PerPage: cmd.PerPage,
		},
		model.VenueSearchFilter{
			OrganisationID: []string{cmd.OrganisationID},
		},
	)
}

func NewListHandler(validator Validator, repo ListHandlerRepo) *ListHandler {
	return &ListHandler{
		repo: repo,
		validator: validator,
	}
}