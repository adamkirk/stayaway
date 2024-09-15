package api

import (
	"github.com/adamkirk-stayaway/venues/pkg/model"
	"github.com/labstack/echo/v4"
)

type OrganisationsV1ControllerConfig interface {}

type OrganisationsV1Controller struct {
	cfg OrganisationsV1ControllerConfig
}

func (c *OrganisationsV1Controller) RegisterRoutes(api *echo.Group) {
	g := api.Group("/v1/organisations")
	g.POST("", c.Create).Name = "v1.organisations.create"
	g.GET("", c.List).Name = "v1.organisations.list"
}

func NewOrganisationsV1Controller(cfg OrganisationsV1ControllerConfig) *OrganisationsV1Controller {
	return &OrganisationsV1Controller{
		cfg: cfg,
	}
}

func (c *OrganisationsV1Controller) List(ctx echo.Context) error {
	req := V1ListOrganisationsRequest{
		OrderDirection: model.SortAsc,
		OrderBy: model.OrganisationOrderBySlug,
		Page: 1,
		PerPage: 50,
	}

	if err := bindRequest(&req, ctx); err != nil {
		return err
	}

	orgs := make(model.Organisations, 3)
	orgs[0] = &model.Organisation{
		ID: *model.NewID(),
		Name: "Test 1",
		Slug: "test-1",
	}
	orgs[1] = &model.Organisation{
		ID: *model.NewID(),
		Name: "Test 2",
		Slug: "test-2",
	}
	orgs[2] = &model.Organisation{
		ID: *model.NewID(),
		Name: "Test 3",
		Slug: "test-3",
	}
	
	resp := V1ListOrganisationsResponse{
		Meta: V1ListOrganisationsMeta{
			OrderDirection: string(req.OrderDirection),
			OrderBy: string(req.OrderBy),
			Page: req.Page,
			PerPage: req.PerPage,
			TotalPages: 0,
		},
		Data: V1OrganisationsFromModels(orgs),
	}

	ctx.JSON(200, resp)

	return nil
}

func (c *OrganisationsV1Controller) Create(ctx echo.Context) error {
	ctx.String(200, "Create")

	return nil
}