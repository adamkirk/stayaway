package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
)

type SortBy string

const (
	SortByName SortBy = "name"
)

type SearchFilter struct {
	VenueID    []string
	NamePrefix *string
}

type PaginationFilter struct {
	OrderBy  SortBy
	OrderDir common.SortDirection
	Page     int
	PerPage  int
}

type Type string

const (
	TypeRoom Type = "room"
)

func AllTypes() []string {
	return []string{
		string(TypeRoom),
	}
}

func (vt Type) IsValid() bool {
	val := string(vt)

	for _, test := range AllTypes() {
		if test == val {
			return true
		}
	}

	return false
}

// Template is a generic type resource because we might assign a template to
// multiple things.
type Template struct {
	Name         string `bson:"name"`
	MinOccupancy int    `bson:"min_occupancy"`
	MaxOccupancy *int   `bson:"max_occupancy"`
	Description  string `bson:"description"`
	Type         Type   `bson:"type"`
}

type VenueTemplate struct {
	Template `bson:"inline,template"`
	ID       string `bson:"_id,omitempty"`
	VenueID  string `bson:"venue_id"`
}

type Templates []Template

type VenueTemplates []*VenueTemplate

type AccommodationGroup struct {
	TemplateID  string
	Name        string
	Slots       int
	Description string
}

type AccommodationGroups []AccommodationGroup

type Validator interface {
	Validate(any) error
}

type VenueTemplatesRepo interface {
	Save(org *VenueTemplate) (*VenueTemplate, error)
	ByNameAndVenue(name string, venueId string) (*VenueTemplate, error)
	Get(id string, venueId string) (*VenueTemplate, error)
	Delete(*VenueTemplate) error
	Paginate(p PaginationFilter, search SearchFilter) (VenueTemplates, common.PaginationResult, error)
}

type VenuesRepo interface {
	Get(id string, orgId string) (*venues.Venue, error)
}

type VenueTemplatesService struct {
	repo       VenueTemplatesRepo
	venuesRepo VenuesRepo
	validator  Validator
	idGen      common.IDGenerator
}

func NewVenueTemplatesService(repo VenueTemplatesRepo, venuesRepo VenuesRepo, v Validator, idGen common.IDGenerator) *VenueTemplatesService {
	return &VenueTemplatesService{
		repo:       repo,
		venuesRepo: venuesRepo,
		validator:  v,
		idGen:      idGen,
	}
}
