package responses

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/templates"
)

type VenueAccommodationTemplateOccupancy struct {
	// The minimum occupancy for any templates using this template.
	Min int `json:"min"`

	// The maximum occupancy for any templates using this template.
	// Null means there is no maximum.
	Max int `json:"max"`
} // @name	V1.Response[Model].VenueAccommodationTemplateOccupancy

type VenueAccommodationTemplate struct {
	// The ID of the template.
	ID string `json:"id"`

	// The ID of the venue this template exists in.
	VenueID string `json:"venue_id"`

	// The name of the template.
	Name string `json:"name"`

	// The type of accommodation this is.
	// Currently only 'room' is supported.
	Type string `json:"type"`

	// A description about the accommodation.
	Description string `json:"description"`

	// The settings for occupancy.
	Occupancy VenueAccommodationTemplateOccupancy `json:"occupancy"`
} // @name	V1.Response[Model].VenueAccommodationTemplate

func VenueAccommodationTemplateFromModel(v *templates.VenueTemplate) VenueAccommodationTemplate {

	return VenueAccommodationTemplate{
		ID:          v.ID,
		VenueID:     v.VenueID,
		Name:        v.Name,
		Type:        string(v.Type),
		Description: v.Description,
		Occupancy: VenueAccommodationTemplateOccupancy{
			Min: v.MinOccupancy,
			Max: v.MaxOccupancy,
		},
	}
}

type VenueAccommodationTemplates []VenueAccommodationTemplate // @name	V1.Response[Model].VenueAccommodationTemplates

func VenueAccommodationTemplatesFromModels(templates templates.VenueTemplates) VenueAccommodationTemplates {
	v1Orgs := make(VenueAccommodationTemplates, len(templates))

	for i, t := range templates {
		v1Orgs[i] = VenueAccommodationTemplateFromModel(t)
	}

	return v1Orgs
}

type PostVenueAccommodationTemplateResponse struct {
	Data VenueAccommodationTemplate `json:"data"`
} // @name	V1.Response.PostVenueAccommodationTemplate

type PatchVenueAccommodationTemplateResponse struct {
	Data VenueAccommodationTemplate `json:"data"`
} // @name	V1.Response.PatchVenueAccommodationTemplate

type GetVenueAccommodationTemplateResponse struct {
	Data VenueAccommodationTemplate `json:"data"`
} // @name	V1.Response.GetVenueAccommodationTemplate

type ListVenueAccommodationTemplatesResponse struct {
	Data VenueAccommodationTemplates `json:"data"`
	Meta ListResponseMeta            `json:"meta"`
} // @name	V1.Response.ListVenueAccommodationTemplates
