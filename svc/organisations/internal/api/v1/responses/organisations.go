package responses

import "github.com/adamkirk-stayaway/organisations/internal/domain/organisations"

type Organisation struct {
	ID string `json:"id" doc:"ID of the organisation."`
	Name string `json:"name" doc:"The name of the organisation."`
	Slug string `json:"slug" doc:"Unique slug used in URI's."`
}

func OrganisationFromModel(org *organisations.Organisation) Organisation {
	return Organisation{
		ID:   org.ID,
		Name: org.Name,
		Slug: org.Slug,
	}
}

type Organisations []Organisation

func OrganisationsFromModels(orgs organisations.Organisations) Organisations {
	v1Orgs := make(Organisations, len(orgs))

	for i, org := range orgs {
		v1Orgs[i] = OrganisationFromModel(org)
	}

	return v1Orgs
}

type OrganisationsListResponseBody struct {
	Meta ListResponseMeta `json:"meta"`
	Data Organisations `json:"data"`
}

type ListOrganisationsResponse struct {
	Body OrganisationsListResponseBody
}

type OrganisationResponseBody struct {
	Data Organisation `json:"data"`
}

type OrganisationResponse struct {
	Body OrganisationResponseBody
}