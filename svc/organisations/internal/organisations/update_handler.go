package organisations

import (
	"github.com/adamkirk-stayaway/organisations/internal/model"
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

		org.Slug = *cmd.Slug
	}

	return h.repo.Save(org);
}

func NewUpdateHandler(repo UpdateHandlerRepo, validator Validator) *UpdateHandler {
	return &UpdateHandler{
		repo: repo,
		validator: validator,
	}
}