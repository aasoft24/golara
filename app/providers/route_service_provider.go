package providers

import (
	"your_project/routes"

	"github.com/aasoft24/golara/wpkg/foundation"
	"github.com/aasoft24/golara/wpkg/routing"
	"github.com/aasoft24/golara/wpkg/view"
)

type RouteServiceProvider struct {
	router         *routing.Router
	templateEngine *view.TemplateEngine
}

func NewRouteServiceProvider(router *routing.Router, templateEngine *view.TemplateEngine) *RouteServiceProvider {
	return &RouteServiceProvider{
		router:         router,
		templateEngine: templateEngine,
	}
}

// Bind router to app container
func (p *RouteServiceProvider) Register(app *foundation.Application) {
	app.Bind((*routing.Router)(nil), p.router)
}

// Boot routes
func (p *RouteServiceProvider) Boot(app *foundation.Application) {
	routes.RegisterWebRoutes(p.router, p.templateEngine)
	routes.RegisterApiRoutes(p.router)
}
