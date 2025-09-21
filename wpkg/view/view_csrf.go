// pkg/view/csrf.go
package view

import (
	"html/template"

	"github.com/aasoft24/golara/wpkg/helpers"
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
