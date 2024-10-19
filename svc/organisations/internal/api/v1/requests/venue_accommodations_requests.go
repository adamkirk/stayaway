package requests

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/accommodations"
)

type PostVenueAccommodationOccupancy struct {
	// The minimum amount of people that must occupy the accommodation.
	// Or null if venue_template_id is supplied.
	Min *int `json:"min" validationmap:"MinOccupancy" validate:"required" minimum:"1" extensions:"x-nullable"`

	// The maximum amount of people that must occupy the accommodation.
	// If null or blank there will be no limit, or it will use the templates value.
	// If provided, it must be greater than the min occupancy.
	Max *int `json:"max" validationmap:"MaxOccupancy" validate:"optional" minimum:"1" extensions:"x-nullable"`
} // @name	V1.Request[Model].VenueAccommodationOccupancyCreate

type PostVenueAccommodationRequest struct {
	VenueID        string `param:"venueId" swaggerignore:"true"`

	VenueTemplateID        *string `json:"venue_template_id" validationmap:"VenueTemplateID"`

	// The name of the template.
	// Or null if venue_template_id is supplied.
	Name *string `json:"name" validationmap:"Name" validate:"required" extensions:"x-nullable" minLength:"3"`

	// The type of accommodation that this template is for.
	// Currently only supports 'room'.
	// Or null if venue_template_id is supplied.
	Type *string `json:"type" validationmap:"Type" validate:"required" extensions:"x-nullable" enums:"room"`

	// Description of the accommodation that this applies to.
	// Or null if venue_template_id is supplied.
	Description *string `json:"description" validationmap:"Description" validate:"required" minimum:"10"`

	Occupancy PostVenueAccommodationOccupancy `json:"occupancy" validate:"required"`
} // @name	V1.Request.CreateVenueAccommodation

func (req PostVenueAccommodationRequest) ToCommand() accommodations.CreateCommand {
	return accommodations.CreateCommand{
		VenueID: &req.VenueID,
		VenueTemplateID: req.VenueTemplateID,
		Name: req.Name,
		Description: req.Description,
		Type: req.Type,
		MinOccupancy: req.Occupancy.Min,
		MaxOccupancy: req.Occupancy.Max,
	}
}
