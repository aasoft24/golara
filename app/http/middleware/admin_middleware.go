package middleware

import (
	"net/http"

	"github.com/aasoft24/golara/wpkg/gola"
)

type UserProvider interface {
	User(*http.Request) (interface{}, error)
}

func SafeUser(u interface{}) map[string]interface{} {
	switch v := u.(type) {
	case map[string]interface{}:
		return v
	default:
		return map[string]interface{}{"ID": u}
	}
}

func AdminMiddleware(a UserProvider) func(*gola.Context, func(*gola.Context)) {
	return func(c *gola.Context, next func(*gola.Context)) {
		user, err := a.User(c.Request)
		if err != nil {
			c.Redirect("/login")
			return
		}

		// Check if user is admin (expects Role in the user map)
		safeUser := SafeUser(user)
		if role, ok := safeUser["Role"]; !ok || role != "admin" {
			c.Redirect("/")
			return
		}

		c.Set("User", safeUser)
		next(c)
	}
}
