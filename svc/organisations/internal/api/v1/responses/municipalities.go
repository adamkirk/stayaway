package responses

import "github.com/adamkirk-stayaway/organisations/internal/domain/municipalities"

type Municipality struct {
	// The ID of the municipality.
	ID string `json:"id"`

	// The "common name" of the town/city/village.
	Name string `json:"name"`

	// An ascii friendly representation of the name.
	NameAscii string `json:"name_ascii"`

	// The latitude of the municipality.
	Lat float64 `json:"lat"`

	// The longitude of the municipality.
	Long float64 `json:"long"`

	// The country in which the municipality resides.
	Country string `json:"country"`

	// ISO code for the country.
	Iso3 string `json:"iso3"`
} // @name	V1.Response[Model].Municipality

func MunicipalityFromModel(v municipalities.Municipality) Municipality {
	return Municipality{
		ID: v.ID,
		Name: v.Name,
		NameAscii: v.NameAscii,
		Lat: v.Lat,
		Long: v.Long,
		Country: v.Country,
		Iso3: v.Iso3,
	}
}

type Municipalities []Municipality // @name	V1.Response[Model].Municipalities

func MunicipalitiesFromModels(venues municipalities.Municipalities) Municipalities {
	v1Orgs := make(Municipalities, len(venues))

	for i, v := range(venues) {
		v1Orgs[i] = MunicipalityFromModel(v)
	}

	return v1Orgs
}

type ListMunicipalitiesResponse struct {
	Data Municipalities `json:"data"`
	Meta ListResponseMeta `json:"meta"`
} // @name	V1.Response.ListMunicipalities
