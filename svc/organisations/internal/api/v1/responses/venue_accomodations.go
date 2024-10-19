package responses

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/accommodations"
	"github.com/adamkirk-stayaway/organisations/internal/util"
)

type VenueAccommodationOccupancy struct {
	// The minimum occupancy for any templates using this template.
	Min int `json:"min"`

	// The maximum occupancy for any templates using this template.
	// Null means there is no maximum.
	Max *int `json:"max"`
} // @name	V1.Response[Model].VenueAccommodationOccupancy

type VenueAccommodationOccupancyOverrides struct {
	// The minimum occupancy for any templates using this template.
	// Can be null if the accommodation has a template.
	Min *int `json:"min"`

	// The maximum occupancy for any templates using this template.
	// Null means there is no maximum.
	Max *int `json:"max"`
} // @name	V1.Response[Model].VenueAccommodationOccupancyOverrides

// The final configuration used by the accommodation. This is the result of 
// merging any configuration from the template and the overrides of this 
// accommodation. If there was no template, this is identical to the overrides.
type VenueAccommodationConfig struct {
	Occupancy VenueAccommodationOccupancy `json:"occupancy"`
	Description  string `json:"description"`
	Type         string   `json:"type"`
} // @name	V1.Response[Model].VenueAccommodationConfig

// The overrides that were supplied upon creation. These may all be set, or may
// all be null. They may only be null if the accommodation uses a template.
type VenueAccommodationOverrides struct {
	Occupancy VenueAccommodationOccupancyOverrides `json:"occupancy"`
	Description  *string `json:"description"`
	Type         *string   `json:"type"`
} // @name	V1.Response[Model].VenueAccommodationOverrides

type VenueAccommodation struct {
	// The ID of the template.
	ID string `json:"id"`

	// The ID of the venue this template exists in.
	VenueID string `json:"venue_id"`
	
	// The ID of the venue template used by this accommodation.
	VenueTemplateID *string `json:"venue_template_id"`

	// The overrides that were supplied upon creation. These may all be set, or may
	// all be null. They may only be null if the accommodation uses a template.
	Overrides VenueAccommodationOverrides `json:"overrides"`

	Name string `json:"name"`


	// The final configuration used by the accommodation. This is the result of 
	// merging any configuration from the template and the overrides of this 
	// accommodation. If there was no template, this is identical to the overrides.
	Config VenueAccommodationConfig `json:"config"`
} // @name	V1.Response[Model].VenueAccommodation

func VenueAccommodationFromModel(a *accommodations.Accommodation) VenueAccommodation {
	var overriddenType *string

	if a.Overrides.Type != nil {
		overriddenType = util.PointTo[string](string(*a.Overrides.Type))
	}

	return VenueAccommodation{
		ID:          a.ID,
		VenueID:     a.VenueID,
		VenueTemplateID:     a.VenueTemplateID,
		Name: a.Name,
		Config: VenueAccommodationConfig{
			Description: a.Config.Description,
			Type: string(a.Config.Type),
			Occupancy: VenueAccommodationOccupancy{
				Min: a.Config.MinOccupancy,
				Max: a.Config.MaxOccupancy,
			},
		},
		Overrides: VenueAccommodationOverrides{
			Description: a.Overrides.Description,
			Type: overriddenType,
			Occupancy: VenueAccommodationOccupancyOverrides{
				Min: a.Overrides.MinOccupancy,
				Max: a.Overrides.MaxOccupancy,
			},
		},
	}
}

type VenueAccommodations []VenueAccommodation // @name	V1.Response[Model].VenueAccommodations

func VenueAccommodationsFromModels(accs accommodations.Accommodations) VenueAccommodations {
	v1Accs := make(VenueAccommodations, len(accs))

	for i, a := range accs {
		v1Accs[i] = VenueAccommodationFromModel(a)
	}

	return v1Accs
}

type PostVenueAccommodationResponse struct {
	Data VenueAccommodation `json:"data"`
} // @name	V1.Response.PostVenueAccommodation

type PatchVenueAccommodationResponse struct {
	Data VenueAccommodation `json:"data"`
} // @name	V1.Response.PatchVenueAccommodation

type GetVenueAccommodationResponse struct {
	Data VenueAccommodation `json:"data"`
} // @name	V1.Response.GetVenueAccommodation

type ListVenueAccommodationsResponse struct {
	Data VenueAccommodations `json:"data"`
	Meta ListResponseMeta            `json:"meta"`
} // @name	V1.Response.ListVenueAccommodations
