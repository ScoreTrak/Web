package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// RouteHandler is a default handler to serve up
// a route page.
func RouteHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("routes.html"))
}
