package model

type OrganisationSortBy string

const (
	OrganisationSortBySlug OrganisationSortBy = "slug"
	OrganisationSortByName OrganisationSortBy = "name"
)

type Organisation struct {
	ID string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Slug string `json:"slug" bson:"slug"`
}

type Organisations []*Organisation

