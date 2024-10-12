package venues

type DeleteHandlerRepo interface {
	Get(id string, orgId string) (*Venue, error)
	Delete(v *Venue) error
}

type DeleteCommand struct {
	ID string `validate:"required"`
	OrganisationID string `validate:"required"`
}

type DeleteHandler struct {
	validator Validator
	repo DeleteHandlerRepo
}

func (h *DeleteHandler) Handle(cmd DeleteCommand) error {
	err := h.validator.Validate(cmd)

	if err != nil {
		return err
	}

	v, err := h.repo.Get(cmd.ID, cmd.OrganisationID)

	if err != nil {
		return err
	}

	return h.repo.Delete(v)
}

func NewDeleteHandler(validator Validator, repo DeleteHandlerRepo) *DeleteHandler {
	return &DeleteHandler{
		validator: validator,
		repo: repo,
	}
}