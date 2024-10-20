package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/templates"
	"github.com/adamkirk-stayaway/organisations/internal/util"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type ErrCannotUseOverridesForConfig struct {}

func (_ ErrCannotUseOverridesForConfig) Error() string {
	return "cannot use config from overrides when some fields are nil"
}

type SortBy string

const (
	SortByReference SortBy = "reference"
)

type Type string

type SearchFilter struct {
	VenueID []string
	ReferencePrefix *string
	VenueTemplateID []string
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
	Reference string `json:"reference" bson:"reference"`

	// This field shouldn't be saved, it will be generated when loading the model
	// It's set to nil in the repository save method
	Config *common.AccommodationConfig `bson:"-,omitempty"`
}

func (a *Accommodation) MergeTemplateConfig(other common.AccommodationConfig) {
	a.Config = &common.AccommodationConfig{
		Name: *util.Default[string](a.Overrides.Name, &other.Name),
		MinOccupancy: *util.Default[int](a.Overrides.MinOccupancy, &other.MinOccupancy),
		MaxOccupancy: *util.Default[int](a.Overrides.MaxOccupancy, &other.MaxOccupancy),
		Description:  *util.Default[string](a.Overrides.Description, &other.Description),
		Type:         *util.Default[common.AccommodationConfigType](a.Overrides.Type, &other.Type),
	}
}

type Accommodations []*Accommodation

func (all Accommodations) VenueTemplateIDs() []string {
	kv := map[string]bool{}

	for _, a := range all {
		if a.VenueTemplateID == nil {
			continue
		}
		kv[*a.VenueTemplateID] = true
		
	}
	
	ids := []string{}
	for k, _ := range kv {
		ids = append(ids, k)
	}

	return ids
}

type Validator interface {
	Validate(any) error
}

type AccommodationsRepo interface {
	ByReferenceAndVenueID(reference string, venueId string) (*Accommodation, error)
	Save(a *Accommodation) (*Accommodation, error)
	Delete(a *Accommodation) (error)
	Get(id string, venueId string) (*Accommodation, error)
	Paginate(p PaginationFilter, search SearchFilter) (Accommodations, common.PaginationResult, error)
}

type TemplatesRepo interface {
	Get(id string, venueId string) (*templates.VenueTemplate, error)
	ByID(ids []string, venueId string) (templates.VenueTemplates, error)
}

type Service struct {
	repo      AccommodationsRepo
	validator Validator
	validationMapper *validation.ValidationMapper
	idGen     common.IDGenerator
	templatesRepo TemplatesRepo

}

func (svc *Service) decorateWithTemplateConfig(t *templates.VenueTemplate, a *Accommodation) (error){
	if t != nil {
		a.MergeTemplateConfig(t.AccommodationConfig)

		return nil
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
	if a.Overrides.Name == nil {
		nilFields = append(nilFields, "Name")
	}
 
	if len(nilFields) > 0 {
		return ErrCannotUseOverridesForConfig{}
	}

	a.Config = &common.AccommodationConfig{
		Name: *a.Overrides.Name,
		MinOccupancy: *a.Overrides.MinOccupancy,
		MaxOccupancy: *a.Overrides.MaxOccupancy,
		Description:  *a.Overrides.Description,
		Type:         *a.Overrides.Type,
	}

	return nil
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
