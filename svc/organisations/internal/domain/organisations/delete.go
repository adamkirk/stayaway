package organisations

type DeleteCommand struct {
	ID string
}

func (svc *Service) Delete(cmd DeleteCommand) error {
	org, err := svc.repo.Get(cmd.ID)

	if err != nil {
		return err
	}

	return svc.repo.Delete(org)
}
