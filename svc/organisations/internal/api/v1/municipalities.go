package v1

import (
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/requests"
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/responses"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/municipalities"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
	"github.com/labstack/echo/v4"
)

type MunicipalitiesListHandler interface {
	Handle(cmd municipalities.ListCommand) (municipalities.Municipalities, common.PaginationResult, error)
}

type MunicipalitiesControllerConfig interface{}

type MunicipalitiesController struct {
	cfg              MunicipalitiesControllerConfig
	list             MunicipalitiesListHandler
	validationMapper *validation.ValidationMapper
}

func (c *MunicipalitiesController) RegisterRoutes(api *echo.Group) {
	g := api.Group("/municipalities")
	g.GET("", c.List).Name = "v1.municipalities.list"
}

func NewMunicipalitiesController(
	cfg MunicipalitiesControllerConfig,
	validationMapper *validation.ValidationMapper,
	list MunicipalitiesListHandler,
) *MunicipalitiesController {
	return &MunicipalitiesController{
		cfg:              cfg,
		list:             list,
		validationMapper: validationMapper,
	}
}

// @Summary		List all municipalities that can be used
// @Tags			Municipalities
// @Accept			json
// @Produce		json
// @Success		200	{object}	responses.ListMunicipalitiesResponse
// @Failure		422	{object}	responses.ValidationErrorResponse
// @Failure		404	{object}	responses.GenericErrorResponse
// @Failure		400	{object}	responses.GenericErrorResponse
// @Failure		500	{object}	responses.GenericErrorResponse
// @Router			/v1/municipalities [get]
// @Param			request	query requests.ListMunicipalitiesRequest	true "Query params"
func (c *MunicipalitiesController) List(ctx echo.Context) error {
	req := requests.ListMunicipalitiesRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	cmd := req.ToCommand()

	results, pagination, err := c.list.Handle(cmd)

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}
		return err
	}

	resp := responses.ListMunicipalitiesResponse{
		Meta: responses.ListResponseMeta{
			SortOptionsResponseMeta: responses.SortOptionsResponseMeta{
				OrderDirection: string(cmd.OrderDirection),
				OrderBy:        string(cmd.OrderBy),
			},
			PaginationResponseMeta: responses.PaginationResponseMeta{
				Page:         pagination.Page,
				PerPage:      pagination.PerPage,
				TotalPages:   pagination.TotalPages,
				TotalResults: pagination.Total,
			},
		},
		Data: responses.MunicipalitiesFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}
