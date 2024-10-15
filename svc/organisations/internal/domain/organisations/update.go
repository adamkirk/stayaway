package organisations

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/mutex"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
)

type UpdateCommand struct {
	ID   string
	Name *string `validate:"omitnil,min=3"`
	Slug *string `validate:"omitnil,required,min=3,slug"`
}

func (svc *Service) Update(cmd UpdateCommand) (*Organisation, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	org, err := svc.repo.Get(cmd.ID)

	if err != nil {
		return nil, err
	}

	if cmd.Name != nil {
		org.Name = *cmd.Name
	}

	if cmd.Slug != nil {
		orgBySlug, err := svc.repo.BySlug(*cmd.Slug)

		if orgBySlug != nil && orgBySlug.ID != org.ID {
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

		org.Slug = *cmd.Slug
	}

	slugMutexKey := slugMutexKey(*cmd.Slug)
	editLockKey := fmt.Sprintf("organisation_edit:%s", org.ID)

	l, err := svc.mutex.MultiClaimWithBackOff([]string{editLockKey, slugMutexKey}, 300*time.Millisecond)

	if err != nil {
		if _, ok := err.(mutex.ErrLockNotClaimed); ok {
			return nil, common.ErrConflict{
				Message: "organisation is already being edited elsewhere",
			}
		}

		return nil, err
	}

	defer func() {
		if err := l.Release(); err != nil {
			slog.Error("failed to release lock", "error", err, "key", editLockKey)
		}
	}()

	return svc.repo.Save(org)
}
