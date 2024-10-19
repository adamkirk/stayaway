package v1

import (
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/requests"
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/responses"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/accommodations"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/labstack/echo/v4"
)

type VenueAccommodationsService interface {
	Create(cmd accommodations.CreateCommand) (*accommodations.Accommodation, error)
}

type VenueAccommodationsControllerConfig interface{}

type VenueAccommodationsController struct {
	svc              VenueAccommodationsService
	validationMapper *validation.ValidationMapper
}

func (c *VenueAccommodationsController) RegisterRoutes(api *echo.Group) {
	g := api.Group("/organisations/:organisationId/venues/:venueId/accommodations")
	g.POST("", c.Create).Name = "v1.organisations.venues.accommodations.create"
	// g.DELETE("/:id", c.Delete).Name = "v1.organisations.venues.accommodations.delete"
	// g.GET("/:id", c.Get).Name = "v1.organisations.venues.accommodations.get"
	// g.PATCH("/:id", c.Patch).Name = "v1.organisations.venues.accommodations.patch"
	// g.GET("", c.List).Name = "v1.organisations.venues.accommodations.list"
}

func NewVenueAccommodationsController(
	svc VenueAccommodationsService,
	validationMapper *validation.ValidationMapper,
) *VenueAccommodationsController {
	return &VenueAccommodationsController{
		svc:              svc,
		validationMapper: validationMapper,
	}
}

// @Summary		Create an accommodation for a venue
// @Tags			Accommodations
// @Accept			json
// @Produce		json
// @Success		201	{object}	responses.PostVenueAccommodationResponse
// @Failure		422	{object}	responses.ValidationErrorResponse
// @Failure		404	{object}	responses.GenericErrorResponse
// @Failure		400	{object}	responses.GenericErrorResponse
// @Failure		500	{object}	responses.GenericErrorResponse
// @Router			/v1/organisations/{orgId}/venues/{venueId}/accommodations [post]
// @Param			orgId	path	string	true	"The Organisations ID"
// @Param			venueId	path	string	true	"The Venues ID"
// @Param			body	body		requests.PostVenueAccommodationRequest	true	"Venue Accommodation definition"
func (c *VenueAccommodationsController) Create(ctx echo.Context) error {
	req := requests.PostVenueAccommodationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	a, err := c.svc.Create(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}
		return err
	}

	resp := responses.PostVenueAccommodationResponse{
		Data: responses.VenueAccommodationFromModel(a),
	}

	ctx.JSON(201, resp)

	return nil
}

// // @Summary		Get a venue accommodation template
// // @Tags			AccommodationTemplates
// // @Accept			json
// // @Produce		json
// // @Success		200	{object}	responses.GetVenueAccommodationTemplateResponse
// // @Failure		422	{object}	responses.ValidationErrorResponse
// // @Failure		404	{object}	responses.GenericErrorResponse
// // @Failure		400	{object}	responses.GenericErrorResponse
// // @Failure		500	{object}	responses.GenericErrorResponse
// // @Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates/{id} [get]
// // @Param			orgId	path	string	true	"The Organisations ID"
// // @Param			venueId	path	string	true	"The Venues ID"
// // @Param			id	path	string	true	"The ID of the template"
// func (c *VenueAccommodationsController) Get(ctx echo.Context) error {
// 	req := requests.GetVenueAccommodationTemplateRequest{}

// 	if err := bindRequest(&req, ctx); err != nil {
// 		return err
// 	}

// 	vt, err := c.svc.Get(req.ToCommand())

// 	if err != nil {
// 		return err
// 	}

// 	resp := responses.GetVenueAccommodationTemplateResponse{
// 		Data: responses.VenueAccommodationTemplateFromModel(vt),
// 	}

// 	ctx.JSON(200, resp)

// 	return nil
// }

