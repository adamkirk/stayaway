package responses

import "github.com/adamkirk-stayaway/organisations/internal/domain/organisations"

type Organisation struct {
	// ID of the organisation.
	ID string `json:"id"`

	// The name of the organisation.
	Name string `json:"name"`

	// Unique slug used in URI's.
	Slug string `json:"slug"`
} // @name	V1.Response[Model].Organisation

func OrganisationFromModel(org *organisations.Organisation) Organisation {
	return Organisation{
		ID:   org.ID,
		Name: org.Name,
		Slug: org.Slug,
	}
}

type Organisations []Organisation // @name	V1.Response[Model].Organisations

func OrganisationsFromModels(orgs organisations.Organisations) Organisations {
	v1Orgs := make(Organisations, len(orgs))

	for i, org := range orgs {
		v1Orgs[i] = OrganisationFromModel(org)
	}

	return v1Orgs
}

type ListOrganisationsResponseBody struct {
	Meta ListResponseMeta `json:"meta"`
	Data Organisations `json:"data"`
}

type ListOrganisationsResponse struct {
	Body ListOrganisationsResponseBody
} // @name	V1.Response.ListOrganisations

type PostOrganisationResponse struct {
	Data Organisation `json:"data"`
} // @name	V1.Response.CreateOrganisation

type GetOrganisationResponse struct {
	Data Organisation `json:"data"`
} // @name	V1.Response.GetOrganisation

type PatchOrganisationResponse struct {
	Data Organisation `json:"data"`
} // @name	V1.Response.PatchOrganisation
