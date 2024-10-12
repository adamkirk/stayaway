package organisations


type DeleteHandlerRepo interface {
	Get(id string) (*Organisation, error)
	Delete(*Organisation) error
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