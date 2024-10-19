package common

type AccommodationConfigType string

const (
	AccommodationConfigTypeRoom AccommodationConfigType = "room"
)

func AllAccommodationConfigTypes() []string {
	return []string{
		string(AccommodationConfigTypeRoom),
	}
}

func (vt AccommodationConfigType) IsValid() bool {
	val := string(vt)

	for _, test := range AllAccommodationConfigTypes() {
		if test == val {
			return true
		}
	}

	return false
}

// AccommodationConfig is a generic type resource because we might assign a template to
// multiple things.
type AccommodationConfig struct {
	MinOccupancy int    `bson:"min_occupancy" validate:"required,min=1"`
	MaxOccupancy *int   `bson:"max_occupancy" validate:"omitnil,gtefield=MinOccupancy"`
	Description  string `bson:"description" validate:"omitnil,min=10"`
	Type         AccommodationConfigType   `bson:"type" validate:"required,accommodationtype"`
}

// IsValid checks that all the fields in the config are compatible.
// Will probably grow quite a lot, but for now, quite simple.
func (ac *AccommodationConfig) IsValid() bool {
	if ac.MaxOccupancy != nil && *ac.MaxOccupancy < ac.MinOccupancy {
		return false
	}

	return true
}

// AccommodationConfigOverrides is exactly the same as AccommodationConfig except
// all values are pointers as they are optional.
type AccommodationConfigOverrides struct {
	MinOccupancy *int    `bson:"min_occupancy"`
	MaxOccupancy *int   `bson:"max_occupancy"`
	Description  *string `bson:"description"`
	Type         *AccommodationConfigType   `bson:"type"`
}