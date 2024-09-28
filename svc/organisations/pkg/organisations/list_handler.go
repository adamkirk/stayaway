package organisations

import "github.com/adamkirk-stayaway/organisations/pkg/model"


type ListHandlerRepo interface {
	Paginate(orderBy model.OrganisationSortBy, orderDir model.SortDirection, page int, perPage int) (model.Organisations, model.PaginationResult, error)
}

type ListCommand struct {
	OrderDirection model.SortDirection
	OrderBy model.OrganisationSortBy
	Page int
	PerPage int
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: model.SortAsc,
		OrderBy: model.OrganisationSortByName,
		Page: 1,
		PerPage: 50,
	}
}

type ListHandler struct {
	repo ListHandlerRepo
}

func (h *ListHandler) Handle(cmd ListCommand) (model.Organisations, model.PaginationResult, error) {
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