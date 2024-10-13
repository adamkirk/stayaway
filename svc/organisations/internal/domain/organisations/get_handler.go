package organisations

type GetHandlerRepo interface {
	Get(id string) (*Organisation, error)
}

type GetCommand struct {
	ID string
}

type GetHandler struct {
	repo GetHandlerRepo
}

func (h *GetHandler) Handle(cmd GetCommand) (*Organisation, error) {
	return h.repo.Get(cmd.ID)
}

func NewGetHandler(repo GetHandlerRepo) *GetHandler {
	return &GetHandler{
		repo: repo,
	}
}
