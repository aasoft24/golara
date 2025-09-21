package facades

import (
	"github.com/aasoft24/golara/wpkg/gola"
	"github.com/aasoft24/golara/wpkg/helpers"

	"github.com/gorilla/sessions"
)

type AuthFacade struct{}

var Auth = &AuthFacade{}

var store = sessions.NewCookieStore([]byte("very-secret-key"))

type SafeUser struct {
	ID     uint
	Name   string
	Email  string
	Mobile string
}

// User returns SafeUser and ok flag
func (a *AuthFacade) User(c *gola.Context) (*helpers.SafeUser, bool) {
	user, ok := c.Get("User").(helpers.SafeUser)
	if !ok {
		return nil, false
	}
	return &user, true
}

// Id returns user ID
func (a *AuthFacade) Id(c *gola.Context) uint {
	if u, ok := a.User(c); ok {
		return u.ID
	}
	return 0
}

// Check if authenticated
func (a *AuthFacade) Check(c *gola.Context) bool {
	_, ok := a.User(c)
	return ok
}

// Guest
func (a *AuthFacade) Guest(c *gola.Context) bool {
	return !a.Check(c)
}

// Auth.attempt(email, password)
// func (a *AuthFacade) Attempt(c *gola.Context, identifier, password string) bool {
// 	var user models.User
// 	db := database.DB

// 	// Try finding user by email, mobile, or name
// 	if err := db.Where("email = ? OR mobile = ?", identifier, identifier).First(&user).Error; err != nil {
// 		return false
// 	}

// 	// Password check
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
// 		return false
// 	}

// 	// Set session safely
// 	if c != nil && c.Request != nil && c.Response != nil {
// 		session, _ := store.Get(c.Request, "user-session")
// 		session.Values["user"] = user.ID
// 		session.Save(c.Request, c.Response)

// 		// Set context
// 		c.Set("User", SafeUser{
// 			ID:     uint(user.ID),
// 			Name:   user.Name,
// 			Email:  user.Email,
// 			Mobile: user.Mobile,
// 		})
// 	}

// 	return true
// }

// Auth.logout()

// Logout
func (a *AuthFacade) Logout(c *gola.Context) {
	if c == nil {
		return
	}

	session, _ := store.Get(c.Request, "user-session")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)

	c.Set("User", nil)
}
