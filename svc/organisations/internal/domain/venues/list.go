package venues

import "github.com/adamkirk-stayaway/organisations/internal/domain/common"

type ListCommand struct {
	OrganisationID string               `validate:"required"`
	OrderDirection common.SortDirection `validate:"required,orderdir"`
	OrderBy        SortBy               `validate:"required,venues_sortfield"`
	Page           int                  `validate:"required,min=1"`
	PerPage        int                  `validate:"required,min=1"`
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: common.SortAsc,
		OrderBy:        SortByName,
		Page:           1,
		PerPage:        50,
	}
}

func (svc *Service) List(cmd ListCommand) (Venues, common.PaginationResult, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	return svc.repo.Paginate(
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
