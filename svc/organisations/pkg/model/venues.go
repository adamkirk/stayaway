package model

type VenueSortBy string

const (
	VenueSortBySlug VenueSortBy = "slug"
	VenueSortByName VenueSortBy = "name"
)

type VenueType string

type VenueSearchFilter struct {
	OrganisationID []string
}

type VenuePaginationFilter struct {
	OrderBy VenueSortBy
	OrderDir SortDirection
	Page int
	PerPage int
}

const (
	VenueTypeHotel = "hotel"
)

func AllVenueTypes() []string {
	return []string{
		VenueTypeHotel,
	}
}

func (vt VenueType) IsValid() bool {
	val := string(vt)

	for _, test := range AllVenueTypes() {
		if test == val {
			return true
		}
	}

	return false
}

type VenueCoordinates struct {
	Lat float64 `json:"lat" bson:"lat"`
	Long float64 `json:"long" bson:"long"`
}

type VenueAddress struct {
	Line1 string `json:"line_1" bson:"line_1"`
	Line2 *string `json:"line_2" bson:"line_2"`
	Municipality string `json:"municipality" bson:"municipality"`
	PostCode string `json:"postcode" bson:"postcode"`
	Coordinates *VenueCoordinates `json:"coordinates" bson:"coordinates"`
}

type Venue struct {
	ID string `json:"id" bson:"_id,omitempty"`
	OrganisationID string `json:"organisation_id" bson:"organisation_id"`
	Name string `json:"name" bson:"name"`
	Slug string `json:"slug" bson:"slug"`
	Type VenueType `json:"type" bson:"type"`
	Address *VenueAddress `json:"address" bson:"address"`
}

type Venues []*Venue