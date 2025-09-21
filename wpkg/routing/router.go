// pkg/routing/router.go
package routing

import (
	"github.com/aasoft24/golara/wpkg/gola"

	"net/http"
	"regexp"
	"strings"

	"github.com/aasoft24/golara/wpkg/middleware"
)

// MiddlewareFunc defines middleware signature
type MiddlewareFunc func(func(ctx *gola.Context)) func(ctx *gola.Context)

// Route struct for each route
type Route struct {
	method      string
	pattern     *regexp.Regexp
	paramNames  []string
	handler     func(ctx *gola.Context)
	middlewares []MiddlewareFunc
}

// Router struct
type Router struct {
	routes         *[]Route
	middleware     []MiddlewareFunc
	TemplateEngine *gola.Context
	prefix         string
}

// NewRouter creates a new router
func NewRouter(templateEngine *gola.Context) *Router {
	routes := []Route{}
	return &Router{
		routes:         &routes,
		TemplateEngine: templateEngine,
	}
}

// AddRoute adds a route with pattern
func (r *Router) AddRoute(method string, pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	// apply group prefix
	fullPattern := r.prefix + pattern

	// Convert :param to regex
	paramNames := []string{}
	regexPattern := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	replacer := regexPattern.ReplaceAllStringFunc(fullPattern, func(m string) string {
		paramNames = append(paramNames, m[1:])
		return "([^/]+)"
	})
	regex := regexp.MustCompile("^" + replacer + "$")

	*r.routes = append(*r.routes, Route{
		method:      method,
		pattern:     regex,
		paramNames:  paramNames,
		handler:     handler,
		middlewares: middlewares,
	})
}

// ==== HTTP Methods ==== //
func (r *Router) Get(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	r.AddRoute("GET", pattern, handler, middlewares...)
}

func (r *Router) PostCSRF(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	// Auto add CSRF middleware for POST requests
	allMiddlewares := append([]MiddlewareFunc{middleware.CSRF}, middlewares...)
	r.AddRoute("POST", pattern, handler, allMiddlewares...)
}

func (r *Router) Post(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	if shouldAddCSRF(pattern) {
		middlewares = append([]MiddlewareFunc{middleware.CSRF}, middlewares...)
	}
	r.AddRoute("POST", pattern, handler, middlewares...)
}

func (r *Router) PostNoCSRF(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	// CSRF without middleware
	r.AddRoute("POST", pattern, handler, middlewares...)
}

func (r *Router) Put(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	// Auto add CSRF middleware for PUT requests
	if shouldAddCSRF(pattern) {
		middlewares = append([]MiddlewareFunc{middleware.CSRF}, middlewares...)
	}
	r.AddRoute("PUT", pattern, handler, middlewares...)
}

func (r *Router) Patch(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	// Auto add CSRF middleware for PATCH requests
	//allMiddlewares := append([]MiddlewareFunc{middleware.CSRF}, middlewares...)
	if shouldAddCSRF(pattern) {
		middlewares = append([]MiddlewareFunc{middleware.CSRF}, middlewares...)
	}
	r.AddRoute("PATCH", pattern, handler, middlewares...)
}

func (r *Router) Delete(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	r.AddRoute("DELETE", pattern, handler, middlewares...)
}

func (r *Router) Options(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	r.AddRoute("OPTIONS", pattern, handler, middlewares...)
}

func (r *Router) Head(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	r.AddRoute("HEAD", pattern, handler, middlewares...)
}

// Any will register handler for all methods
func (r *Router) Any(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	for _, m := range methods {
		r.AddRoute(m, pattern, handler, middlewares...)
	}
}

func (r *Router) PutNoCSRF(pattern string, handler func(ctx *gola.Context), middlewares ...MiddlewareFunc) {
	r.AddRoute("PUT", pattern, handler, middlewares...)
}

// Exclude multiple prefixes
func shouldAddCSRF(pattern string) bool {
	excludedPrefixes := []string{"/api/", "/webhook/", "/health/"}
	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(pattern, prefix) {
			return false
		}
	}
	return true
}

func (r *Router) ServeFiles(prefix string, dir string) {
	// Ensure prefix ends with a slash
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	fs := http.FileServer(http.Dir(dir))
	r.Get(prefix+":filepath", func(ctx *gola.Context) {
		http.StripPrefix(prefix, fs).ServeHTTP(ctx.Writer, ctx.Request)
	})
}

// ==== Group ==== //
func (r *Router) Group(prefix string, middlewares ...MiddlewareFunc) *Router {
	return &Router{
		routes:         r.routes, // share the same slice pointer
		middleware:     append(append([]MiddlewareFunc{}, r.middleware...), middlewares...),
		TemplateEngine: r.TemplateEngine,
		prefix:         r.prefix + prefix,
	}
}

func (r *Router) ServeStatic(urlPrefix string, folder string) {
	fs := http.FileServer(http.Dir(folder))
	cleanPrefix := urlPrefix
	if cleanPrefix[len(cleanPrefix)-1] != '/' {
		cleanPrefix += "/"
	}

	r.Get(cleanPrefix+":filepath", func(ctx *gola.Context) {
		// Strip the prefix and serve the file
		http.StripPrefix(cleanPrefix, fs).ServeHTTP(ctx.Writer, ctx.Request)
	})
}

// ==== Middleware ==== //
func (r *Router) Use(middleware MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware)
}

// ==== ServeHTTP ==== //
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	for _, route := range *r.routes {
		if route.method != method {
			continue
		}

		matches := route.pattern.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		// build context
		ctx := &gola.Context{
			Writer:         w,
			Request:        req,
			Params:         map[string]string{},
			TemplateEngine: r.TemplateEngine.TemplateEngine,
		}

		// extract params
		for i, name := range route.paramNames {
			ctx.Params[name] = matches[i+1]
		}

		handler := route.handler

		// apply route middleware
		for i := len(route.middlewares) - 1; i >= 0; i-- {
			handler = route.middlewares[i](handler)
		}

		// apply global middleware
		for i := len(r.middleware) - 1; i >= 0; i-- {
			handler = r.middleware[i](handler)
		}

		handler(ctx)
		return
	}

	http.NotFound(w, req)
}
