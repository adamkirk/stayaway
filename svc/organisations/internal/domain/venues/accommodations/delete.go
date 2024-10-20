package accommodations

type DeleteCommand struct {
	ID *string `validate:"required"`
	VenueID         *string `validate:"required"`
}

func (svc *Service) Delete(cmd DeleteCommand) (error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return err
	}

	a, err := svc.repo.Get(*cmd.ID, *cmd.VenueID)

	if err != nil {
		return err
	}

	return svc.repo.Delete(a)
}
