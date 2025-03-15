package handler

import (
	"server/views/layouts"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pocketbase/pocketbase/core"
)

func Render(e *core.RequestEvent, status int, t templ.Component) error {
	e.Response.WriteHeader(status)
	e.Response.Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	isAlpine := e.Request.Header.Get("X-Requested-With") == "AlpineJS"
	route := e.Request.URL.Path

	if !isAlpine {
		return layouts.BaseLayout(route).Render(e.Request.Context(), e.Response)
	}

	return t.Render(e.Request.Context(), e.Response)
}
