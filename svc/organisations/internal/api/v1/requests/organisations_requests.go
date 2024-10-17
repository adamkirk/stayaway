package requests

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/organisations"
)

type ListOrganisationsRequest struct {
	// The direction to order the results by.
	OrderDirection *string `query:"order_dir" json:"order_dir" validationmap:"OrderDir" validate:"optional" enums:"asc,desc"`

	// The field by which to order the results.
	OrderBy *string `query:"order_by" json:"order_by" validationmap:"OrderBy" validate:"optional" enums:"name,slug"`

	// The page to display.
	// An empty list may be returned if going beyond the last page of results.
	Page *int `query:"page" json:"page" validationmap:"Page" validate:"optional" minimum:"1"`

	// The amount of results to display per page.
	PerPage *int `query:"per_page" json:"per_page" validationmap:"PerPage" validate:"optional" minimum:"1" maximum:"100"`
}

func (req ListOrganisationsRequest) ToCommand() organisations.ListCommand {
	cmd := organisations.NewListCommand()

	if req.OrderDirection != nil {
		cmd.OrderDirection = common.SortDirection(*req.OrderDirection)
	}

	if req.OrderBy != nil {
		cmd.OrderBy = organisations.SortBy(*req.OrderBy)
	}

	if req.Page != nil {
		cmd.Page = *req.Page
	}

	if req.PerPage != nil {
		cmd.PerPage = *req.PerPage
	}

	return cmd
}

type PostOrganisationRequest struct {
	// The name of the organisation.
	Name *string `json:"name" validationmap:"Name" validate:"required" minLength:"3"`

	// A slug that will be used in URI's.
	// Must only contain alphanumeric and hyphen characters
	Slug *string `json:"slug" validationmap:"Slug" validate:"required" minLength:"3"`
} // @name	V1.Request.CreateOrganisation

func (req PostOrganisationRequest) ToCommand() organisations.CreateCommand {
	return organisations.CreateCommand{
		Name: req.Name,
		Slug: req.Slug,
	}
}

type DeleteOrganisationRequest struct {
	ID string `param:"id"`
}

func (req DeleteOrganisationRequest) ToCommand() organisations.DeleteCommand {
	return organisations.DeleteCommand{
		ID: req.ID,
	}
}

type GetOrganisationRequest struct {
	ID string `param:"id"`
}

func (req GetOrganisationRequest) ToCommand() organisations.GetCommand {
	return organisations.GetCommand{
		ID: req.ID,
	}
}

type PatchOrganisationRequest struct {
	ID string `param:"id" swaggerignore:"true"`

	// The name of the organisation.
	Name *string `json:"name,omitempty" validationmap:"Name" validate:"optional" minLength:"3" extensions:"x-nullable"`

	// A slug that will be used in URI's.
	// Must only contain alphanumeric and hyphen characters.
	// Must be unique across all other organisations.
	Slug *string `json:"slug,omitempty" validationmap:"Slug" validate:"optional" minLength:"3" extensions:"x-nullable"`
} // @name	V1.Request.UpdateOrganisation

func (req PatchOrganisationRequest) ToCommand() organisations.UpdateCommand {
	return organisations.UpdateCommand{
		ID:   req.ID,
		Name: req.Name,
		Slug: req.Slug,
	}
}
