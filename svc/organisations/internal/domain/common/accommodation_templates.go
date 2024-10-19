package common

type AccommodationTemplateType string

const (
	AccommodationTemplateTypeRoom AccommodationTemplateType = "room"
)

func AllAccommodationTemplateTypes() []string {
	return []string{
		string(AccommodationTemplateTypeRoom),
	}
}

func (vt AccommodationTemplateType) IsValid() bool {
	val := string(vt)

	for _, test := range AllAccommodationTemplateTypes() {
		if test == val {
			return true
		}
	}

	return false
}

// AccommodationTemplate is a generic type resource because we might assign a template to
// multiple things.
type AccommodationTemplate struct {
	Name         string `bson:"name"`
	MinOccupancy int    `bson:"min_occupancy"`
	MaxOccupancy *int   `bson:"max_occupancy"`
	Description  string `bson:"description"`
	Type         AccommodationTemplateType   `bson:"type"`
}