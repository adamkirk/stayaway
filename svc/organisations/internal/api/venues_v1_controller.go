package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/adamkirk-stayaway/organisations/pkg/venues"
	"github.com/labstack/echo/v4"
)

type VenuesGetHandler interface {
	Handle(venues.GetCommand) (*model.Venue, error)
}

type VenuesListHandler interface {
	Handle(cmd venues.ListCommand) (model.Venues, model.PaginationResult, error)
}

type VenuesCreateHandler interface {
	Handle(cmd venues.CreateCommand) (*model.Venue, error)
}

type VenuesDeleteHandler interface {
	Handle(cmd venues.DeleteCommand) error
}


type VenuesUpdateHandler interface {
	Handle(cmd venues.UpdateCommand) (*model.Venue, error)
}

type VenuesV1ControllerConfig interface {}

type VenuesV1Controller struct {
	cfg VenuesV1ControllerConfig
	// repo VenuesRepo
	get VenuesGetHandler
	list VenuesListHandler
	create VenuesCreateHandler
	delete VenuesDeleteHandler
	update VenuesUpdateHandler
	validationMapper *ValidationMapper
}

func (c *VenuesV1Controller) RegisterRoutes(api *echo.Group) {
	g := api.Group("/v1/organisations/:organisationId/venues")
	g.POST("", c.Create).Name = "v1.organisations.venues.create"
	g.DELETE("/:id", c.Delete).Name = "v1.organisations.venues.delete"
	g.GET("/:id", c.Get).Name = "v1.organisations.venues.get"
	g.PATCH("/:id", c.Patch).Name = "v1.organisations.venues.patch"
	g.GET("", c.List).Name = "v1.organisations.venues.list"
}

func NewVenuesV1Controller(
	cfg VenuesV1ControllerConfig,
	create VenuesCreateHandler,
	validationMapper *ValidationMapper,
	get VenuesGetHandler,
	list VenuesListHandler,
	delete VenuesDeleteHandler,
	update VenuesUpdateHandler,
) *VenuesV1Controller {
	return &VenuesV1Controller{
		cfg: cfg,
		get: get,
		list: list,
		create: create,
		delete: delete,
		update: update,
		validationMapper: validationMapper,
	}
}

//	@Summary		List all venues for an organisation
//	@Tags			Venues
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	V1ListVenuesResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues [get]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			request	query V1ListVenuesRequest	true "Query params"
func (c *VenuesV1Controller) List(ctx echo.Context) error {
	req := V1ListVenuesRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	cmd := req.ToCommand()

	results, pagination, err := c.list.Handle(cmd)

	if err != nil {
		return err
	}
	
	resp := V1ListVenuesResponse{
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
		Data: V1VenuesFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Create a venue for an organisation
//	@Tags			Venues
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	V1PostVenueResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues [post]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			Venue	body		V1PostVenueRequest	true	"Venue definition"
func (c *VenuesV1Controller) Create(ctx echo.Context) error {
	req := V1PostVenueRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	v, err := c.create.Handle(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}
		return err
	}

	resp := V1PostVenueResponse{
		Data: V1VenueFromModel(v),
	}

	ctx.JSON(201, resp)

	return nil
}

//	@Summary		Get a venue for an organisation
//	@Tags			Venues
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	V1PostVenueResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{id} [get]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			id	path	string	true	"The Venues ID"
func (c *VenuesV1Controller) Get(ctx echo.Context) error {
	req := V1GetVenueRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.get.Handle(req.ToCommand())

	if err != nil {
		return err
	}

	resp := V1GetVenueResponse{
		Data: V1VenueFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Update a venue for an organisation
//	@Tags			Venues
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	V1PatchVenueResponse
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{id} [patch]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			id	path	string	true	"The Venues ID"
//	@Param			Changes	body		V1PatchVenueRequest	true	"Venue changes"
func (c *VenuesV1Controller) Patch(ctx echo.Context) error {
	req := V1PatchVenueRequest{}
	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	venue, err := c.update.Handle(req.ToCommand())

	if err != nil {
		return err
	}

	resp := V1PatchVenueResponse{
		Data: V1VenueFromModel(venue),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Delete a venue from an organisation
//	@Tags			Venues
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		422	{object}	V1ValidationErrorResponse
//	@Failure		404	{object}	V1GenericErrorResponse
//	@Failure		400	{object}	V1GenericErrorResponse
//	@Failure		500	{object}	V1GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{id} [delete]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			id	path	string	true	"The Venues ID"
func (c *VenuesV1Controller) Delete(ctx echo.Context) error {
	req := V1DeleteVenueRequest{}

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