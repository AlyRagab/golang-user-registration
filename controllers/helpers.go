package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

// parseForm to decode the http request in Gorilla schema
func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
