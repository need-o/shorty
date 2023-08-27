package api

import (
	"errors"
	"fmt"
	"net/http"
	"shorty/internal/models"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

func HandleCreateShorty(shortener shortener) echo.HandlerFunc {
	type (
		request struct {
			ID  string `json:"id,omitempty" validate:"omitempty,alphanum"`
			URL string `json:"url" validate:"required,url"`
		}
		response struct {
			ID      string `json:"id,omitempty"`
			Address string `json:"address,omitempty"`
		}
	)

	return func(c echo.Context) error {
		req := request{}

		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(req); err != nil {
			return err
		}

		input := models.ShortyInput{
			ID:  req.ID,
			URL: req.URL,
		}

		sh, err := shortener.Create(c.Request().Context(), input)
		if err != nil {
			if errors.Is(err, models.ErrShortyExists) {
				return echo.NewHTTPError(http.StatusConflict, models.ErrShortyExists.Error())
			}

			log.Errorf("error create shorty %v: %v", req.URL, err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, response{
			ID:      sh.ID,
			Address: fmt.Sprintf("%v/%v", c.Request().Host, sh.ID),
		})
	}
}

func HandleGetShorty(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		sh, err := shortener.Get(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrShortyNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			log.Errorf("error get shorty %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "error get shorty")
		}

		return c.JSON(http.StatusOK, sh)
	}
}

func HandleRedirect(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		visit := models.VisitInput{
			ShortyID:  id,
			Referer:   c.Request().Referer(),
			UserIP:    c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}

		url, err := shortener.Redirect(c.Request().Context(), visit)
		if err != nil {
			if errors.Is(err, models.ErrShortyNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			log.Errorf("error redirect shorty %v: %v", url, err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if len(c.QueryParams()) > 0 {
			url.RawQuery = c.QueryParams().Encode()
		}

		return c.Redirect(http.StatusMovedPermanently, url.String())
	}
}
