package model

import (
	"github.com/google/uuid"
)

// ID represents a unique ID of something in the system
// TODO: Custom unmarshalling to/from JSON
type ID struct {
	uuidV7 uuid.UUID
}

func (id *ID) Equals(other *ID) bool {
	return id.uuidV7.String() == other.uuidV7.String()
}

func (id *ID) String() string {
	return id.uuidV7.String()
}

func NewID() *ID {
	return &ID{
		uuidV7: uuid.New(),
	}
}

