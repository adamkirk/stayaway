package venues

type GetHandlerRepo interface {
	Get(id string, orgId string) (*Venue, error)
}

type GetCommand struct {
	ID             string `validate:"required"`
	OrganisationID string `validate:"required"`
}

type GetHandler struct {
	validator Validator
	repo      GetHandlerRepo
}

func (h *GetHandler) Handle(cmd GetCommand) (*Venue, error) {
	err := h.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	return h.repo.Get(cmd.ID, cmd.OrganisationID)
}

func NewGetHandler(validator Validator, repo GetHandlerRepo) *GetHandler {
	return &GetHandler{
		validator: validator,
		repo:      repo,
	}
}
