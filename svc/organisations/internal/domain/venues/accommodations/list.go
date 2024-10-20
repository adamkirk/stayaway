package accommodations

import (
	"log/slog"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/templates"
)

type ListCommand struct {
	OrganisationID string               `validate:"required"`
	VenueID        string               `validate:"required"`
	ReferencePrefix     *string              `validate:"omitnil,min=3"`
	VenueTemplateID *string               `validate:"omitnil"`
	OrderDirection common.SortDirection `validate:"required,orderdir"`
	OrderBy        SortBy               `validate:"required,venueaccommodation_sortfield"`
	Page           int                  `validate:"required,min=1"`
	PerPage        int                  `validate:"required,min=1"`
}

func NewListCommand() ListCommand {
	return ListCommand{
		OrderDirection: common.SortAsc,
		OrderBy:        SortByReference,
		Page:           1,
		PerPage:        50,
	}
}

func (svc *Service) List(cmd ListCommand) (Accommodations, common.PaginationResult, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	p := PaginationFilter{
		OrderBy:  cmd.OrderBy,
		OrderDir: cmd.OrderDirection,
		Page:     cmd.Page,
		PerPage:  cmd.PerPage,
	}

	s := SearchFilter{
		ReferencePrefix: cmd.ReferencePrefix,
		VenueID: []string{
			cmd.VenueID,
		},
	}

	if cmd.VenueTemplateID != nil {
		s.VenueTemplateID = []string{
			*cmd.VenueTemplateID,
		}
	}

	accs, pagination, err := svc.repo.Paginate(p, s)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}


	templs, err := svc.templatesRepo.ByID(accs.VenueTemplateIDs(), cmd.VenueID)

	slog.Debug("found templates for accommodations", "templates", templs, "ids", accs.VenueTemplateIDs())

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	templatesKV := templs.KeyByID()

	for _, a := range accs {
		var t *templates.VenueTemplate
		
		if a.VenueTemplateID != nil {
			t = templatesKV[*a.VenueTemplateID]
		}

		if err := svc.decorateWithTemplateConfig(t, a); err != nil {
			slog.Debug("failed to decorate accommodation", "accommodation", a, "template", templatesKV[*a.VenueTemplateID])
			return nil, common.PaginationResult{}, err
		}
	}

	return accs, pagination, nil
}
