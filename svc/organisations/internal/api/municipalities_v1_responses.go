package api

import "github.com/adamkirk-stayaway/organisations/pkg/model"

type V1Municipality struct {
	ID string `json:"id"`
	Name string `json:"name"`
	NameAscii string `json:"name_ascii"`
	Lat float64 `json:"lat"`
	Long float64 `json:"long"`
	Country string `json:"country"`
	Iso3 string `json:"iso3"`
}

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

type V1Municipalities []V1Municipality

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
}
