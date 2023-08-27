package api

import (
	"context"
	"net/http"
	"shorty/internal/models"

	"github.com/labstack/echo/v4"
)

type (
	shortener interface {
		Get(ctx context.Context, id string) (*models.Shortening, error)
		Create(ctx context.Context, in models.InputShortening) (string, error)
		Redirect(ctx context.Context, id string) (*models.RedirectShortening, error)
	}

	Api struct {
		echo      *echo.Echo
		shortener shortener
	}
)

func New(shortener shortener) *Api {
	return &Api{
		echo:      echo.New(),
		shortener: shortener,
	}
}

func (s *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.echo.ServeHTTP(w, r)
}
