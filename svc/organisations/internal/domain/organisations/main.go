package organisations

import (
	"fmt"
	"time"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/mutex"
)

type DistributedMutex interface {
	ClaimWithBackOff(key string, ttl time.Duration) (mutex.DistributedMutex, error)
	MultiClaimWithBackOff(keys []string, ttl time.Duration) (mutex.DistributedMutex, error)
}

func slugMutexKey(slug string) string {
	return fmt.Sprintf("organisation_slug:%s", slug)
}

type SortBy string

const (
	SortBySlug SortBy = "slug"
	SortByName SortBy = "name"
)

type Organisation struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Slug string `json:"slug" bson:"slug"`
}

type Organisations []*Organisation


type Validator interface {
	Validate(any) error
}

type OrganisationsRepo interface {
	Save(org *Organisation) (*Organisation, error)
	BySlug(slug string) (*Organisation, error)
	Get(id string) (*Organisation, error)
	Delete(*Organisation) error
	Paginate(orderBy SortBy, orderDir common.SortDirection, page int, perPage int) (Organisations, common.PaginationResult, error)
}

type Service struct {
	repo OrganisationsRepo
	validator Validator
	mutex DistributedMutex
}

func NewService(repo OrganisationsRepo, v Validator, mutex DistributedMutex) *Service {
	return &Service{
		repo: repo,
		validator: v,
		mutex: mutex,
	}
}