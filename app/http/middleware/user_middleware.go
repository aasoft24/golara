package middleware

import (
	"github.com/aasoft24/golara/app/models"
	"github.com/aasoft24/golara/wpkg/database"
	"github.com/aasoft24/golara/wpkg/gola"
	"github.com/aasoft24/golara/wpkg/helpers"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("very-secret-key"))

func UserMiddleware(next func(c *gola.Context)) func(c *gola.Context) {
	return func(c *gola.Context) {

		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")

		session, _ := store.Get(c.Request, "user-session")
		userID, ok := session.Values["user"]
		if !ok || userID == 0 {
			c.Redirect("/login")
			return
		}

		var user models.User
		db := database.DB
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			c.Redirect("/login")
			return
		}

		safeUser := helpers.SafeUser{
			ID:    uint(user.ID),
			Name:  user.Name,
			Email: user.Email,
		}
		c.Set("User", safeUser)

		next(c) // continue to next handler
	}
}
