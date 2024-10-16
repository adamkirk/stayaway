package municipalities

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/spf13/afero"
)

type SortBy string

const (
	SortByName SortBy = "name"
)

type Type string

type SearchFilter struct {
	Country    []string
	NamePrefix *string
}

type PaginationFilter struct {
	OrderBy  SortBy
	OrderDir common.SortDirection
	Page     int
	PerPage  int
}

type Municipality struct {
	ID        string  `bson:"_id,omitempty"`
	Name      string  `csv:"city" bson:"name"`
	NameAscii string  `csv:"city_ascii" bson:"name_ascii"`
	Lat       float64 `csv:"lat"`
	Long      float64 `csv:"lng"`
	Country   string  `csv:"country"`
	Iso3      string  `csv:"iso3"`
	ImportID  int     `csv:"id" bson:"import_id"`
}

type Municipalities []Municipality

type BatchUpdateResult struct {
	Updated int
	Created int
}

type SyncResult struct {
	Processed int
	Path      string
}

type MunicipalitiesRepo interface {
	Paginate(p PaginationFilter, search SearchFilter) (Municipalities, common.PaginationResult, error)
	UpdateBatch(batch []Municipality) (BatchUpdateResult, error)
}

type Config interface {
	MunicipalitiesSyncBatchSize() int
	MunicipalitiesSyncMaxProcesses() int
	MunicipalitiesSyncCountries() []string
}

type Validator interface {
	Validate(any) error
}

type Service struct {
	repo      MunicipalitiesRepo
	validator Validator
	cfg       Config
	fs        afero.Fs
}

func NewService(repo MunicipalitiesRepo, v Validator, cfg Config, fs afero.Fs) *Service {
	return &Service{
		repo:      repo,
		validator: v,
		cfg:       cfg,
		fs:        fs,
	}
}
