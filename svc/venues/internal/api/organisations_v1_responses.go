package api

import "github.com/adamkirk-stayaway/venues/pkg/model"

type V1Organisation struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func V1OrganisationFromModel(org *model.Organisation) V1Organisation {
	return V1Organisation{
		ID: org.ID.String(),
		Name: org.Name,
		Slug: org.Slug,
	}
}

type V1Organisations []V1Organisation 

func V1OrganisationsFromModels(orgs model.Organisations) V1Organisations {
	v1Orgs := make(V1Organisations, len(orgs))

	for i, org := range(orgs) {
		v1Orgs[i] = V1OrganisationFromModel(org)
	}

	return v1Orgs
}

type V1ListOrganisationsMeta struct {
	OrderDirection string `json:"order_dir"`
	OrderBy string `json:"order_by"`
	Page int `json:"page"`
	PerPage int `json:"per_page"`
	TotalPages int 
}

type V1ListOrganisationsResponse struct {
	Data V1Organisations `json:"data"`
	Meta V1ListOrganisationsMeta `json:"meta"`
}