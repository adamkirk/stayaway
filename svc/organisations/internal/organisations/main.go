package organisations

import (
	"fmt"
	"time"

	"github.com/adamkirk-stayaway/organisations/internal/mutex"
)

type Validator interface {
	Validate(any) error
}

type DistributedMutex interface {
	ClaimWithBackOff(key string, ttl time.Duration) (mutex.DistributedMutex, error)
	MultiClaimWithBackOff(keys []string, ttl time.Duration) (mutex.DistributedMutex, error)
}

func slugMutexKey(slug string) string {
	return fmt.Sprintf("organisation_slug:%s", slug)
}