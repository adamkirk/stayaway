package templates

import (
	"log/slog"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
)

type CreateVenueTemplateCommand struct {
	OrganisationID *string `validate:"required"`
	VenueID        *string `validate:"required"`
	Name           *string `validate:"required,min=3"`
	Type           *string `validate:"required,accommodationtype"`
	MinOccupancy   *int    `validate:"required,min=1"`
	// 100 seems an appropriate max
	MaxOccupancy *int    `validate:"omitnil,gtefield=MinOccupancy"`
	Description  *string `validate:"required,min=10"`
}

func (svc *VenueTemplatesService) Create(cmd CreateVenueTemplateCommand) (*VenueTemplate, error) {
	err := svc.validator.Validate(cmd)

	if err != nil {
		return nil, err
	}

	_, err = svc.venuesRepo.Get(*cmd.VenueID, *cmd.OrganisationID)

	if err != nil {
		if _, ok := err.(common.ErrNotFound); ok {
			return nil, err
		}

		return nil, err
	}

	blah, err := svc.repo.ByNameAndVenue(*cmd.Name, *cmd.VenueID)

	slog.Debug("found", "vt", blah, "err", err)

	if err == nil {
		return nil, validation.ValidationError{
			Errs: []validation.FieldError{
				{
					Key:    "Name",
					Errors: []string{"must be unique"},
				},
			},
		}
	} else {
		_, ok := err.(common.ErrNotFound)

		if !ok {
			return nil, err
		}
	}

	vt := &VenueTemplate{
		ID:      svc.idGen.Generate(),
		VenueID: *cmd.VenueID,
		AccommodationTemplate: common.AccommodationTemplate{
			Name:         *cmd.Name,
			MinOccupancy: *cmd.MinOccupancy,
			MaxOccupancy: cmd.MaxOccupancy,
			Description:  *cmd.Description,
			Type:         common.AccommodationTemplateType(*cmd.Type),
		},
	}

	return svc.repo.Save(vt)
}
