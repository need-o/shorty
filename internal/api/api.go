package api

import (
	"context"
	"net/http"
	"net/url"
	"shorty/internal/models"

	"github.com/labstack/echo/v4"
)

type (
	shortener interface {
		Get(ctx context.Context, id string) (*models.Shortening, error)
		Create(ctx context.Context, in models.ShortyInput) (*models.Shortening, error)
		Redirect(ctx context.Context, id string) (*url.URL, error)
	}

	CloseFunc func(context.Context) error

	Api struct {
		echo      *echo.Echo
		shortener shortener
		closers   []CloseFunc
	}
)

func New(shortener shortener) *Api {
	api := Api{
		shortener: shortener,
	}

	api.router()
	api.AddCloser(api.echo.Shutdown)

	return &api
}

func (s *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.echo.ServeHTTP(w, r)
}

func (a *Api) AddCloser(closer CloseFunc) {
	a.closers = append(a.closers, closer)
}

func (a *Api) Shutdown(ctx context.Context) error {
	for _, close := range a.closers {
		if err := close(ctx); err != nil {
			return err
		}
	}

	return nil
}
