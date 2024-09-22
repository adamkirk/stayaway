package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/organisations"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/labstack/echo/v4"
)


type OrganisationsRepo interface {
	Paginate(orderBy model.OrganisationSortBy, orderDir model.SortDirection, page int, perPage int) (model.Organisations, model.PaginationResult, error)
	Save(org *model.Organisation) (*model.Organisation, error)
	Get(id string) (*model.Organisation, error)
	Delete(org *model.Organisation) (error)
}

type OrganisationsGetHandler interface {
	Handle(organisations.GetCommand) (*model.Organisation, error)
}

type OrganisationsListHandler interface {
	Handle(cmd organisations.ListCommand) (model.Organisations, model.PaginationResult, error)
}

type OrganisationsCreateHandler interface {
	Handle(cmd organisations.CreateCommand) (*model.Organisation, error)
}

type OrganisationsDeleteHandler interface {
	Handle(cmd organisations.DeleteCommand) error
}


type OrganisationsUpdateHandler interface {
	Handle(cmd organisations.UpdateCommand) (*model.Organisation, error)
}

type OrganisationsV1ControllerConfig interface {}

type OrganisationsV1Controller struct {
	cfg OrganisationsV1ControllerConfig
	repo OrganisationsRepo
	get OrganisationsGetHandler
	list OrganisationsListHandler
	create OrganisationsCreateHandler
	delete OrganisationsDeleteHandler
	update OrganisationsUpdateHandler
	validationMapper *ValidationMapper
}

func (c *OrganisationsV1Controller) RegisterRoutes(api *echo.Group) {
	g := api.Group("/v1/organisations")
	g.POST("", c.Create).Name = "v1.organisations.create"
	g.DELETE("/:id", c.Delete).Name = "v1.organisations.delete"
	g.GET("/:id", c.Get).Name = "v1.organisations.get"
	g.PATCH("/:id", c.Patch).Name = "v1.organisations.patch"
	g.GET("", c.List).Name = "v1.organisations.list"
}

func NewOrganisationsV1Controller(
	cfg OrganisationsV1ControllerConfig,
	repo OrganisationsRepo,
	get OrganisationsGetHandler,
	list OrganisationsListHandler,
	create OrganisationsCreateHandler,
	delete OrganisationsDeleteHandler,
	update OrganisationsUpdateHandler,
	validationMapper *ValidationMapper,
) *OrganisationsV1Controller {
	return &OrganisationsV1Controller{
		cfg: cfg,
		repo: repo,
		get: get,
		list: list,
		create: create,
		delete: delete,
		update: update,
		validationMapper: validationMapper,
	}
}

func (c *OrganisationsV1Controller) List(ctx echo.Context) error {
	req := V1ListOrganisationsRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	cmd := req.ToCommand()

	results, pagination, err := c.list.Handle(cmd)

	if err != nil {
		return err
	}
	
	resp := V1ListOrganisationsResponse{
		Meta: V1ListResponseMeta{
			V1SortOptionsResponseMeta: V1SortOptionsResponseMeta{
				OrderDirection: string(cmd.OrderDirection),
				OrderBy: string(cmd.OrderBy),
			},
			V1PaginationResponseMeta: V1PaginationResponseMeta{
				Page: pagination.Page,
				PerPage: pagination.PerPage,
				TotalPages: pagination.TotalPages,
				TotalResults: pagination.Total,
			},
		},
		Data: V1OrganisationsFromModels(results),
	}

	ctx.JSON(200, resp)

	return nil
}

func (c *OrganisationsV1Controller) Create(ctx echo.Context) error {
	req := V1PostOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.create.Handle(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}
		return err
	}

	resp := V1PostOrganisationResponse{
		Data: V1OrganisationFromModel(org),
	}

	ctx.JSON(201, resp)

	return nil
}


func (c *OrganisationsV1Controller) Get(ctx echo.Context) error {
	req := V1GetOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.get.Handle(req.ToCommand())

	if err != nil {
		return err
	}

	resp := V1GetOrganisationResponse{
		Data: V1OrganisationFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

func (c *OrganisationsV1Controller) Patch(ctx echo.Context) error {
	req := V1PatchOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.update.Handle(req.ToCommand())

	if err != nil {
		if err, ok := err.(validation.ValidationError); ok {
			return c.validationMapper.Map(err, req)
		}

		return err
	}

	resp := V1PatchOrganisationResponse{
		Data: V1OrganisationFromModel(org),
	}

	ctx.JSON(200, resp)

	return nil
}

func (c *OrganisationsV1Controller) Delete(ctx echo.Context) error {
	req := V1DeleteOrganisationRequest{}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	err := c.delete.Handle(req.ToCommand())

	if err != nil {
		return err
	}
	
	ctx.NoContent(204)

	return nil
}