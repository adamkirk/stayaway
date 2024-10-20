package requests

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/accommodations"
	"github.com/adamkirk-stayaway/organisations/internal/util"
)

type PostVenueAccommodationOccupancy struct {
	// The minimum amount of people that must occupy the accommodation.
	// Or null if venue_template_id is supplied.
	Min *int `json:"min" validationmap:"MinOccupancy" validate:"required" minimum:"1" extensions:"x-nullable"`

	// The maximum amount of people that must occupy the accommodation.
	// If null or blank there will be no limit, or it will use the templates value.
	// If provided, it must be greater than the min occupancy.
	Max *int `json:"max" validationmap:"MaxOccupancy" validate:"optional" minimum:"1" extensions:"x-nullable"`
} // @name	V1.Request[Model].VenueAccommodationOccupancyCreate

type PostVenueAccommodationRequest struct {
	VenueID        string `param:"venueId" swaggerignore:"true"`

	VenueTemplateID        *string `json:"venue_template_id" validationmap:"VenueTemplateID"`

	// The name of the template.
	// Or null if venue_template_id is supplied.
	Name *string `json:"name" validationmap:"Name" validate:"required" extensions:"x-nullable" minLength:"3"`

	// A reference for this accommodation.
	Reference *string `json:"reference" validationmap:"Reference" validate:"required" minLength:"3"`

	// The type of accommodation that this template is for.
	// Currently only supports 'room'.
	// Or null if venue_template_id is supplied.
	Type *string `json:"type" validationmap:"Type" validate:"required" extensions:"x-nullable" enums:"room"`

	// Description of the accommodation that this applies to.
	// Or null if venue_template_id is supplied.
	Description *string `json:"description" validationmap:"Description" validate:"required" minimum:"10"`

	Occupancy PostVenueAccommodationOccupancy `json:"occupancy" validate:"required"`
} // @name	V1.Request.CreateVenueAccommodation

func (req PostVenueAccommodationRequest) ToCommand() accommodations.CreateCommand {
	return accommodations.CreateCommand{
		VenueID: &req.VenueID,
		VenueTemplateID: req.VenueTemplateID,
		Reference: req.Reference,
		Name: req.Name,
		Description: req.Description,
		Type: req.Type,
		MinOccupancy: req.Occupancy.Min,
		MaxOccupancy: req.Occupancy.Max,
	}
}

type GetVenueAccommodationRequest struct {
	ID             string `param:"id"`
	OrganisationID string `param:"organisationId"`
	VenueID        string `param:"venueId"`
} // @name	V1.Request.GetVenueAccommodation

func (req GetVenueAccommodationRequest) ToCommand() accommodations.GetCommand {
	return accommodations.GetCommand{
		ID:             &req.ID,
		VenueID:        &req.VenueID,
	}
}

type ListVenueAccommodationsRequest struct {
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
	ReferencePrefix *string `query:"prefix" json:"prefix" validationmap:"ReferencePrefix" validate:"optional" minimum:"3"`

	VenueTemplateID        *string `query:"venue_template_id" json:"venue_template_id" validate:"optional" validationmap:"VenueTemplateID"`

	OrganisationID *string `param:"organisationId" swaggerignore:"true"`
	VenueID        *string `param:"venueId" swaggerignore:"true"`
}

func (req ListVenueAccommodationsRequest) ToCommand() accommodations.ListCommand {
	cmd := accommodations.NewListCommand()

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

	cmd.ReferencePrefix = req.ReferencePrefix
	cmd.VenueTemplateID = req.VenueTemplateID

	cmd.OrganisationID = *req.OrganisationID
	cmd.VenueID = *req.VenueID

	return cmd
}

type PatchVenueAccommodationOccupancy struct {
	// The minimum amount of people that must occupy the accommodation.
	// Or null to remove this setting and use the templates value.
	Min *int `json:"min" validationmap:"MinOccupancy" validate:"optional" minimum:"1" extensions:"x-nullable"`

	// The maximum amount of people that must occupy the accommodation.
	// If null or blank there will be no limit, or it will use the templates value.
	Max *int `json:"max" validationmap:"MaxOccupancy" validate:"optional" minimum:"1" extensions:"x-nullable"`
} // @name	V1.Request[Model].VenueAccommodationOccupancyUpdate

type PatchVenueAccommodationRequest struct {
	raw map[string]any

	ID string `param:"id" swaggerignore:"true"`
	VenueID        string `param:"venueId" swaggerignore:"true"`

	VenueTemplateID        *string `json:"venue_template_id" validate:"optional" extensions:"x-nullable" validationmap:"VenueTemplateID"`

	// The name of the template.
	// Or null to remove this setting and use the templates value.
	Name *string `json:"name" validationmap:"Name" validate:"optional" extensions:"x-nullable" minLength:"3"`

	// A reference for this accommodation.
	Reference *string `json:"reference" validationmap:"Reference" validate:"optional" minLength:"3"`

	// The type of accommodation that this template is for.
	// Currently only supports 'room'.
	// Or null to remove this setting and use the templates value.
	Type *string `json:"type" validationmap:"Type" validate:"optional" extensions:"x-nullable" enums:"room" extensions:"x-nullable"`

	// Description of the accommodation that this applies to.
	// Or null to remove this setting and use the templates value.
	Description *string `json:"description" validationmap:"Description" validate:"optional" minimum:"10" extensions:"x-nullable"`

	Occupancy PatchVenueAccommodationOccupancy `json:"occupancy" validate:"optional"`
} // @name	V1.Request.PatchVenueAccommodation

func (req PatchVenueAccommodationRequest) ToCommand() accommodations.UpdateCommand {
	return accommodations.UpdateCommand{
		ID: &req.ID,
		VenueID: &req.VenueID,
		Reference: req.Reference,

		VenueTemplateID: req.VenueTemplateID,
		NullifyVenueTemplateID: req.FieldWasPresent("venue_template_id") && req.VenueTemplateID == nil,

		Name: req.Name,
		NullifyName: req.FieldWasPresent("name") && req.Name == nil,

		Description: req.Description,
		NullifyDescription: req.FieldWasPresent("description") && req.Description == nil,

		Type: req.Type,
		NullifyType: req.FieldWasPresent("type") && req.Type == nil,

		MinOccupancy: req.Occupancy.Min,
		NullifyMinOccupancy: req.FieldWasPresent("occupancy.min") && req.Occupancy.Min == nil,

		MaxOccupancy: req.Occupancy.Max,
		NullifyMaxOccupancy: req.FieldWasPresent("occupancy.max") && req.Occupancy.Max == nil,
	}
}


func (req *PatchVenueAccommodationRequest) IncludeRawBody(raw map[string]any) {
	req.raw = raw
}

func (req *PatchVenueAccommodationRequest) FieldWasPresent(fld string) bool {
	return util.KeyExistsInMap(req.raw, fld)
}

type DeleteVenueAccommodationRequest struct {
	ID             string `param:"id"`
	OrganisationID string `param:"organisationId"`
	VenueID        string `param:"venueId"`
} // @name	V1.Request.DeleteVenueAccommodation

func (req DeleteVenueAccommodationRequest) ToCommand() accommodations.DeleteCommand {
	return accommodations.DeleteCommand{
		ID:             &req.ID,
		VenueID:        &req.VenueID,
	}
}