// // @Summary		Update a venue accommodation template
// // @Tags			AccommodationTemplates
// // @Accept			json
// // @Produce		json
// // @Success		200	{object}	responses.PatchVenueAccommodationTemplateResponse
// // @Failure		422	{object}	responses.ValidationErrorResponse
// // @Failure		404	{object}	responses.GenericErrorResponse
// // @Failure		400	{object}	responses.GenericErrorResponse
// // @Failure		500	{object}	responses.GenericErrorResponse
// // @Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates/{id} [patch]
// // @Param			orgId	path	string	true	"The Organisations ID"
// // @Param			venueId	path	string	true	"The Venues ID"
// // @Param			id	path	string	true	"The ID of the template"
// // @Param			Changes	body		requests.PatchVenueAccommodationTemplateRequest	true	"Venue changes"
// func (c *VenueAccommodationsController) Patch(ctx echo.Context) error {
// 	req := requests.PatchVenueAccommodationTemplateRequest{}

// 	if err := bindRequest(&req, ctx); err != nil {
// 		return err
// 	}

// 	vt, err := c.svc.Update(req.ToCommand())

// 	if err != nil {
// 		if err, ok := err.(validation.ValidationError); ok {
// 			return c.validationMapper.Map(err, req)
// 		}
// 		return err
// 	}

// 	resp := responses.PatchVenueAccommodationTemplateResponse{
// 		Data: responses.VenueAccommodationTemplateFromModel(vt),
// 	}

// 	ctx.JSON(200, resp)

// 	return nil
// }

// // @Summary		List all accommodation templates for a venue
// // @Tags			AccommodationTemplates
// // @Accept			json
// // @Produce		json
// // @Success		200	{object}	responses.ListVenueAccommodationsResponse
// // @Failure		422	{object}	responses.ValidationErrorResponse
// // @Failure		404	{object}	responses.GenericErrorResponse
// // @Failure		400	{object}	responses.GenericErrorResponse
// // @Failure		500	{object}	responses.GenericErrorResponse
// // @Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates [get]
// // @Param			orgId	path	string	true	"The Organisations ID"
// // @Param			venueId	path	string	true	"The Venues ID"
// // @Param			request	query requests.ListVenueAccommodationsRequest	true "Query params"
// func (c *VenueAccommodationsController) List(ctx echo.Context) error {
// 	req := requests.ListVenueAccommodationsRequest{}

// 	if err := bindRequest(&req, ctx); err != nil {
// 		return err
// 	}

// 	cmd := req.ToCommand()

// 	results, pagination, err := c.svc.List(cmd)

// 	if err != nil {
// 		if err, ok := err.(validation.ValidationError); ok {
// 			return c.validationMapper.Map(err, req)
// 		}
// 		return err
// 	}

// 	resp := responses.ListVenueAccommodationsResponse{
// 		Meta: responses.ListResponseMeta{
// 			SortOptionsResponseMeta: responses.SortOptionsResponseMeta{
// 				OrderDirection: string(cmd.OrderDirection),
// 				OrderBy:        string(cmd.OrderBy),
// 			},
// 			PaginationResponseMeta: responses.PaginationResponseMeta{
// 				Page:         pagination.Page,
// 				PerPage:      pagination.PerPage,
// 				TotalPages:   pagination.TotalPages,
// 				TotalResults: pagination.Total,
// 			},
// 		},
// 		Data: responses.VenueAccommodationsFromModels(results),
// 	}

// 	ctx.JSON(200, resp)

// 	return nil
// }

// // @Summary		Delete a venue accommodation template
// // @Tags			AccommodationTemplates
// // @Accept			json
// // @Produce		json
// // @Success		204
// // @Failure		422	{object}	responses.ValidationErrorResponse
// // @Failure		404	{object}	responses.GenericErrorResponse
// // @Failure		400	{object}	responses.GenericErrorResponse
// // @Failure		500	{object}	responses.GenericErrorResponse
// // @Router			/v1/organisations/{orgId}/venues/{venueId}/accommodation-templates/{id} [delete]
// // @Param			orgId	path	string	true	"The Organisations ID"
// // @Param			venueId	path	string	true	"The Venues ID"
// // @Param			id	path	string	true	"The template ID"
// func (c *VenueAccommodationsController) Delete(ctx echo.Context) error {
// 	req := requests.DeleteVenueAccommodationTemplateRequest{}

// 	if err := bindRequest(&req, ctx); err != nil {
// 		return err
// 	}

// 	err := c.svc.Delete(req.ToCommand())

// 	if err != nil {
// 		return err
// 	}

// 	ctx.NoContent(204)

// 	return nil
// }
