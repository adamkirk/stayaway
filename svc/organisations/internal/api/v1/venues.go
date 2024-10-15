package v1

import (
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/requests"
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/responses"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
	"github.com/labstack/echo/v4"
)

type VenuesService interface {
	Get(venues.GetCommand) (*venues.Venue, error)
	Create(cmd venues.CreateCommand) (*venues.Venue, error)
	Delete(cmd venues.DeleteCommand) error
	List(cmd venues.ListCommand) (venues.Venues, common.PaginationResult, error)
	Update(cmd venues.UpdateCommand) (*venues.Venue, error)
}

type VenuesControllerConfig interface{}

type VenuesController struct {
	cfg              VenuesControllerConfig
	svc           VenuesService
	validationMapper *validation.ValidationMapper
}

func (c *VenuesController) RegisterRoutes(api *echo.Group) {
	g := api.Group("/organisations/:organisationId/venues")
	g.POST("", c.Create).Name = "v1.organisations.venues.create"
	g.DELETE("/:id", c.Delete).Name = "v1.organisations.venues.delete"
	g.GET("/:id", c.Get).Name = "v1.organisations.venues.get"
	g.PATCH("/:id", c.Patch).Name = "v1.organisations.venues.patch"
	g.GET("", c.List).Name = "v1.organisations.venues.list"
}

func NewVenuesController(
	cfg VenuesControllerConfig,
	validationMapper *validation.ValidationMapper,
	svc VenuesService,
) *VenuesController {
	return &VenuesController{
		cfg:              cfg,
		svc:           svc,
		validationMapper: validationMapper,
	}
}

// @Summary		List all venues for an organisation
// @Tags			Venues
// @Accept			json
// @Produce		json
// @Success		200	{object}	responses.ListVenuesResponse
// @Failure		422	{object}	responses.ValidationErrorResponse
// @Failure		404	{object}	responses.GenericErrorResponse
// @Failure		400	{object}	responses.GenericErrorResponse
// @Failure		500	{object}	responses.GenericErrorResponse
// @Router			/v1/organisations/{orgId}/venues [get]
// @Param			orgId	path	string	true	"The Organisations ID"
// @Param			request	query requests.ListVenuesRequest	true "Query params"
func (c *VenuesController) List(ctx echo.Context) error {
	req := requests.ListVenuesRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	cmd := req.ToCommand()

	results, pagination, err := c.svc.List(cmd)

	if err != nil {
		return err
	}

	resp := responses.ListVenuesResponse{
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
		Data: responses.VenuesFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}

// @Summary		Create a venue for an organisation
// @Tags			Venues
// @Accept			json
// @Produce		json
// @Success		201	{object}	responses.PostVenueResponse
// @Failure		422	{object}	responses.ValidationErrorResponse
// @Failure		404	{object}	responses.GenericErrorResponse
// @Failure		400	{object}	responses.GenericErrorResponse
// @Failure		500	{object}	responses.GenericErrorResponse
// @Router			/v1/organisations/{orgId}/venues [post]
// @Param			orgId	path	string	true	"The Organisations ID"
// @Param			Venue	body		requests.PostVenueRequest	true	"Venue definition"
func (c *VenuesController) Create(ctx echo.Context) error {
	req := requests.PostVenueRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	v, err := c.svc.Create(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}
		return err
	}

	resp := responses.PostVenueResponse{
		Data: responses.VenueFromModel(v),
	}

	ctx.JSON(201, resp)

	return nil
}

// @Summary		Get a venue for an organisation
// @Tags			Venues
// @Accept			json
// @Produce		json
// @Success		200	{object}	responses.PostVenueResponse
// @Failure		422	{object}	responses.ValidationErrorResponse
// @Failure		404	{object}	responses.GenericErrorResponse
// @Failure		400	{object}	responses.GenericErrorResponse
// @Failure		500	{object}	responses.GenericErrorResponse
// @Router			/v1/organisations/{orgId}/venues/{id} [get]
// @Param			orgId	path	string	true	"The Organisations ID"
// @Param			id	path	string	true	"The Venues ID"
func (c *VenuesController) Get(ctx echo.Context) error {
	req := requests.GetVenueRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.svc.Get(req.ToCommand())

	if err != nil {
		return err
	}

	resp := responses.GetVenueResponse{
		Data: responses.VenueFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

// @Summary		Update a venue for an organisation
// @Tags			Venues
// @Accept			json
// @Produce		json
// @Success		200	{object}	responses.PatchVenueResponse
// @Failure		422	{object}	responses.ValidationErrorResponse
// @Failure		404	{object}	responses.GenericErrorResponse
// @Failure		400	{object}	responses.GenericErrorResponse
// @Failure		500	{object}	responses.GenericErrorResponse
// @Router			/v1/organisations/{orgId}/venues/{id} [patch]
// @Param			orgId	path	string	true	"The Organisations ID"
// @Param			id	path	string	true	"The Venues ID"
// @Param			Changes	body		requests.PatchVenueRequest	true	"Venue changes"
func (c *VenuesController) Patch(ctx echo.Context) error {
	req := requests.PatchVenueRequest{}
	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	venue, err := c.svc.Update(req.ToCommand())

	if err != nil {
		return err
	}

	resp := responses.PatchVenueResponse{
		Data: responses.VenueFromModel(venue),
	}

	ctx.JSON(200, resp)

	return nil
}

// @Summary		Delete a venue from an organisation
// @Tags			Venues
// @Accept			json
// @Produce		json
// @Success		204
// @Failure		422	{object}	responses.ValidationErrorResponse
// @Failure		404	{object}	responses.GenericErrorResponse
// @Failure		400	{object}	responses.GenericErrorResponse
// @Failure		500	{object}	responses.GenericErrorResponse
// @Router			/v1/organisations/{orgId}/venues/{id} [delete]
// @Param			orgId	path	string	true	"The Organisations ID"
// @Param			id	path	string	true	"The Venues ID"
func (c *VenuesController) Delete(ctx echo.Context) error {
	req := requests.DeleteVenueRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	err := c.svc.Delete(req.ToCommand())

	if err != nil {
		return err
	}

	ctx.NoContent(204)

	return nil
}
