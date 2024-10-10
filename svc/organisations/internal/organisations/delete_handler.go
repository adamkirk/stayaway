package organisations

import "github.com/adamkirk-stayaway/organisations/internal/model"


type DeleteHandlerRepo interface {
	Get(id string) (*model.Organisation, error)
	Delete(*model.Organisation) error
}

type DeleteCommand struct {
	ID string
}

type DeleteHandler struct {
	repo DeleteHandlerRepo
}

func (h *DeleteHandler) Handle(cmd DeleteCommand) (error) {
	org, err := h.repo.Get(cmd.ID)

	if err != nil {
		return err
	}

	return h.repo.Delete(org)
}

func NewDeleteHandler(repo DeleteHandlerRepo) *DeleteHandler {
	return &DeleteHandler{
		repo: repo,
	}
}