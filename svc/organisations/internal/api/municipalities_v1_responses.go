package api

import "github.com/adamkirk-stayaway/organisations/internal/model"

type V1Municipality struct {
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

func V1MunicipalityFromModel(v model.Municipality) V1Municipality {
	return V1Municipality{
		ID: v.ID,
		Name: v.Name,
		NameAscii: v.NameAscii,
		Lat: v.Lat,
		Long: v.Long,
		Country: v.Country,
		Iso3: v.Iso3,
	}
}

type V1Municipalities []V1Municipality // @name	V1.Response[Model].Municipalities

func V1MunicipalitiesFromModels(venues model.Municipalities) V1Municipalities {
	v1Orgs := make(V1Municipalities, len(venues))

	for i, v := range(venues) {
		v1Orgs[i] = V1MunicipalityFromModel(v)
	}

	return v1Orgs
}

type V1ListMunicipalitiesResponse struct {
	Data V1Municipalities `json:"data"`
	Meta V1ListResponseMeta `json:"meta"`
} // @name	V1.Response.ListMunicipalities
