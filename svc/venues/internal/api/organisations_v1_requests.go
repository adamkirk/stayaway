package api

import "github.com/adamkirk-stayaway/venues/pkg/model"


type V1ListOrganisationsRequest struct {
	OrderDirection model.SortDirection `query:"order_dir"`
	OrderBy model.OrganisationSortBy `query:"order_by"`
	Page int `query:"page"`
	PerPage int `query:"per_page"`
}

type V1PostOrganisationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}