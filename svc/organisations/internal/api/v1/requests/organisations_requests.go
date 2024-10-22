package requests

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/organisations"
)

type ListOrganisationsRequest struct {
	OrderDirection string `query:"order_dir" json:"order_dir" validationmap:"OrderDir" validate:"optional" enums:"asc,desc" doc:"The direction to order the results by."`
	OrderBy string `query:"order_by" json:"order_by" validationmap:"OrderBy" validate:"optional" enums:"name,slug" doc:"The field by which to order the results."`
	Page int `query:"page" json:"page" validationmap:"Page" validate:"optional" minimum:"1" doc:"The page to display.<br />An empty list may be returned if going beyond the last page of results."`
	PerPage int `query:"per_page" json:"per_page" validationmap:"PerPage" validate:"optional" minimum:"1" maximum:"100" doc:"The amount of results to display per page."`
}

func (req *ListOrganisationsRequest) ToCommand() organisations.ListCommand {
	cmd := organisations.NewListCommand()

	if req.OrderDirection != "" {
		cmd.OrderDirection = common.SortDirection(req.OrderDirection)
	}

	if req.OrderBy != "" {
		cmd.OrderBy = organisations.SortBy(req.OrderBy)
	}


	if req.Page != 0 {
		cmd.Page = req.Page
	}

	if req.PerPage != 0 {
		cmd.PerPage = req.PerPage
	}

	return cmd
}

type PostOrganisationRequest struct {
	Body PostOrganisationRequestBody
}

type PostOrganisationRequestBody struct {
	Name *string `json:"name,omitempty" required:"true" minLength:"3" example:"My Organisation" validationmap:"Name" doc:"The name of the organisation."`

	Slug *string `json:"slug,omitempty" required:"true" minLength:"3" pattern:"^[a-z0-9]{1}[a-z0-9\\-]*$" patternDescription:"alphanum + hyphen" validationmap:"Slug" example:"my-organisation" doc:"A slug that will be used in URI's. Must be unique across all other organisations."`
}

func (req *PostOrganisationRequest) ToCommand() organisations.CreateCommand {
	return organisations.CreateCommand{
		Name: req.Body.Name,
		Slug: req.Body.Slug,
	}
}

type DeleteOrganisationRequest struct {
	ID string `path:"id"`
}

func (req DeleteOrganisationRequest) ToCommand() organisations.DeleteCommand {
	return organisations.DeleteCommand{
		ID: req.ID,
	}
}

type GetOrganisationRequest struct {
	ID string `path:"id"`
}

func (req GetOrganisationRequest) ToCommand() organisations.GetCommand {
	return organisations.GetCommand{
		ID: req.ID,
	}
}

type PatchOrganisationRequest struct {
	ID string `path:"id"`
	Body PatchOrganisationRequestBody
}

type PatchOrganisationRequestBody struct {
	Name *string `json:"name,omitempty" required:"false" example:"My Organisation" validationmap:"Name" validate:"optional" minLength:"3" doc:"The name of the organisation."`
	Slug *string `json:"slug,omitempty" required:"false" example:"my-organisation" validationmap:"Slug" pattern:"^[a-z0-9]{1}[a-z0-9\\-]*$" patternDescription:"alphanum + hyphen" minLength:"3" extensions:"x-nullable" doc:"A slug that will be used in URI's. Must be unique across all other organisations."`
}

func (req PatchOrganisationRequest) ToCommand() organisations.UpdateCommand {
	return organisations.UpdateCommand{
		ID:   req.ID,
		Name: req.Body.Name,
		Slug: req.Body.Slug,
	}
}
