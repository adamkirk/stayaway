package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/organisations"
)


type V1ListOrganisationsRequest struct {
	// The direction to order the results by.
	OrderDirection *string `query:"order_dir" json:"order_dir" validate:"optional" enums:"asc,desc"`
	
	// The field by which to order the results.
	OrderBy *string `query:"order_by" json:"order_by" validate:"optional" enums:"name,slug"`

	// The page to display.
	// An empty list may be returned if going beyond the last page of results.
	Page *int `query:"page" json:"page" validate:"optional" minimum:"1"`

	// The amount of results to display per page.
	PerPage *int `query:"per_page" json:"per_page" validate:"optional" minimum:"1" maximum:"100"`
}

func (req V1ListOrganisationsRequest) ToCommand() organisations.ListCommand {
	cmd := organisations.NewListCommand()

	if req.OrderDirection != nil {
		cmd.OrderDirection = model.SortDirection(*req.OrderDirection)
	}

	if req.OrderBy != nil {
		cmd.OrderBy = model.OrganisationSortBy(*req.OrderBy)
	}

	if req.Page != nil {
		cmd.Page = *req.Page
	}

	if req.PerPage != nil {
		cmd.PerPage = *req.PerPage
	}

	return cmd
}

type V1PostOrganisationRequest struct {
	// The name of the organisation.
	Name *string `json:"name" validationmap:"Name" validate:"required" minLength:"3"`
	
	// A slug that will be used in URI's.
	// Must only contain alphanumeric and hyphen characters
	Slug *string `json:"slug" validationmap:"Slug" validate:"required" minLength:"3"`
} // @name	V1.Request.CreateOrganisation

func (req V1PostOrganisationRequest) ToCommand() organisations.CreateCommand {
	return organisations.CreateCommand{
		Name: req.Name,
		Slug: req.Slug,
	}
}

type V1DeleteOrganisationRequest struct {
	ID string `param:"id"`
}

func (req V1DeleteOrganisationRequest) ToCommand() organisations.DeleteCommand {
	return organisations.DeleteCommand{
		ID: req.ID,
	}
}

type V1GetOrganisationRequest struct {
	ID string `param:"id"`
}

func (req V1GetOrganisationRequest) ToCommand() organisations.GetCommand {
	return organisations.GetCommand{
		ID: req.ID,
	}
}

type V1PatchOrganisationRequest struct {
	ID string `param:"id" swaggerignore:"true"`

	// The name of the organisation.
	Name *string `json:"name,omitempty" validationmap:"Name" validate:"optional" minLength:"3" extensions:"x-nullable"`

	// A slug that will be used in URI's.
	// Must only contain alphanumeric and hyphen characters.
	// Must be unique across all other organisations.
	Slug *string `json:"slug,omitempty" validationmap:"Slug" validate:"optional" minLength:"3" extensions:"x-nullable"`
} // @name	V1.Request.UpdateOrganisation

func (req V1PatchOrganisationRequest) ToCommand() organisations.UpdateCommand {
	return organisations.UpdateCommand{
		ID: req.ID,
		Name: req.Name,
		Slug: req.Slug,
	}
}