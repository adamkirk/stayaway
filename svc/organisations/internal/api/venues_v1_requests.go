package api

import (
	"github.com/adamkirk-stayaway/organisations/internal/model"
	"github.com/adamkirk-stayaway/organisations/internal/util"
	"github.com/adamkirk-stayaway/organisations/internal/venues"
)

type V1ListVenuesRequest struct {
	OrganisationID *string `param:"organisationId" swaggerignore:"true"`
	OrderDirection *string `query:"order_dir" json:"order_dir"`
	OrderBy *string `query:"order_by" json:"order_by"`
	Page *int `query:"page" json:"page"`
	PerPage *int `query:"per_page" json:"per_page"`
}

func (req V1ListVenuesRequest) ToCommand() venues.ListCommand {
	cmd := venues.NewListCommand()

	if req.OrderDirection != nil {
		cmd.OrderDirection = model.SortDirection(*req.OrderDirection)
	}

	if req.OrderBy != nil {
		cmd.OrderBy = model.VenueSortBy(*req.OrderBy)
	}

	if req.Page != nil {
		cmd.Page = *req.Page
	}

	if req.PerPage != nil {
		cmd.PerPage = *req.PerPage
	}

	cmd.OrganisationID = *req.OrganisationID

	return cmd
}

type V1PostVenueAddress struct {
	// Line 1 of the address, typically number/name and street.
	Line1 *string `json:"line_1" validationmap:"AddressLine1" validate:"required" minLength:"1"` 

	// Line 2, extra information for the address if needed, optional.
	Line2 *string `json:"line_2" validationmap:"AddressLine2" validate:"optional" minLength:"1" extensions:"x-nullable"`

	// The town/city/village that the venue is in.
	Municipality *string `json:"municipality" validationmap:"Municipality" validate:"required" minLength:"1"`

	// A valid UK postcode, following standard formats.
	PostCode *string `json:"postcode" validationmap:"PostCode"  validate:"required"`

	//Latitude of the venue.
	Lat *float64 `json:"lat" validationmap:"Lat" validate:"required" minimum:"0"`

	//Longitude of the venue.
	Long *float64 `json:"long" validationmap:"Long" validate:"required" minimum:"0"`
} // @name	V1.Request[Model].VenueAddress

type V1PostVenueRequest struct {
	OrganisationID string `param:"organisationId" swaggerignore:"true"`

	// The name of the venue.
	Name *string `json:"name" validationmap:"Name" validate:"required" minLength:"3"`

	// The slug of the venue, used in URI's.
	// Must only contain alphanumeric and hyphen characters.
	// Must be unique within the organisation.
	Slug *string `json:"slug" validationmap:"Slug" validate:"required" minLength:"3"`

	// The type of venue.
	// Currently only supports 'hotel'
	Type *string `json:"type" validationmap:"Type" validate:"required" enums:"hotel"`

	// The address of the venue.
	Address V1PostVenueAddress `json:"address" validate:"required"`
} // @name	V1.Request.CreateVenue

func (req V1PostVenueRequest) ToCommand() venues.CreateCommand {
	return venues.CreateCommand{
		OrganisationID: &req.OrganisationID,
		Name: req.Name,
		Slug: req.Slug,
		Type: req.Type,
		AddressLine1: req.Address.Line1,
		AddressLine2: req.Address.Line2,
		PostCode: req.Address.PostCode,
		Municipality: req.Address.Municipality,
		Lat: req.Address.Lat,
		Long: req.Address.Long,
	}
}


type V1DeleteVenueRequest struct {
	ID string `param:"id"`
	OrganisationID string `param:"organisationId"`
}

func (req V1DeleteVenueRequest) ToCommand() venues.DeleteCommand {
	return venues.DeleteCommand{
		ID: req.ID,
		OrganisationID: req.OrganisationID,
	}
}
type V1GetVenueRequest struct {
	ID string `param:"id"`
	OrganisationID string `param:"organisationId"`
}

func (req V1GetVenueRequest) ToCommand() venues.GetCommand {
	return venues.GetCommand{
		ID: req.ID,
		OrganisationID: req.OrganisationID,
	}
}

// @Description The changes to apply to the venue. Only include fields to change.
type V1PatchVenueAddress struct {
	// Line 1 of the address, typically number/name and street.
	Line1 *string `json:"line_1" validationmap:"AddressLine1" validate:"optional" minLength:"1" extensions:"x-nullable"` 

	// Line 2, extra information for the address if needed, optional.
	Line2 *string `json:"line_2" validationmap:"AddressLine2" validate:"optional" minLength:"1" extensions:"x-nullable"`
	
	// The town/city/village that the venue is in.
	Municipality *string `json:"municipality" validationmap:"Municipality" validate:"optional" minLength:"1" extensions:"x-nullable"`
	
	// A valid UK postcode, following standard formats.
	PostCode *string `json:"postcode" validationmap:"PostCode" validate:"optional" extensions:"x-nullable"`

	//Latitude of the venue.
	Lat *float64 `json:"lat" validationmap:"Lat" validate:"optional" minimum:"0" extensions:"x-nullable"`
	
	//Longitude of the venue.
	Long *float64 `json:"long" validationmap:"Long" validate:"optional" minimum:"0" extensions:"x-nullable"`
} // @name	V1.Request[Model].VenueAddress

type V1PatchVenueRequest struct {
	raw map[string]any

	ID string `param:"id" swaggerignore:"true"`
	OrganisationID string `param:"organisationId" swaggerignore:"true"`


	// The name of the venue.
	Name *string `json:"name" validationmap:"Name" validate:"optional" minLength:"3" extensions:"x-nullable"`

	// The slug of the venue, used in URI's.
	// Must only contain alphanumeric and hyphen characters.
	// Must be unique within the organisation.
	Slug *string `json:"slug" validationmap:"Slug" validate:"optional" minLength:"3" extensions:"x-nullable"`
	
	// The type of venue.
	// Currently only supports 'hotel'
	Type *string `json:"type" validationmap:"Type" validate:"optional" enums:"hotel" extensions:"x-nullable"`

	// The address of the venue.
	Address V1PatchVenueAddress `json:"address" validate:"optional" extensions:"x-nullable"`
} // @name	V1.Request.UpdateVenue

func (req *V1PatchVenueRequest) IncludeRawBody(raw map[string]any) {
	req.raw = raw
}

func (req *V1PatchVenueRequest) FieldWasPresent(fld string) bool {
	return util.KeyExistsInMap(req.raw, fld)
}

func (req V1PatchVenueRequest) ToCommand() venues.UpdateCommand {
	return venues.UpdateCommand{
		ID: &req.ID,
		OrganisationID: &req.OrganisationID,
		Name: req.Name,
		Slug: req.Slug,
		Type: req.Slug,
		AddressLine1: req.Address.Line1,
		AddressLine2: req.Address.Line2,
		NullifyAddressLine2: req.FieldWasPresent("address.line_2") && req.Address.Line2 == nil,
		Municipality: req.Address.Municipality,
		PostCode: req.Address.PostCode,
		Lat: req.Address.Lat,
		Long: req.Address.Long,
	}
}