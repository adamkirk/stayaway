package venues

type GetCommand struct {
	ID             string `validate:"required"`
	OrganisationID string `validate:"required"`
}

func (svc *Service) Get(cmd GetCommand) (*Venue, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	return svc.repo.Get(cmd.ID, cmd.OrganisationID)
}
