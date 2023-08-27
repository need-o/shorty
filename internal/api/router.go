package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *Api) router() {
	a.echo = echo.New()
	a.echo.Use(middleware.RequestID())
	a.echo.Validator = NewValidator()

	api := a.echo.Group("/api")
	api.POST("/shortening", HandleCreateShortening(a.shortener))
	api.GET("/shortening/:id", HandleGetShortening(a.shortener))

	a.echo.GET("/:id", HandleRedirect(a.shortener))
}
