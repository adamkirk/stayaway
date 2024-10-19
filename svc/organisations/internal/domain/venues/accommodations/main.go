package accommodations

import (
	"errors"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/templates"
	"github.com/adamkirk-stayaway/organisations/internal/util"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type SortBy string

const (
	SortByName SortBy = "name"
)

type Type string

type SearchFilter struct {
	VenueID []string
}

type PaginationFilter struct {
	OrderBy  SortBy
	OrderDir common.SortDirection
	Page     int
	PerPage  int
}

type Accommodation struct {
	ID              string                               `json:"id" bson:"_id,omitempty"`
	VenueID         string                               `json:"venue_id" bson:"venue_id"`
	VenueTemplateID *string                              `json:"venue_template_id" bson:"venue_template_id"`
	Overrides       *common.AccommodationConfigOverrides `bson:"overrides"`
	Name string `json:"name" bson:"name"`

	// This field shouldn't be saved, it will be generated when loading the model
	// It's set to nil in the repository save method
	Config *common.AccommodationConfig `bson:"-,omitempty"`
}

func (a *Accommodation) MergeTemplateConfig(other common.AccommodationConfig) {
	a.Config = &common.AccommodationConfig{
		MinOccupancy: *util.Default[int](a.Overrides.MinOccupancy, &other.MinOccupancy),
		MaxOccupancy: util.Default[int](a.Overrides.MaxOccupancy, other.MaxOccupancy),
		Description:  *util.Default[string](a.Overrides.Description, &other.Description),
		Type:         *util.Default[common.AccommodationConfigType](a.Overrides.Type, &other.Type),
	}
}

type Accommodations []*Accommodation

type Validator interface {
	Validate(any) error
}

type AccommodationsRepo interface {
	ByNameAndVenueID(name string, venueId string) (*Accommodation, error)
	Save(a *Accommodation) (*Accommodation, error)
}

type TemplatesRepo interface {
	Get(id string, venueId string) (*templates.VenueTemplate, error)
}

type Service struct {
	repo      AccommodationsRepo
	validator Validator
	validationMapper *validation.ValidationMapper
	idGen     common.IDGenerator
	templatesRepo TemplatesRepo

}

func (svc *Service) decorateWithTemplateConfig(t *templates.VenueTemplate, a *Accommodation) (*Accommodation, error){
	if t != nil {
		a.MergeTemplateConfig(t.AccommodationConfig)

		return a, nil
	}

	nilFields := []string{}

	// This is tedious but only other way i can think to do it is to use reflection 
	// and loop through fields, which just sounds a bit bleh. This'll do until 
	// the amount of fields gets out of hand.
	if a.Overrides.MinOccupancy == nil {
		nilFields = append(nilFields, "MinOccupancy")
	}
	if a.Overrides.Description == nil {
		nilFields = append(nilFields, "Description")
	}
	if a.Overrides.Type == nil {
		nilFields = append(nilFields, "Type")
	}

	// In theory this should never return an error; if it does it means the 
	// accommodation is in an invalid state. If the template id was nil (which 
	// it has to be to get here) then all of the fields should've been populated. 
	if len(nilFields) > 0 {
		return nil, errors.New("cannot use config from overrides when some fields are nil")
	}

	a.Config = &common.AccommodationConfig{
		MinOccupancy: *a.Overrides.MinOccupancy,
		MaxOccupancy: a.Overrides.MaxOccupancy,
		Description:  *a.Overrides.Description,
		Type:         *a.Overrides.Type,
	}

	return a, nil
}

func NewService(repo AccommodationsRepo, templatesRepo TemplatesRepo, validationMapper *validation.ValidationMapper, v Validator, idGen common.IDGenerator) *Service {
	return &Service{
		repo:      repo,
		validator: v,
		idGen:     idGen,
		validationMapper: validationMapper,
		templatesRepo: templatesRepo,
	}
}
