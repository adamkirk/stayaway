package venues

import "github.com/adamkirk-stayaway/organisations/internal/domain/common"

type ListHandlerRepo interface {
	Paginate(p PaginationFilter, search SearchFilter) (Venues, common.PaginationResult, error)
}

type ListCommand struct {
	OrganisationID string               `validate:"required"`
	OrderDirection common.SortDirection `validate:"required"`
	OrderBy        SortBy               `validate:"required"`
	Page           int                  `validate:"required,min=1"`
	PerPage        int                  `validate:"required"`
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: common.SortAsc,
		OrderBy:        SortByName,
		Page:           1,
		PerPage:        50,
	}
}

type ListHandler struct {
	repo      ListHandlerRepo
	validator Validator
}

func (h *ListHandler) Handle(cmd ListCommand) (Venues, common.PaginationResult, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	return h.repo.Paginate(
		PaginationFilter{
			OrderBy:  cmd.OrderBy,
			OrderDir: cmd.OrderDirection,
			Page:     cmd.Page,
			PerPage:  cmd.PerPage,
		},
		SearchFilter{
			OrganisationID: []string{cmd.OrganisationID},
		},
	)
}

func NewListHandler(validator Validator, repo ListHandlerRepo) *ListHandler {
	return &ListHandler{
		repo:      repo,
		validator: validator,
	}
}
