package organisations

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/adamkirk-stayaway/organisations/internal/model"
	"github.com/adamkirk-stayaway/organisations/internal/mutex"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
)


type UpdateHandlerRepo interface {
	Get(id string) (*model.Organisation, error)
	Save(org *model.Organisation) (*model.Organisation, error)
	BySlug(slug string) (*model.Organisation, error)
}

type UpdateCommand struct {
	ID string
	Name *string `validate:"omitnil,min=3"`
	Slug *string `validate:"omitnil,required,min=3,slug"`
}

type UpdateHandler struct {
	repo UpdateHandlerRepo
	validator Validator
	mutex DistributedMutex
}

func (h *UpdateHandler) Handle(cmd UpdateCommand) (*model.Organisation, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	org, err := h.repo.Get(cmd.ID)

	if err != nil {
		return nil, err
	}
	

	
	if cmd.Name != nil {
		org.Name = *cmd.Name
	}

	slugMutexKey := slugMutexKey(*cmd.Slug)
	
	if cmd.Slug != nil {
		orgBySlug, err := h.repo.BySlug(*cmd.Slug)

		if orgBySlug != nil && orgBySlug.ID != org.ID {
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
			if _, ok := err.(model.ErrNotFound); !ok {
				return nil, err
			}
		}
		
		// claim a mutex so no one else can create/edit another resource to use 
		// this slug.
		l, err := h.mutex.ClaimWithBackOff(slugMutexKey, 300 * time.Millisecond)

		if err != nil {
			if _, ok:= err.(mutex.ErrLockNotClaimed); ok {
				return nil, model.ErrConflict{
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

		org.Slug = *cmd.Slug
	}

	editLockKey := fmt.Sprintf("organisation_edit:%s", org.ID)

	l, err := h.mutex.ClaimWithBackOff(editLockKey, 300 * time.Millisecond)
	
	if err != nil {
		if _, ok:= err.(mutex.ErrLockNotClaimed); ok {
			return nil, model.ErrConflict{
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

	return h.repo.Save(org);
}

func NewUpdateHandler(repo UpdateHandlerRepo, validator Validator, mutex DistributedMutex) *UpdateHandler {
	return &UpdateHandler{
		repo: repo,
		validator: validator,
		mutex: mutex,
	}
}