package routes

import (
	"github.com/aasoft24/golara/wpkg/gola"
	"github.com/aasoft24/golara/wpkg/routing"
	"github.com/aasoft24/golara/wpkg/view"
)

func RegisterWebRoutes(router *routing.Router, templateEngine *view.TemplateEngine) {

	router.Get("/", func(c *gola.Context) {
		c.Session.Set("name", "Sayed")
		c.Session.Save()

		c.String(200, "Hello, World!")
	})
}
