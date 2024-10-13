package v1

import (
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/requests"
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/responses"
	"github.com/adamkirk-stayaway/organisations/internal/domain/accommodations"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
	"github.com/labstack/echo/v4"
)

type VenueAccommodationTemplateCreateHandler interface {
	Handle(cmd accommodations.CreateVenueTemplateCommand) (*accommodations.VenueTemplate, error)
}

type VenueAccommodationTemplateGetHandler interface {
	Handle(cmd accommodations.GetVenueTemplateCommand) (*accommodations.VenueTemplate, error)
}

type VenueAccommodationTemplatesListHandler interface {
	Handle(cmd accommodations.ListVenueTemplatesCommand) (accommodations.VenueTemplates, common.PaginationResult, error)
}

type VenueAccommodationTemplateDeleteHandler interface {
	Handle(cmd accommodations.DeleteVenueTemplateCommand) (error)
}

type VenueAccommodationTemplateUpdateHandler interface {
	Handle(cmd accommodations.UpdateVenueTemplateCommand) (*accommodations.VenueTemplate, error)
}


type VenueAccommodationTemplatesControllerConfig interface {}

type VenueAccommodationTemplatesController struct {
	cfg VenuesControllerConfig
	create VenueAccommodationTemplateCreateHandler
	get VenueAccommodationTemplateGetHandler
	delete VenueAccommodationTemplateDeleteHandler
	update VenueAccommodationTemplateUpdateHandler
	list VenueAccommodationTemplatesListHandler
	validationMapper *validation.ValidationMapper
}

func (c *VenueAccommodationTemplatesController) RegisterRoutes(api *echo.Group) {
	g := api.Group("/organisations/:organisationId/venues/:venueId/accommodation-templates")
	g.POST("", c.Create).Name = "v1.organisations.venues.accommodation_templates.create"
	g.DELETE("/:id", c.Delete).Name = "v1.organisations.venues.accommodation_templates.delete"
	g.GET("/:id", c.Get).Name = "v1.organisations.venues.accommodation_templates.get"
	g.PATCH("/:id", c.Patch).Name = "v1.organisations.venues.accommodation_templates.patch"
	g.GET("", c.List).Name = "v1.organisations.venues.accommodation_templates.list"
}

func NewVenueAccommodationTemplatesController(
	cfg VenuesControllerConfig,
	create VenueAccommodationTemplateCreateHandler,
	get VenueAccommodationTemplateGetHandler,
	list VenueAccommodationTemplatesListHandler,
	delete VenueAccommodationTemplateDeleteHandler,
	update VenueAccommodationTemplateUpdateHandler,
	validationMapper *validation.ValidationMapper,
) *VenueAccommodationTemplatesController {
	return &VenueAccommodationTemplatesController{
		cfg: cfg,
		create: create,
		list: list,
		get: get,
		delete: delete,
		update: update,
		validationMapper: validationMapper,
	}
}

//	@Summary		Create a venue accommodation template for an organisation
//	@Tags			AccommodationTemplates
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	responses.PostVenueAccommodationTemplateResponse
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates [post]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			venueId	path	string	true	"The Venues ID"
//	@Param			VenueAccommodationTemplate	body		requests.PostVenueAccommodationTemplateRequest	true	"Venue Accommodation Template definition"
func (c *VenueAccommodationTemplatesController) Create(ctx echo.Context) error {
	req := requests.PostVenueAccommodationTemplateRequest{}

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

	resp := responses.PostVenueAccommodationTemplateResponse{
		Data: responses.VenueAccommodationTemplateFromModel(v),
	}

	ctx.JSON(201, resp)

	return nil
}

//	@Summary		Get a venue accommodation template
//	@Tags			AccommodationTemplates
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.GetVenueAccommodationTemplateResponse
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates/{id} [get]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			venueId	path	string	true	"The Venues ID"
//	@Param			id	path	string	true	"The ID of the template"
func (c *VenueAccommodationTemplatesController) Get(ctx echo.Context) error {
	req := requests.GetVenueAccommodationTemplateRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	vt, err := c.get.Handle(req.ToCommand())

	if err != nil {
		return err
	}

	resp := responses.GetVenueAccommodationTemplateResponse{
		Data: responses.VenueAccommodationTemplateFromModel(vt),
	}

	ctx.JSON(200, resp)

	return nil
}


//	@Summary		Update a venue accommodation template
//	@Tags			AccommodationTemplates
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.PatchVenueAccommodationTemplateResponse
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates/{id} [patch]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			venueId	path	string	true	"The Venues ID"
//	@Param			id	path	string	true	"The ID of the template"
//	@Param			Changes	body		requests.PatchVenueAccommodationTemplateRequest	true	"Venue changes"
func (c *VenueAccommodationTemplatesController) Patch(ctx echo.Context) error {
	req := requests.PatchVenueAccommodationTemplateRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	vt, err := c.update.Handle(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}
		return err
	}

	resp := responses.PatchVenueAccommodationTemplateResponse{
		Data: responses.VenueAccommodationTemplateFromModel(vt),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		List all accommodation templates for a venue
//	@Tags			AccommodationTemplates
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.ListVenueAccommodationTemplatesResponse
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates [get]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			venueId	path	string	true	"The Venues ID"
//	@Param			request	query requests.ListVenueAccommodationTemplatesRequest	true "Query params"
func (c *VenueAccommodationTemplatesController) List(ctx echo.Context) error {
	req := requests.ListVenueAccommodationTemplatesRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	cmd := req.ToCommand()

	results, pagination, err := c.list.Handle(cmd)

	if err != nil {
		return err
	}
	
	resp := responses.ListVenueAccommodationTemplatesResponse{
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
		Data: responses.VenueAccommodationTemplatesFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}

//	@Summary		Delete a venue accommodation template
//	@Tags			AccommodationTemplates
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		422	{object}	responses.ValidationErrorResponse
//	@Failure		404	{object}	responses.GenericErrorResponse
//	@Failure		400	{object}	responses.GenericErrorResponse
//	@Failure		500	{object}	responses.GenericErrorResponse
//	@Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates/{id} [delete]
//	@Param			orgId	path	string	true	"The Organisations ID"
//	@Param			venueId	path	string	true	"The Venues ID"
//	@Param			id	path	string	true	"The template ID"
func (c *VenueAccommodationTemplatesController) Delete(ctx echo.Context) error {
	req := requests.DeleteVenueAccommodationTemplateRequest{}

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