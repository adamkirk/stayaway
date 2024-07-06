package api

import (
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
	ctx.String(200, "List")

	return nil
}

func (c *OrganisationsV1Controller) Create(ctx echo.Context) error {
	ctx.String(200, "Create")

	return nil
}