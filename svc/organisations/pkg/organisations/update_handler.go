package organisations

import "github.com/adamkirk-stayaway/organisations/pkg/model"


type UpdateHandlerRepo interface {
	Get(id string) (*model.Organisation, error)
	Save(org *model.Organisation) (*model.Organisation, error)
}

type UpdateCommand struct {
	ID string
	Name *string
	Slug *string
}

type UpdateHandler struct {
	repo UpdateHandlerRepo
}

func (h *UpdateHandler) Handle(cmd UpdateCommand) (*model.Organisation, error) {
	org, err := h.repo.Get(cmd.ID)

	if err != nil {
		return nil, err
	}
	
	if cmd.Name != nil {
		org.Name = *cmd.Name
	}

	if cmd.Slug != nil {
		org.Slug = *cmd.Slug
	}

	return h.repo.Save(org);
}

func NewUpdateHandler(repo UpdateHandlerRepo) *UpdateHandler {
	return &UpdateHandler{
		repo: repo,
	}
}