package accommodations

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/templates"
)

type GetCommand struct {
	ID *string `validate:"required"`
	VenueID         *string `validate:"required"`
}

func (svc *Service) Get(cmd GetCommand) (*Accommodation, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	a, err := svc.repo.Get(*cmd.ID, *cmd.VenueID)

	if err != nil {
		return nil, err
	}

	var template *templates.VenueTemplate

	if a.VenueTemplateID != nil {
		template, err = svc.templatesRepo.Get(*a.VenueTemplateID, *cmd.VenueID)

		if err != nil {
			// Maybe wanna wrap this error if its not found...it shouldn't ever 
			// be possible, but should do something smarter than returning a not
			// found
			return nil, err
		}
	}

	svc.decorateWithTemplateConfig(template, a)

	return a, nil
}
