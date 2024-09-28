package municipalities

import "github.com/adamkirk-stayaway/organisations/pkg/model"


type ListHandlerRepo interface {
	Paginate(p model.MunicipalityPaginationFilter, search model.MunicipalitySearchFilter) (model.Municipalities, model.PaginationResult, error)
}

type ListCommand struct {
	OrderDirection model.SortDirection `validate:"required"`
	OrderBy model.MunicipalitySortBy `validate:"required"`
	Page int `validate:"required,min=1"`
	PerPage int `validate:"required,min=1,max=100"`
	Country []string
	NamePrefix *string `validate:"omitnil,min=3"`
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: model.SortAsc,
		OrderBy: model.MunicipalitySortByName,
		Page: 1,
		PerPage: 50,
		Country: []string{},
	}
}

type ListHandler struct {
	repo ListHandlerRepo
	validator Validator
}

func (h *ListHandler) Handle(cmd ListCommand) (model.Municipalities, model.PaginationResult, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, model.PaginationResult{}, err
	}

	return h.repo.Paginate(
		model.MunicipalityPaginationFilter{
			Page: cmd.Page,
			PerPage: cmd.PerPage,
			OrderBy: cmd.OrderBy,
			OrderDir: cmd.OrderDirection,
		},
		model.MunicipalitySearchFilter{
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