package api

import "github.com/adamkirk-stayaway/organisations/pkg/model"

type V1VenueAddress struct {
	// Line 1 of the address. Typically number/name and street/road.
	Line1 string `json:"line_1"`

	// Line 2 of the address. May be null.
	Line2 *string `json:"line_2"`

	// The town/city in which the venue is situated.
	Municipality string `json:"municipality"`
	
	// The postcode of the venue.
	PostCode string `json:"postcode"`

	// The latitude of the venue.
	Lat float64 `json:"lat"`

	// The longiitude of the venue.
	Long float64 `json:"long"`
} // @name	V1.Response[Model].VenueAddress

type V1Venue struct {
	// The ID of the venue.
	ID string `json:"id"`

	// The ID of the organisation that this venue is a part of.
	OrganisationID string `json:"organisation_id"`

	// The name of the venue.
	Name string `json:"name"`

	// Slug used in URI's. Unique across all venues in the organisation.
	Slug string `json:"slug"`

	// The type of venue this is.
	// Currently only 'hotel' is supported.
	Type string `json:"type"`

	// The address of the venue.
	Address V1VenueAddress `json:"address"`
} // @name	V1.Response[Model].Venue

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

type V1Venues []V1Venue  // @name	V1.Response[Model].Venues

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
} // @name	V1.Response.ListVenues

type V1PostVenueResponse struct {
	Data V1Venue `json:"data"`
} // @name	V1.Response.PostVenue

type V1GetVenueResponse struct {
	Data V1Venue `json:"data"`
} // @name	V1.Response.GetVenue

type V1PatchVenueResponse struct {
	Data V1Venue `json:"data"`
} // @name	V1.Response.PatchVenue