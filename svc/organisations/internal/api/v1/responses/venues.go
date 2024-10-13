package responses

import "github.com/adamkirk-stayaway/organisations/internal/domain/venues"

type VenueAddress struct {
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

type Venue struct {
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
	Address VenueAddress `json:"address"`
} // @name	V1.Response[Model].Venue

func VenueFromModel(v *venues.Venue) Venue {
	return Venue{
		ID:             v.ID,
		Name:           v.Name,
		Slug:           v.Slug,
		Type:           string(v.Type),
		OrganisationID: v.OrganisationID,
		Address: VenueAddress{
			Line1:        v.Address.Line1,
			Line2:        v.Address.Line2,
			Municipality: v.Address.Municipality,
			PostCode:     v.Address.PostCode,
			Lat:          v.Address.Coordinates.Lat,
			Long:         v.Address.Coordinates.Long,
		},
	}
}

type Venues []Venue // @name	V1.Response[Model].Venues

func VenuesFromModels(venues venues.Venues) Venues {
	v1Orgs := make(Venues, len(venues))

	for i, v := range venues {
		v1Orgs[i] = VenueFromModel(v)
	}

	return v1Orgs
}

type ListVenuesResponse struct {
	Data Venues           `json:"data"`
	Meta ListResponseMeta `json:"meta"`
} // @name	V1.Response.ListVenues

type PostVenueResponse struct {
	Data Venue `json:"data"`
} // @name	V1.Response.PostVenue

type GetVenueResponse struct {
	Data Venue `json:"data"`
} // @name	V1.Response.GetVenue

type PatchVenueResponse struct {
	Data Venue `json:"data"`
} // @name	V1.Response.PatchVenue
