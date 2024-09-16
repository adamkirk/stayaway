package api

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/labstack/echo/v4"
)


type OrganisationsRepo interface {
	Paginate(orderBy model.OrganisationSortBy, orderDir model.SortDirection, page int, perPage int) (model.Organisations, model.PaginationResult, error)
	Save(org *model.Organisation) (*model.Organisation, error)
	Get(id string) (*model.Organisation, error)
	Delete(org *model.Organisation) (error)
}

type OrganisationsV1ControllerConfig interface {}

type OrganisationsV1Controller struct {
	cfg OrganisationsV1ControllerConfig
	repo OrganisationsRepo
}

func (c *OrganisationsV1Controller) RegisterRoutes(api *echo.Group) {
	g := api.Group("/v1/organisations")
	g.POST("", c.Create).Name = "v1.organisations.create"
	g.DELETE("/:id", c.Delete).Name = "v1.organisations.delete"
	g.GET("/:id", c.Get).Name = "v1.organisations.get"
	g.PATCH("/:id", c.Patch).Name = "v1.organisations.patch"
	g.GET("", c.List).Name = "v1.organisations.list"
}

func NewOrganisationsV1Controller(cfg OrganisationsV1ControllerConfig, repo OrganisationsRepo) *OrganisationsV1Controller {
	return &OrganisationsV1Controller{
		cfg: cfg,
		repo: repo,
	}
}

func (c *OrganisationsV1Controller) List(ctx echo.Context) error {
	req := V1ListOrganisationsRequest{
		OrderDirection: model.SortAsc,
		OrderBy: model.OrganisationSortBySlug,
		Page: 1,
		PerPage: 50,
	}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	results, pagination, err := c.repo.Paginate(
		req.OrderBy,
		req.OrderDirection,
		req.Page,
		req.PerPage,
	)

	if err != nil {
		return err
	}
	
	resp := V1ListOrganisationsResponse{
		Meta: V1ListResponseMeta{
			V1SortOptionsResponseMeta: V1SortOptionsResponseMeta{
				OrderDirection: string(req.OrderDirection),
				OrderBy: string(req.OrderBy),
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

	// TODO validation...
	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org := &model.Organisation{
		Name: req.Name,
		Slug: req.Slug,
	}
	
	org, err := c.repo.Save(org)

	if err != nil {
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

	// TODO validation...
	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.repo.Get(req.ID)

	if err != nil {
		if notFound, ok := err.(model.ErrNotFound); ok {
			return ErrNotFound{
				ResourceName: notFound.ResourceName,
			}
		}

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

	// TODO validation...
	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.repo.Get(req.ID)

	if err != nil {
		if notFound, ok := err.(model.ErrNotFound); ok {
			return ErrNotFound{
				ResourceName: notFound.ResourceName,
			}
		}

		return err
	}
	
	if req.Name != nil {
		org.Name = *req.Name
	}

	if req.Slug != nil {
		org.Slug = *req.Slug
	}

	org, err = c.repo.Save(org);

	if err != nil {
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

	// TODO validation...
	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	org, err := c.repo.Get(req.ID)

	if err != nil {
		if notFound, ok := err.(model.ErrNotFound); ok {
			return ErrNotFound{
				ResourceName: notFound.ResourceName,
			}
		}

		return err
	}

	if err := c.repo.Delete(org); err != nil {
		return err
	}
	
	ctx.NoContent(204)

	return nil
}