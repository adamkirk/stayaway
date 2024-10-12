package v1

import (
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/requests"
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/responses"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/organisations"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
	"github.com/labstack/echo/v4"
)

type OrganisationsGetHandler interface {
	Handle(organisations.GetCommand) (*organisations.Organisation, error)
}

type OrganisationsListHandler interface {
	Handle(cmd organisations.ListCommand) (organisations.Organisations, common.PaginationResult, error)
}

type OrganisationsCreateHandler interface {
	Handle(cmd organisations.CreateCommand) (*organisations.Organisation, error)
}

type OrganisationsDeleteHandler interface {
	Handle(cmd organisations.DeleteCommand) error
}


type OrganisationsUpdateHandler interface {
	Handle(cmd organisations.UpdateCommand) (*organisations.Organisation, error)
}

type OrganisationsControllerConfig interface {}

type OrganisationsController struct {
	cfg OrganisationsControllerConfig
	get OrganisationsGetHandler
	list OrganisationsListHandler
	create OrganisationsCreateHandler
	delete OrganisationsDeleteHandler
	update OrganisationsUpdateHandler
	validationMapper *validation.ValidationMapper
}

func (c *OrganisationsController) RegisterRoutes(api *echo.Group) {
	g := api.Group("/organisations")
	g.POST("", c.Create).Name = "v1.organisations.create"
	g.DELETE("/:id", c.Delete).Name = "v1.organisations.delete"
	g.GET("/:id", c.Get).Name = "v1.organisations.get"
	g.PATCH("/:id", c.Patch).Name = "v1.organisations.patch"
	g.GET("", c.List).Name = "v1.organisations.list"
}

func NewOrganisationsController(
	cfg OrganisationsControllerConfig,
	get OrganisationsGetHandler,
	list OrganisationsListHandler,
	create OrganisationsCreateHandler,
	delete OrganisationsDeleteHandler,
	update OrganisationsUpdateHandler,
	validationMapper *validation.ValidationMapper,
) *OrganisationsController {
	return &OrganisationsController{
		cfg: cfg,
		get: get,
		list: list,
		create: create,
		delete: delete,
		update: update,
		validationMapper: validationMapper,
	}
}

//	@Summary		List all organisations
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ListOrganisationsResponse
//	@Failure		422	{object}	ValidationErrorResponse
//	@Failure		404	{object}	GenericErrorResponse
//	@Failure		400	{object}	GenericErrorResponse
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/v1/organisations [get]
//	@Param			request	query	requests.ListOrganisationsRequest	true "Query params"
func (c *OrganisationsController) List(ctx echo.Context) error {
	req := requests.ListOrganisationsRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	cmd := req.ToCommand()

	results, pagination, err := c.list.Handle(cmd)

	if err != nil {
		return err
	}
	
	resp := responses.ListOrganisationsResponse{
		Meta: responses.ListResponseMeta{
			SortOptionsResponseMeta: responses.SortOptionsResponseMeta{
				OrderDirection: string(cmd.OrderDirection),
				OrderBy: string(cmd.OrderBy),
			},
			PaginationResponseMeta: responses.PaginationResponseMeta{
				Page: pagination.Page,
				PerPage: pagination.PerPage,
				TotalPages: pagination.TotalPages,
				TotalResults: pagination.Total,
			},
		},
		Data: responses.OrganisationsFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Create an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	responses.PostOrganisationResponse
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations [post]
//	@Param			Organisation	body		requests.PostOrganisationRequest	true	"Organisation definition"
func (c *OrganisationsController) Create(ctx echo.Context) error {
	req := requests.PostOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.create.Handle(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}
		return err
	}

	resp := responses.PostOrganisationResponse{
		Data: responses.OrganisationFromModel(org),
	}

	ctx.JSON(201, resp)

	return nil
}

//	@Summary		Get an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.GetOrganisationResponse
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{id} [get]
//	@Param			id	path	string	true	"The Organisation ID"
func (c *OrganisationsController) Get(ctx echo.Context) error {
	req := requests.GetOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.get.Handle(req.ToCommand())

	if err != nil {
		return err
	}

	resp := responses.GetOrganisationResponse{
		Data: responses.OrganisationFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Update an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.PatchOrganisationResponse
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{id} [patch]
//	@Param			id	path	string	true	"The Organisation ID"
//	@Param			Changes	body		requests.PatchOrganisationRequest	true	"Organisation definition"
func (c *OrganisationsController) Patch(ctx echo.Context) error {
	req := requests.PatchOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.update.Handle(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}

		return err
	}

	resp := responses.PatchOrganisationResponse{
		Data: responses.OrganisationFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Delete an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{id} [delete]
//	@Param			id	path	string	true	"The Organisation ID"
func (c *OrganisationsController) Delete(ctx echo.Context) error {
	req := requests.DeleteOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	err := c.delete.Handle(req.ToCommand())

	if err != nil {
		return err
	}
	
	ctx.NoContent(204)

	return nil
}