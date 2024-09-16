package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/organisations"
)


type V1ListOrganisationsRequest struct {
	OrderDirection *string `query:"order_dir"`
	OrderBy *string `query:"order_by"`
	Page *int `query:"page"`
	PerPage *int `query:"per_page"`
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
	Name string `json:"name"`
	Slug string `json:"slug"`
}

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
	ID string `param:"id"`
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

func (req V1PatchOrganisationRequest) ToCommand() organisations.UpdateCommand {
	return organisations.UpdateCommand{
		ID: req.ID,
		Name: req.Name,
		Slug: req.Slug,
	}
}