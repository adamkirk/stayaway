package organisations

import "github.com/adamkirk-stayaway/organisations/internal/domain/common"

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

func (svc *Service) List(cmd ListCommand) (Organisations, common.PaginationResult, error) {
	return svc.repo.Paginate(
		cmd.OrderBy,
		cmd.OrderDirection,
		cmd.Page,
		cmd.PerPage,
	)
}
