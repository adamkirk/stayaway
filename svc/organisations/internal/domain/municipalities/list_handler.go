package municipalities

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
)


type ListHandlerRepo interface {
	Paginate(p PaginationFilter, search SearchFilter) (Municipalities, common.PaginationResult, error)
}

type ListCommand struct {
	OrderDirection common.SortDirection `validate:"required"`
	OrderBy SortBy `validate:"required"`
	Page int `validate:"required,min=1"`
	PerPage int `validate:"required,min=1,max=100"`
	Country []string
	NamePrefix *string `validate:"omitnil,min=3"`
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: common.SortAsc,
		OrderBy: SortByName,
		Page: 1,
		PerPage: 50,
		Country: []string{},
	}
}

type ListHandler struct {
	repo ListHandlerRepo
	validator Validator
}

func (h *ListHandler) Handle(cmd ListCommand) (Municipalities, common.PaginationResult, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	return h.repo.Paginate(
		PaginationFilter{
			Page: cmd.Page,
			PerPage: cmd.PerPage,
			OrderBy: cmd.OrderBy,
			OrderDir: cmd.OrderDirection,
		},
		SearchFilter{
			Country: cmd.Country,
			NamePrefix: cmd.NamePrefix,
		},
	)
}

func NewListHandler(validator Validator, repo ListHandlerRepo) *ListHandler {
	return &ListHandler{
		repo: repo,
		validator: validator,
	}
}