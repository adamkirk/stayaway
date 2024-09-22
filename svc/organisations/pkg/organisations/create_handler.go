package organisations

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type CreateHandlerRepo interface {
	Save(org *model.Organisation) (*model.Organisation, error)
	BySlug(slug string) (*model.Organisation, error)
}

type CreateCommand struct {
	Name *string `validate:"required,min=3"`
	Slug *string `validate:"required,min=3,slug"`
}

type CreateHandler struct {
	repo CreateHandlerRepo
	validator Validator
}

func (h *CreateHandler) Handle(cmd CreateCommand) (*model.Organisation, error) {
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
		if _, ok := err.(model.ErrNotFound); !ok {
			return nil, err
		}
	}

	org := &model.Organisation{
		Name: *cmd.Name,
		Slug: *cmd.Slug,
	}
	
	return h.repo.Save(org)
}

func NewCreateHandler(repo CreateHandlerRepo, validator Validator) *CreateHandler {
	return &CreateHandler{
		repo: repo,
		validator: validator,
	}
}