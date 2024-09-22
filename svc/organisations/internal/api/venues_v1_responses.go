package api

import "github.com/adamkirk-stayaway/organisations/pkg/model"

type V1VenueAddress struct {
	Line1 string `json:"line_1"`
	Line2 *string `json:"line_2"`
	Municipality string `json:"municipality"`
	PostCode string `json:"postcode"`
	Lat float64 `json:"lat"`
	Long float64 `json:"long"`
}

type V1Venue struct {
	ID string `json:"id"`
	OrganisationID string `json:"organisation_id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Type string `json:"type"`
	Address V1VenueAddress `json:"address"`
}

func V1VenueFromModel(v *model.Venue) V1Venue {
	return V1Venue{
		ID: v.ID,
		Name: v.Name,
		Slug: v.Slug,
		Type: string(v.Type),
		OrganisationID: v.OrganisationID,
		Address: V1VenueAddress{
			Line1: v.Address.Line1,
			Line2: v.Address.Line2,
			Municipality: v.Address.Municipality,
			PostCode: v.Address.PostCode,
			Lat: v.Address.Coordinates.Lat,
			Long: v.Address.Coordinates.Long,
		},
	}
}

type V1Venues []V1Venue

func V1VenuesFromModels(venues model.Venues) V1Venues {
	v1Orgs := make(V1Venues, len(venues))

	for i, v := range(venues) {
		v1Orgs[i] = V1VenueFromModel(v)
	}

	return v1Orgs
}

type V1ListVenuesResponse struct {
	Data V1Venues `json:"data"`
	Meta V1ListResponseMeta `json:"meta"`
}

type V1PostVenueResponse struct {
	Data V1Venue `json:"data"`
}

type V1GetVenueResponse struct {
	Data V1Venue `json:"data"`
}

// type V1PatchOrganisationResponse struct {
// 	Data V1Organisation `json:"data"`
// }