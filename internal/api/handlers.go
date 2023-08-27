package api

import (
	"errors"
	"fmt"
	"net/http"
	"shorty/internal/models"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

func HandleCreateShortening(shortener shortener) echo.HandlerFunc {
	type (
		request struct {
			ID  string `json:"id,omitempty" validate:"omitempty,alphanum"`
			URL string `json:"url" validate:"required,url"`
		}
		response struct {
			ID      string      `json:"id,omitempty"`
			URL     string      `json:"url,omitempty"`
			Message interface{} `json:"message,omitempty"`
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
			if errors.Is(err, models.ErrShorteningExists) {
				return echo.NewHTTPError(http.StatusConflict, response{
					Message: models.ErrShorteningExists.Error(),
				})
			}

			log.Errorf("error create shortening %v: %v", req.URL, err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, response{
			ID:  sh.ID,
			URL: fmt.Sprintf("https://%v/%v", c.Request().Host, sh.ID),
		})
	}
}

func HandleGetShortening(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		sh, err := shortener.Get(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrShorteningNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			log.Errorf("error get shortening %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "error get shortening")
		}

		return c.JSON(http.StatusOK, sh)
	}
}

func HandleRedirect(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		url, err := shortener.Redirect(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrShorteningNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			log.Errorf("error redirect shortening %v: %v", url, err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.Redirect(http.StatusMovedPermanently, url.String())
	}
}
