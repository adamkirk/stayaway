package v1

import (
	"github.com/labstack/echo/v4"
)

type ProbesController struct {}

func (c *ProbesController) RegisterRoutes(api *echo.Group) {
	g := api.Group("/_probes")
	g.GET("/startup", c.Startup).Name = "v1.probes.startup"
}

func NewProbesController(
) *ProbesController {
	return &ProbesController{}
}

// @Summary		Check is the app is listening
// @Tags			Probes
// @Accept			json
// @Produce		json
// @Success		204
// @Failure		502
// @Failure		504
// @Router			/v1/_probes/startup [get]
func (c *ProbesController) Startup(ctx echo.Context) error {
	ctx.NoContent(204)

	return nil
}
