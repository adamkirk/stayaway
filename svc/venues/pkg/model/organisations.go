package model

type OrganisationOrderBy string

const (
	OrganisationOrderBySlug OrganisationOrderBy = "slug"
	OrganisationOrderByName OrganisationOrderBy = "name"
	OrganisationOrderByID OrganisationOrderBy = "id"
)

type Organisation struct {
	ID ID `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Organisations []*Organisation

