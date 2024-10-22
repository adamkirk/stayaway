package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/adamkirk-stayaway/organisations/internal/api/operations"
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/requests"
	"github.com/adamkirk-stayaway/organisations/internal/api/v1/responses"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/organisations"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/danielgtaylor/huma/v2"
)

type OrganisationsService interface {
	Get(organisations.GetCommand) (*organisations.Organisation, error)
	Create(cmd organisations.CreateCommand) (*organisations.Organisation, error)
	Delete(cmd organisations.DeleteCommand) error
	List(cmd organisations.ListCommand) (organisations.Organisations, common.PaginationResult, error)
	Update(cmd organisations.UpdateCommand) (*organisations.Organisation, error)
}

type OrganisationsControllerConfig interface{
	ApiServerDebugErrorsEnabled() bool 
}

type OrganisationsController struct {
	cfg              OrganisationsControllerConfig
	svc              OrganisationsService
	validationMapper *validation.ValidationMapper
}

func (c *OrganisationsController) RegisterRoutes(api huma.API) {
	huma.Register[requests.ListOrganisationsRequest, responses.ListOrganisationsResponse](api, huma.Operation{
		OperationID:  "v1.organisations.list",
		Method:       http.MethodGet,
		Path:         "/organisations",
		Summary:      "List all organisations",
		DefaultStatus: http.StatusOK,
		Metadata: map[string]any{
			operations.OptDisableNotFound: true,
		},
	}, ErrorHandler(c.cfg.ApiServerDebugErrorsEnabled(), c.validationMapper, c.List))

	huma.Register[requests.PostOrganisationRequest, responses.OrganisationResponse](api, huma.Operation{
		OperationID:  "v1.organisations.post",
		Method:       http.MethodPost,
		Path:         "/organisations",
		Summary:      "Create an organisation.",
		DefaultStatus: http.StatusCreated,
	}, ErrorHandler(c.cfg.ApiServerDebugErrorsEnabled(), c.validationMapper, c.Post))

	huma.Register[requests.PatchOrganisationRequest, responses.OrganisationResponse](api, huma.Operation{
		OperationID:  "v1.organisations.patch",
		Method:       http.MethodPatch,
		Path:         "/organisations/{id}",
		Summary:      "Patch an organisation.",
		DefaultStatus: http.StatusOK,
	}, ErrorHandler(c.cfg.ApiServerDebugErrorsEnabled(), c.validationMapper, c.Patch))

	huma.Register[requests.DeleteOrganisationRequest, responses.NoContent](api, huma.Operation{
		OperationID:  "v1.organisations.delete",
		Method:       http.MethodDelete,
		Path:         "/organisations/{id}",
		Summary:      "Delete an organisation.",
		DefaultStatus: http.StatusNoContent,
	}, ErrorHandler(c.cfg.ApiServerDebugErrorsEnabled(), c.validationMapper, c.Delete))

	huma.Register[requests.GetOrganisationRequest, responses.OrganisationResponse](api, huma.Operation{
		OperationID:  "v1.organisations.get",
		Method:       http.MethodGet,
		Path:         "/organisations/{id}",
		Summary:      "Get an organisation.",
		DefaultStatus: http.StatusOK,
	}, ErrorHandler(c.cfg.ApiServerDebugErrorsEnabled(), c.validationMapper, c.Get))
}


func NewOrganisationsController(
	cfg OrganisationsControllerConfig,
	svc OrganisationsService,
	validationMapper *validation.ValidationMapper,
) *OrganisationsController {
	return &OrganisationsController{
		cfg:              cfg,
		svc:              svc,
		validationMapper: validationMapper,
	}
}

func (c *OrganisationsController) List(ctx context.Context, req *requests.ListOrganisationsRequest) (*responses.ListOrganisationsResponse, error) {
	cmd := req.ToCommand()

	results, pagination, err := c.svc.List(cmd)

	if err != nil {
		return nil, err
	}

	resp := &responses.ListOrganisationsResponse{
		Body: responses.OrganisationsListResponseBody{
			Meta: responses.ListResponseMeta{
				SortOptionsResponseMeta: responses.SortOptionsResponseMeta{
					OrderDirection: string(cmd.OrderDirection),
					OrderBy:        string(cmd.OrderBy),
				},
				PaginationResponseMeta: responses.PaginationResponseMeta{
					Page:         pagination.Page,
					PerPage:      pagination.PerPage,
					TotalPages:   pagination.TotalPages,
					TotalResults: pagination.Total,
				},
			},
			Data: responses.OrganisationsFromModels(results),
		},
	}

	return resp, nil
}

func (c *OrganisationsController) Post(ctx context.Context, req *requests.PostOrganisationRequest) (*responses.OrganisationResponse, error) {
	org, err := c.svc.Create(req.ToCommand())

	if err != nil {
		slog.Debug("err from command")
		return nil, err
	}

	resp := &responses.OrganisationResponse{
		Body: responses.OrganisationResponseBody{
			Data: responses.OrganisationFromModel(org),
		},
	}

	return resp, nil
}

func (c *OrganisationsController) Get(ctx context.Context, req *requests.GetOrganisationRequest) (*responses.OrganisationResponse, error) {
	org, err := c.svc.Get(req.ToCommand())

	if err != nil {
		return nil, err
	}

	resp := &responses.OrganisationResponse{
		Body: responses.OrganisationResponseBody{
			Data: responses.OrganisationFromModel(org),
		},
	}

	return resp, nil
}

func (c *OrganisationsController) Patch(ctx context.Context, req *requests.PatchOrganisationRequest) (*responses.OrganisationResponse, error) {
	org, err := c.svc.Update(req.ToCommand())

	if err != nil {
		return nil, err
	}

	resp := &responses.OrganisationResponse{
		Body: responses.OrganisationResponseBody{
			Data: responses.OrganisationFromModel(org),
		},
	}

	return resp, nil
}

func (c *OrganisationsController) Delete(ctx context.Context, req *requests.DeleteOrganisationRequest) (*responses.NoContent, error) {
	err := c.svc.Delete(req.ToCommand())

	if err != nil {
		return nil, err
	}

	return &responses.NoContent{
		Status: http.StatusNoContent,
	}, nil
}
