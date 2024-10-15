package organisations

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/mutex"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
)

type CreateCommand struct {
	Name *string `validate:"required,min=3"`
	Slug *string `validate:"required,min=3,slug"`
}

func (svc *Service) Create(cmd CreateCommand) (*Organisation, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	orgBySlug, err := svc.repo.BySlug(*cmd.Slug)

	if orgBySlug != nil {
		return nil, validation.ValidationError{
			Errs: []validation.FieldError{
				{
					Key:    "Slug",
					Errors: []string{"must be unique"},
				},
			},
		}
	}

	if err != nil {
		if _, ok := err.(common.ErrNotFound); !ok {
			return nil, err
		}
	}

	slugMutexKey := fmt.Sprintf("organisation_slug:%s", *cmd.Slug)
	l, err := svc.mutex.ClaimWithBackOff(slugMutexKey, 300*time.Millisecond)

	if err != nil {
		if _, ok := err.(mutex.ErrLockNotClaimed); ok {
			return nil, common.ErrConflict{
				Message: "slug is being used by another resource",
			}
		}

		return nil, err
	}

	defer func() {
		if err := l.Release(); err != nil {
			slog.Error("failed to release lock", "error", err, "key", slugMutexKey)
		}
	}()

	org := &Organisation{
		ID: svc.idGen.Generate(),
		Name: *cmd.Name,
		Slug: *cmd.Slug,
	}

	return svc.repo.Save(org)
}
