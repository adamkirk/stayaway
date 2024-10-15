package municipalities

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
)

type ListCommand struct {
	OrderDirection common.SortDirection `validate:"required"`
	OrderBy        SortBy               `validate:"required"`
	Page           int                  `validate:"required,min=1"`
	PerPage        int                  `validate:"required,min=1,max=100"`
	Country        []string
	NamePrefix     *string `validate:"omitnil,min=3"`
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: common.SortAsc,
		OrderBy:        SortByName,
		Page:           1,
		PerPage:        50,
		Country:        []string{},
	}
}


func (svc *Service) List(cmd ListCommand) (Municipalities, common.PaginationResult, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	return svc.repo.Paginate(
		PaginationFilter{
			Page:     cmd.Page,
			PerPage:  cmd.PerPage,
			OrderBy:  cmd.OrderBy,
			OrderDir: cmd.OrderDirection,
		},
		SearchFilter{
			Country:    cmd.Country,
			NamePrefix: cmd.NamePrefix,
		},
	)
}
