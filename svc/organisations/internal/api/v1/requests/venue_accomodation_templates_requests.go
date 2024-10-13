package requests

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/accommodations"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/util"
)


type ListVenueAccommodationTemplatesRequest struct {
	// The direction to order the results by.
	OrderDirection *string `query:"order_dir" json:"order_dir" validationmap:"OrderDirection" validate:"optional" enums:"asc,desc"`

	// The field by which to order the results.
	OrderBy *string `query:"order_by" json:"order_by" validationmap:"OrderBy" validate:"optional" enums:"name"`

	// The page to display.
	// An empty list may be returned if going beyond the last page of results.
	Page *int `query:"page" json:"page" validationmap:"Page" validate:"optional" minimum:"1"`

	// The amount of results to display per page.
	PerPage *int `query:"per_page" json:"per_page" validationmap:"PerPage" validate:"optional" minimum:"1" maximum:"100"`

	// Characters to use as a prefix in searching by the name.
	// Useful for a "typeahead" widget.
	NamePrefix *string `query:"prefix" json:"prefix" validationmap:"NamePrefix" validate:"optional" minimum:"3"`

	OrganisationID *string `param:"organisationId" swaggerignore:"true"`
	VenueID *string `param:"venueId" swaggerignore:"true"`
}

func (req ListVenueAccommodationTemplatesRequest) ToCommand() accommodations.ListVenueTemplatesCommand {
	cmd := accommodations.NewListVenueTemplatesCommand()

	if req.OrderDirection != nil {
		cmd.OrderDirection = common.SortDirection(*req.OrderDirection)
	}

	if req.OrderBy != nil {
		cmd.OrderBy = accommodations.SortBy(*req.OrderBy)
	}

	if req.Page != nil {
		cmd.Page = *req.Page
	}

	if req.PerPage != nil {
		cmd.PerPage = *req.PerPage
	}

	cmd.NamePrefix = req.NamePrefix

	cmd.OrganisationID = *req.OrganisationID
	cmd.VenueID = *req.VenueID

	return cmd
}

type PostVenueAccommodationTemplateOccupancy struct {
	// The minimum amount of people that must occupy the accommodation.
	Min *int `json:"min" validationmap:"MinOccupancy" validate:"required" minimum:"1"`

	// The maximum amount of people that must occupy the accommodation.
	// If null or blank there will be no limit.
	// If provided, it must be greater than the min occupancy.
	Max *int `json:"max" validationmap:"MaxOccupancy" validate:"optional" minimum:"1" extensions:"x-nullable"`
} // @name	V1.Request[Model].VenueAccommodationTemplateOccupancyCreate

type PostVenueAccommodationTemplateRequest struct {
	OrganisationID string `param:"organisationId" swaggerignore:"true"`
	VenueID string `param:"venueId" swaggerignore:"true"`

	// The name of the template.
	Name *string `json:"name" validationmap:"Name" validate:"required" minLength:"3"`

	// The type of accommodation that this template is for.
	// Currently only supports 'room'
	Type *string `json:"type" validationmap:"Type" validate:"required" enums:"room"`

	// Description of the accommodation that this applies to.
	Description *string `json:"description" validationmap:"Description" validate:"required" minimum:"10"`

	Occupancy PostVenueAccommodationTemplateOccupancy `json:"occupancy" validate:"required"`
} // @name	V1.Request.CreateAccommodationVenueTemplate

func (req PostVenueAccommodationTemplateRequest) ToCommand() accommodations.CreateVenueTemplateCommand {
	return accommodations.CreateVenueTemplateCommand{
		OrganisationID: &req.OrganisationID,
		VenueID: &req.VenueID,
		Name: req.Name,
		Type: req.Type,
		Description: req.Description,
		MinOccupancy: req.Occupancy.Min,
		MaxOccupancy: req.Occupancy.Max,
	}
}


type PatchVenueAccommodationTemplateOccupancy struct {
	// The minimum amount of people that must occupy the accommodation.
	Min *int `json:"min" validationmap:"MinOccupancy" validate:"optional" minimum:"1" extensions:"x-nullable"`

	// The maximum amount of people that must occupy the accommodation.
	// If null or blank there will be no limit.
	// If provided, it must be greater than the min occupancy.
	Max *int `json:"max" validationmap:"MaxOccupancy" validate:"optional" minimum:"1" extensions:"x-nullable"`
} // @name	V1.Request[Model].VenueAccommodationTemplateOccupancyUpdate

type PatchVenueAccommodationTemplateRequest struct {
	raw map[string]any

	OrganisationID string `param:"organisationId" swaggerignore:"true"`
	VenueID string `param:"venueId" swaggerignore:"true"`
	ID string `param:"id" swaggerignore:"true"`

	// The name of the template.
	Name *string `json:"name" validationmap:"Name" validate:"optional" minLength:"3" extensions:"x-nullable"`

	// The type of accommodation that this template is for.
	// Currently only supports 'room'
	Type *string `json:"type" validationmap:"Type" validate:"optional" enums:"room" extensions:"x-nullable"`

	// Description of the accommodation that this applies to.
	Description *string `json:"description" validationmap:"Description" validate:"optional" minLength:"10" extensions:"x-nullable"` 

	Occupancy PatchVenueAccommodationTemplateOccupancy `json:"occupancy" validate:"optional" minLength:"10" extensions:"x-nullable"`
} // @name	V1.Request.UpdateAccommodationVenueTemplate

func (req *PatchVenueAccommodationTemplateRequest) IncludeRawBody(raw map[string]any) {
	req.raw = raw
}

func (req *PatchVenueAccommodationTemplateRequest) FieldWasPresent(fld string) bool {
	return util.KeyExistsInMap(req.raw, fld)
}

func (req PatchVenueAccommodationTemplateRequest) ToCommand() accommodations.UpdateVenueTemplateCommand {
	return accommodations.UpdateVenueTemplateCommand{
		OrganisationID: req.OrganisationID,
		VenueID: req.VenueID,
		ID: req.ID,
		Name: req.Name,
		Type: req.Type,
		Description: req.Description,
		MinOccupancy: req.Occupancy.Min,
		MaxOccupancy: req.Occupancy.Max,
		NullifyMaxOccupancy: req.FieldWasPresent("occupancy.max") && req.Occupancy.Max == nil,
	}
}

type GetVenueAccommodationTemplateRequest struct {
	ID string `param:"id"`
	OrganisationID string `param:"organisationId"`
	VenueID string `param:"venueId"`
}

func (req GetVenueAccommodationTemplateRequest) ToCommand() accommodations.GetVenueTemplateCommand {
	return accommodations.GetVenueTemplateCommand{
		ID: req.ID,
		OrganisationID: req.OrganisationID,
		VenueID: req.VenueID,
	}
}

type DeleteVenueAccommodationTemplateRequest struct {
	ID string `param:"id"`
	OrganisationID string `param:"organisationId"`
	VenueID string `param:"venueId"`
}

func (req DeleteVenueAccommodationTemplateRequest) ToCommand() accommodations.DeleteVenueTemplateCommand {
	return accommodations.DeleteVenueTemplateCommand{
		ID: req.ID,
		OrganisationID: req.OrganisationID,
		VenueID: req.VenueID,
	}
}

