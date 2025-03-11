package handlers

import (
	"net/http"

	"github.com/utiiz/autodarts/views/templates"
)

func Index(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, templates.Index())
}
