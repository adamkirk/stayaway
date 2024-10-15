package venues

import "github.com/adamkirk-stayaway/organisations/internal/domain/common"

type SortBy string

const (
	SortBySlug SortBy = "slug"
	SortByName SortBy = "name"
)

type Type string

type SearchFilter struct {
	OrganisationID []string
}

type PaginationFilter struct {
	OrderBy  SortBy
	OrderDir common.SortDirection
	Page     int
	PerPage  int
}

const (
	TypeHotel = "hotel"
)

func AllTypes() []string {
	return []string{
		TypeHotel,
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

type Coordinates struct {
	Lat  float64 `json:"lat" bson:"lat"`
	Long float64 `json:"long" bson:"long"`
}

type Address struct {
	Line1        string       `json:"line_1" bson:"line_1"`
	Line2        *string      `json:"line_2" bson:"line_2"`
	Municipality string       `json:"municipality" bson:"municipality"`
	PostCode     string       `json:"postcode" bson:"postcode"`
	Coordinates  *Coordinates `json:"coordinates" bson:"coordinates"`
}

type Venue struct {
	ID             string   `json:"id" bson:"_id,omitempty"`
	OrganisationID string   `json:"organisation_id" bson:"organisation_id"`
	Name           string   `json:"name" bson:"name"`
	Slug           string   `json:"slug" bson:"slug"`
	Type           Type     `json:"type" bson:"type"`
	Address        *Address `json:"address" bson:"address"`
}

type Venues []*Venue


type Validator interface {
	Validate(any) error
}

type VenuesRepo interface {
	Save(org *Venue) (*Venue, error)
	BySlugAndOrganisation(slug string, orgId string) (*Venue, error)
	Get(id string, orgId string) (*Venue, error)
	Delete(v *Venue) error
	Paginate(p PaginationFilter, search SearchFilter) (Venues, common.PaginationResult, error)
}

type Service struct {
	repo VenuesRepo
	validator Validator
}

func NewService(repo VenuesRepo, v Validator) *Service {
	return &Service{
		repo: repo,
		validator: v,
	}
}