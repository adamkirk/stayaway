package organisations

import "github.com/adamkirk-stayaway/organisations/pkg/model"


type CreateHandlerRepo interface {
	Save(org *model.Organisation) (*model.Organisation, error)
}

type CreateCommand struct {
	Name string
	Slug string
}

type CreateHandler struct {
	repo CreateHandlerRepo
}

func (h *CreateHandler) Handle(cmd CreateCommand) (*model.Organisation, error) {
	org := &model.Organisation{
		Name: cmd.Name,
		Slug: cmd.Slug,
	}
	
	return h.repo.Save(org)
}

func NewCreateHandler(repo CreateHandlerRepo) *CreateHandler {
	return &CreateHandler{
		repo: repo,
	}
}