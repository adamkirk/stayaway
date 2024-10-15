package venues

type DeleteCommand struct {
	ID             string `validate:"required"`
	OrganisationID string `validate:"required"`
}

func (svc *Service) Delete(cmd DeleteCommand) error {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return err
	}

	v, err := svc.repo.Get(cmd.ID, cmd.OrganisationID)

	if err != nil {
		return err
	}

	return svc.repo.Delete(v)
}
