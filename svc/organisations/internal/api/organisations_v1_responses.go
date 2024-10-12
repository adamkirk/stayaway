package api

import "github.com/adamkirk-stayaway/organisations/internal/domain/organisations"

type V1Organisation struct {
	// ID of the organisation.
	ID string `json:"id"`

	// The name of the organisation.
	Name string `json:"name"`

	// Unique slug used in URI's.
	Slug string `json:"slug"`
} // @name	V1.Response[Model].Organisation

func V1OrganisationFromModel(org *organisations.Organisation) V1Organisation {
	return V1Organisation{
		ID: org.ID,
		Name: org.Name,
		Slug: org.Slug,
	}
}

type V1Organisations []V1Organisation  // @name	V1.Response[Model].Organisations

func V1OrganisationsFromModels(orgs organisations.Organisations) V1Organisations {
	v1Orgs := make(V1Organisations, len(orgs))

	for i, org := range(orgs) {
		v1Orgs[i] = V1OrganisationFromModel(org)
	}

	return v1Orgs
}

type V1ListOrganisationsResponse struct {
	Data V1Organisations `json:"data"`
	Meta V1ListResponseMeta `json:"meta"`
} // @name	V1.Response.ListOrganisations

type V1PostOrganisationResponse struct {
	Data V1Organisation `json:"data"`
} // @name	V1.Response.CreateOrganisation

type V1GetOrganisationResponse struct {
	Data V1Organisation `json:"data"`
} // @name	V1.Response.GetOrganisation

type V1PatchOrganisationResponse struct {
	Data V1Organisation `json:"data"`
} // @name	V1.Response.PatchOrganisation