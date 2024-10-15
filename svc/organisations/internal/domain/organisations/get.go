package organisations

type GetCommand struct {
	ID string
}

func (svc *Service) Get(cmd GetCommand) (*Organisation, error) {
	return svc.repo.Get(cmd.ID)
}
