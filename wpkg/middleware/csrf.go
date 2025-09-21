// pkg/middleware/csrf.go
package middleware

import (
	"net/http"

	"github.com/aasoft24/golara/wpkg/gola"
)

// Correct middleware signature for your router
func CSRF(next func(ctx *gola.Context)) func(ctx *gola.Context) {
	return func(ctx *gola.Context) {
		// Skip safe methods
		if ctx.Request.Method == "GET" ||
			ctx.Request.Method == "HEAD" ||
			ctx.Request.Method == "OPTIONS" {
			next(ctx)
			return
		}

		// Check CSRF token
		token := ctx.Request.FormValue("_token")
		if token == "" {
			// Also check header for AJAX requests
			token = ctx.Request.Header.Get("X-CSRF-Token")
		}

		if token == "" || len(token) < 10 {
			if ctx.Request.Header.Get("Content-Type") == "application/json" {
				ctx.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "CSRF token missing or invalid",
				})
				return
			}
			//ctx.SetFlash("error", "CSRF token missing or invalid")
			ctx.HTML(http.StatusForbidden, "CSRF token missing or invalid")
			//next(ctx)

			return
		}

		next(ctx)
	}
}
