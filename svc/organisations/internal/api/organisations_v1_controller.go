package api

import (
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

type OrganisationsV1ControllerConfig interface {}

type OrganisationsV1Controller struct {
	cfg OrganisationsV1ControllerConfig
	get OrganisationsGetHandler
	list OrganisationsListHandler
	create OrganisationsCreateHandler
	delete OrganisationsDeleteHandler
	update OrganisationsUpdateHandler
	validationMapper *ValidationMapper
}

func (c *OrganisationsV1Controller) RegisterRoutes(api *echo.Group) {
	g := api.Group("/v1/organisations")
	g.POST("", c.Create).Name = "v1.organisations.create"
	g.DELETE("/:id", c.Delete).Name = "v1.organisations.delete"
	g.GET("/:id", c.Get).Name = "v1.organisations.get"
	g.PATCH("/:id", c.Patch).Name = "v1.organisations.patch"
	g.GET("", c.List).Name = "v1.organisations.list"
}

func NewOrganisationsV1Controller(
	cfg OrganisationsV1ControllerConfig,
	get OrganisationsGetHandler,
	list OrganisationsListHandler,
	create OrganisationsCreateHandler,
	delete OrganisationsDeleteHandler,
	update OrganisationsUpdateHandler,
	validationMapper *ValidationMapper,
) *OrganisationsV1Controller {
	return &OrganisationsV1Controller{
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
//	@Success		200	{object}	V1ListOrganisationsResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations [get]
//	@Param			request	query	V1ListOrganisationsRequest	true "Query params"
func (c *OrganisationsV1Controller) List(ctx echo.Context) error {
	req := V1ListOrganisationsRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	cmd := req.ToCommand()

	results, pagination, err := c.list.Handle(cmd)

	if err != nil {
		return err
	}
	
	resp := V1ListOrganisationsResponse{
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
		Data: V1OrganisationsFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Create an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	V1PostOrganisationResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations [post]
//	@Param			Organisation	body		V1PostOrganisationRequest	true	"Organisation definition"
func (c *OrganisationsV1Controller) Create(ctx echo.Context) error {
	req := V1PostOrganisationRequest{}

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

	resp := V1PostOrganisationResponse{
		Data: V1OrganisationFromModel(org),
	}

	ctx.JSON(201, resp)

	return nil
}

//	@Summary		Get an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	V1GetOrganisationResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{id} [get]
//	@Param			id	path	string	true	"The Organisation ID"
func (c *OrganisationsV1Controller) Get(ctx echo.Context) error {
	req := V1GetOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.get.Handle(req.ToCommand())

	if err != nil {
		return err
	}

	resp := V1GetOrganisationResponse{
		Data: V1OrganisationFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Update an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	V1PatchOrganisationResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{id} [patch]
//	@Param			id	path	string	true	"The Organisation ID"
//	@Param			Changes	body		V1PatchOrganisationRequest	true	"Organisation definition"
func (c *OrganisationsV1Controller) Patch(ctx echo.Context) error {
	req := V1PatchOrganisationRequest{}

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

	resp := V1PatchOrganisationResponse{
		Data: V1OrganisationFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Delete an organisation
//	@Tags			Organisations
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{id} [delete]
//	@Param			id	path	string	true	"The Organisation ID"
func (c *OrganisationsV1Controller) Delete(ctx echo.Context) error {
	req := V1DeleteOrganisationRequest{}

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