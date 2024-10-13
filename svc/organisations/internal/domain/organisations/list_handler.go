package organisations

import "github.com/adamkirk-stayaway/organisations/internal/domain/common"

type ListHandlerRepo interface {
	Paginate(orderBy SortBy, orderDir common.SortDirection, page int, perPage int) (Organisations, common.PaginationResult, error)
}

type ListCommand struct {
	OrderDirection common.SortDirection
	OrderBy        SortBy
	Page           int
	PerPage        int
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
	repo ListHandlerRepo
}

func (h *ListHandler) Handle(cmd ListCommand) (Organisations, common.PaginationResult, error) {
	return h.repo.Paginate(
		cmd.OrderBy,
		cmd.OrderDirection,
		cmd.Page,
		cmd.PerPage,
	)
}

func NewListHandler(repo ListHandlerRepo) *ListHandler {
	return &ListHandler{
		repo: repo,
	}
}
