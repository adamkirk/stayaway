package organisations

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
)

type CreateHandlerRepo interface {
	Save(org *model.Organisation) (*model.Organisation, error)
}

type CreateCommand struct {
	Name *string `validate:"nonnil"`
	Slug *string `validate:"nonnil"`
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

	org := &model.Organisation{
		Name: *cmd.Name,
		Slug: *cmd.Slug,
	}
	
	return org, nil
	// return h.repo.Save(org)
}

func NewCreateHandler(repo CreateHandlerRepo, validator Validator) *CreateHandler {
	return &CreateHandler{
		repo: repo,
		validator: validator,
	}
}