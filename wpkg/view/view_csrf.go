// pkg/view/csrf.go
package view

import (
	"html/template"
	"net/http"

	"github.com/aasoft24/golara/wpkg/helpers"
	"github.com/aasoft24/golara/wpkg/session"
)

// CSRF functions for templates
func CSRFunctions() template.FuncMap {
	return template.FuncMap{
		"csrf": func() string {
			return helpers.GenerateRandomString(32)
		},
		"csrf_field": func() template.HTML {
			token := helpers.GenerateRandomString(32)
			return template.HTML(`<input type="hidden" name="_token" value="` + token + `">`)
		},
		"csrf_meta": func() template.HTML {
			token := helpers.GenerateRandomString(32)
			return template.HTML(`<meta name="csrf-token" content="` + token + `">`)
		},
	}
}

// CSRF token generate korar function
func CSRFToken(r *http.Request, w http.ResponseWriter) string {
	// Session manager create korun
	mgr := session.NewManager(session.NewMemoryStore(), "myapp-session")

	// Session start korun
	sess, err := mgr.Start(w, r)
	if err != nil {
		return ""
	}

	// Existing token check korun
	token := sess.Get("csrf_token")
	if token != nil {
		return token.(string)
	}

	// New token generate korun
	newToken := helpers.GenerateRandomString(32)
	sess.Set("csrf_token", newToken)
	sess.Save()

	return newToken
}

// CSRF hidden field
func CSRFField(r *http.Request, w http.ResponseWriter) template.HTML {
	token := CSRFToken(r, w)
	return template.HTML(`<input type="hidden" name="_token" value="` + token + `">`)
}

// CSRF meta tag
func CSRFMetaTag(r *http.Request, w http.ResponseWriter) template.HTML {
	token := CSRFToken(r, w)
	return template.HTML(`<meta name="csrf-token" content="` + token + `">`)
}

func dict(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		panic("dict function expects even number of arguments")
	}
	m := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			panic("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m
}
