package rooms

// Template is a generic type resource because we might assign a template to 
// multiple things.
type Template struct {
	ID string
	Name string
	MinOccupancy int
	MaxOccupancy int
}

type VenueTemplate struct {
	Template
	VenueID string
}
