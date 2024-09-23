package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/util"
	"github.com/adamkirk-stayaway/organisations/pkg/venues"
)

type V1ListVenuesRequest struct {
	OrganisationID *string `param:"organisationId"`
	OrderDirection *string `query:"order_dir"`
	OrderBy *string `query:"order_by"`
	Page *int `query:"page"`
	PerPage *int `query:"per_page"`
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
	Line1 *string `json:"line_1" validationmap:"AddressLine1"`
	Line2 *string `json:"line_2" validationmap:"AddressLine2"`
	Municipality *string `json:"municipality" validationmap:"Municipality"`
	PostCode *string `json:"postcode" validationmap:"PostCode"`
	Lat *float64 `json:"lat" validationmap:"Lat"`
	Long *float64 `json:"long" validationmap:"Long"`
}

type V1PostVenueRequest struct {
	OrganisationID string `param:"organisationId"`
	Name *string `json:"name" validationmap:"Name"`
	Slug *string `json:"slug" validationmap:"Slug"`
	Type *string `json:"type" validationmap:"Type"`
	Address V1PostVenueAddress `json:"address"`
}

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

type V1PatchVenueRequest struct {
	raw map[string]any

	ID string `param:"id"`
	OrganisationID string `param:"organisationId"`
	Name *string `json:"name" validationmap:"Name"`
	Slug *string `json:"slug" validationmap:"Slug"`
	Type *string `json:"type" validationmap:"Type"`
	Address V1PostVenueAddress `json:"address"`
}

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