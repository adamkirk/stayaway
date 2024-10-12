package organisations

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/mutex"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
)

type CreateHandlerRepo interface {
	Save(org *Organisation) (*Organisation, error)
	BySlug(slug string) (*Organisation, error)
}

type CreateCommand struct {
	Name *string `validate:"required,min=3"`
	Slug *string `validate:"required,min=3,slug"`
}

type CreateHandler struct {
	repo CreateHandlerRepo
	validator Validator
	mutex DistributedMutex
}

func (h *CreateHandler) Handle(cmd CreateCommand) (*Organisation, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	

	orgBySlug, err := h.repo.BySlug(*cmd.Slug)

	if orgBySlug != nil {
		return nil, validation.ValidationError{
			Errs:[]validation.FieldError{
				{
					Key: "Slug",
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
	l, err := h.mutex.ClaimWithBackOff(slugMutexKey, 300 * time.Millisecond)

	if err != nil {
		if _, ok:= err.(mutex.ErrLockNotClaimed); ok {
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
		Name: *cmd.Name,
		Slug: *cmd.Slug,
	}
	
	return h.repo.Save(org)
}

func NewCreateHandler(repo CreateHandlerRepo, validator Validator, mutex DistributedMutex) *CreateHandler {
	return &CreateHandler{
		repo: repo,
		validator: validator,
		mutex: mutex,
	}
}