package api

import (
	"github.com/adamkirk-stayaway/organisations/internal/model"
	"github.com/adamkirk-stayaway/organisations/internal/municipalities"
)

type V1ListMunicipalitiesRequest struct {
	// The direction to order the results by.
	OrderDirection *string `query:"order_dir" json:"order_dir" validationmap:"OrderDirection" validate:"optional" enums:"asc,desc"`

	// The field by which to order the results.
	OrderBy *string `query:"order_by" json:"order_by" validationmap:"OrderBy" validate:"optional" enums:"name"`

	// The page to display.
	// An empty list may be returned if going beyond the last page of results.
	Page *int `query:"page" json:"page" validationmap:"Page" validate:"optional" minimum:"1"`

	// The amount of results to display per page.
	PerPage *int `query:"per_page" json:"per_page" validationmap:"PerPage" validate:"optional" minimum:"1" maximum:"100"`

	// Countries to filter the municipalities by.
	// Currently we only support United Kingdom anyway.
	Country []string `query:"country[]" json:"country" validationmap:"Country" validate:"optional"`

	// Characters to use as a prefix in searching by the name.
	// Useful for a "typeahead" widget.
	NamePrefix *string `query:"prefix" json:"prefix" validationmap:"NamePrefix" validate:"optional" minimum:"3"`
}

func (req V1ListMunicipalitiesRequest) ToCommand() municipalities.ListCommand {
	cmd := municipalities.NewListCommand()

	if req.OrderDirection != nil {
		cmd.OrderDirection = model.SortDirection(*req.OrderDirection)
	}

	if req.OrderBy != nil {
		cmd.OrderBy = model.MunicipalitySortBy(*req.OrderBy)
	}

	if req.Page != nil {
		cmd.Page = *req.Page
	}

	if req.PerPage != nil {
		cmd.PerPage = *req.PerPage
	}

	cmd.Country = req.Country
	cmd.NamePrefix = req.NamePrefix

	return cmd
}
