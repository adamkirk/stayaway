package model


type MunicipalitySortBy string

const (
	MunicipalitySortByName MunicipalitySortBy = "city"
)

type MunicipalityType string

type MunicipalitySearchFilter struct {
	Country []string
	NamePrefix *string
}

type MunicipalityPaginationFilter struct {
	OrderBy MunicipalitySortBy
	OrderDir SortDirection
	Page int
	PerPage int
}

type Municipality struct {
	ID string `bson:"_id,omitempty"`
	Name string `csv:"city" bson:"name"`
	NameAscii string `csv:"city_ascii" bson:"name_ascii"`
	Lat float64 `csv:"lat"`
	Long float64 `csv:"lng"`
	Country string `csv:"country"`
	Iso3 string `csv:"iso3"`
	ImportID int `csv:"id" bson:"import_id"`
}

type Municipalities []Municipality

type BatchUpdateResult struct {
	Updated int
	Created int
}

type SyncResult struct {
	Processed int
	Path string
}