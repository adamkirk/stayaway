package model



type RoomFeature struct {
	Enabled bool
}

type EarlyCheckInFeature struct {
	RoomFeature
	Time string 
}

type LateCheckOutFeature struct {
	RoomFeature
	Time string

}

type RoomFeatures struct {
	EarlyCheckIn EarlyCheckInFeature 
	LateCheckOut LateCheckOutFeature 
}

type RoomSpec struct {
	ID string 
	VenueID string
	Name string
	Features RoomFeatures 
	MinimumOccupancy int
	MaximumOccupancy int
	CheckInTime string
	CheckOutTime string
}