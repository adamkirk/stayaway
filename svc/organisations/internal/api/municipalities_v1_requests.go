package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/municipalities"
)

type V1ListMunicipalitiesRequest struct {
	OrderDirection *string `query:"order_dir" json:"order_dir" validationmap:"OrderDirection"`
	OrderBy *string `query:"order_by" validationmap:"OrderBy"`
	Page *int `query:"page" validationmap:"Page"`
	PerPage *int `query:"per_page" json:"per_page" validationmap:"PerPage"`
	Country []string `query:"country[]" validationmap:"Country"`
	NamePrefix *string `query:"prefix" json:"prefix" validationmap:"NamePrefix"`
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
