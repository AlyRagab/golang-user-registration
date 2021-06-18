package controllers

import (
	"github.com/AlyRagab/golang-user-registration/views"
)

// Static struct for static views/pages
type Static struct {
	Home    *views.View
	Contact *views.View
}

// NewStatic for creating a static pages and parsing the template
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}
