package organisations

import "github.com/adamkirk-stayaway/organisations/internal/model"


type GetHandlerRepo interface {
	Get(id string) (*model.Organisation, error)
}

type GetCommand struct {
	ID string
}

type GetHandler struct {
	repo GetHandlerRepo
}

func (h *GetHandler) Handle(cmd GetCommand) (*model.Organisation, error) {
	return h.repo.Get(cmd.ID)
}

func NewGetHandler(repo GetHandlerRepo) *GetHandler {
	return &GetHandler{
		repo: repo,
	}
}