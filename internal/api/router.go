package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func (a *Api) router() {
	a.echo = echo.New()
	a.echo.Use(middleware.RequestID())
	a.echo.Use(requestLogMiddleware())
	a.echo.Validator = NewValidator()

	api := a.echo.Group("/api")
	api.POST("/shortening", HandleCreateShortening(a.shortener))
	api.GET("/shortening/:id", HandleGetShortening(a.shortener))

	a.echo.GET("/:id", HandleRedirect(a.shortener))
}

func requestLogMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	})
}
