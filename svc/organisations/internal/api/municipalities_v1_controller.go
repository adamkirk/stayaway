package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/municipalities"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/labstack/echo/v4"
)


type MunicipalitiesListHandler interface {
	Handle(cmd municipalities.ListCommand) (model.Municipalities, model.PaginationResult, error)
}

type MunicipalitiesV1ControllerConfig interface {}

type MunicipalitiesV1Controller struct {
	cfg MunicipalitiesV1ControllerConfig
	list MunicipalitiesListHandler
	validationMapper *ValidationMapper
}

func (c *MunicipalitiesV1Controller) RegisterRoutes(api *echo.Group) {
	g := api.Group("/v1/municipalities")
	g.GET("", c.List).Name = "v1.municipalities.list"
}

func NewMunicipalitiesV1Controller(
	cfg MunicipalitiesV1ControllerConfig,
	validationMapper *ValidationMapper,
	list MunicipalitiesListHandler,
) *MunicipalitiesV1Controller {
	return &MunicipalitiesV1Controller{
		cfg: cfg,
		list: list,
		validationMapper: validationMapper,
	}
}

//	@Summary		List all venues for an organisation
//	@Tags			Municipalities
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	V1ListMunicipalitiesResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/municipalities [get]
//	@Param			request	query V1ListMunicipalitiesRequest	true "Query params"
func (c *MunicipalitiesV1Controller) List(ctx echo.Context) error {
	req := V1ListMunicipalitiesRequest{}

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

	resp := V1ListMunicipalitiesResponse{
		Meta: V1ListResponseMeta{
			V1SortOptionsResponseMeta: V1SortOptionsResponseMeta{
				OrderDirection: string(cmd.OrderDirection),
				OrderBy: string(cmd.OrderBy),
			},
			V1PaginationResponseMeta: V1PaginationResponseMeta{
				Page: pagination.Page,
				PerPage: pagination.PerPage,
				TotalPages: pagination.TotalPages,
				TotalResults: pagination.Total,
			},
		},
		Data: V1MunicipalitiesFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}
