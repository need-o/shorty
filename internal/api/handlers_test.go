package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shorty/internal/models"
	"shorty/internal/shorty"
	"shorty/internal/storage/memory"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCreateShorty(t *testing.T) {
	type (
		response struct {
			ID      string `json:"id,omitempty"`
			Address string `json:"address,omitempty"`
		}
	)

	t.Run("create shorty by valid url", func(t *testing.T) {
		var body = `{"url": "https://examplt.com"}`
		var resp response

		recorder := httptest.NewRecorder()
		shortener := shorty.New(memory.NewShortyStorage())
		handler := HandleCreateShorty(shortener)
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

		e := echo.New()
		c := e.NewContext(request, recorder)

		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, recorder.Code)

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))
		assert.NotEmpty(t, resp.ID)
		assert.NotEmpty(t, resp.Address)
	})

	t.Run("create shorty by existing id", func(t *testing.T) {
		var body = `{"id": "test", "url": "https://example.com"}`
		var err *echo.HTTPError

		recorder := httptest.NewRecorder()
		shortener := shorty.New(memory.NewShortyStorage())
		handler := HandleCreateShorty(shortener)
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

		e := echo.New()
		c := e.NewContext(request, recorder)

		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, recorder.Code)

		recorder = httptest.NewRecorder()
		request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c = e.NewContext(request, recorder)

		require.ErrorAs(t, handler(c), &err)
		assert.Equal(t, err.Code, http.StatusConflict)
		assert.Contains(t, err.Message, models.ErrShortyExists.Error())
	})

	t.Run("create shorty by invalid data", func(t *testing.T) {
		var body = `{"url": ""}`
		var err *echo.HTTPError

		recorder := httptest.NewRecorder()
		shortener := shorty.New(memory.NewShortyStorage())
		handler := HandleCreateShorty(shortener)
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

		e := echo.New()
		c := e.NewContext(request, recorder)

		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.ErrorAs(t, handler(c), &err)
		assert.Equal(t, err.Code, http.StatusBadRequest)
		assert.Contains(t, err.Message, "Field validation for 'URL' failed on the 'required' tag")
	})
}

func TestHandleGetShorty(t *testing.T) {
	ctx := context.Background()

	t.Run("get shorty by valid id", func(t *testing.T) {
		var resp models.Shorty

		recorder := httptest.NewRecorder()
		shortener := shorty.New(memory.NewShortyStorage())
		handler := HandleGetShorty(shortener)
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		input := models.ShortyInput{
			ID:  "test",
			URL: "https://example.com",
		}

		e := echo.New()
		c := e.NewContext(request, recorder)
		c.SetPath("/shorty/:id")
		c.SetParamNames("id")
		c.SetParamValues(input.ID)

		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		_, err := shortener.Create(ctx, input)
		assert.NoError(t, err)

		require.NoError(t, handler(c))
		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

		assert.Equal(t, resp.ID, input.ID)
		assert.Equal(t, resp.URL, input.URL)
		assert.True(t, len(resp.Visits) == 0)
	})

	t.Run("get shorty by not existing id", func(t *testing.T) {
		var err *echo.HTTPError

		recorder := httptest.NewRecorder()
		shortener := shorty.New(memory.NewShortyStorage())
		handler := HandleGetShorty(shortener)
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		input := models.ShortyInput{
			ID:  "test",
			URL: "https://example.com",
		}

		e := echo.New()
		c := e.NewContext(request, recorder)
		c.SetPath("/shorty/:id")
		c.SetParamNames("id")
		c.SetParamValues(input.ID)

		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.ErrorAs(t, handler(c), &err)
		assert.Equal(t, err.Code, http.StatusNotFound)
	})
}

func TestHandleRedirect(t *testing.T) {
	ctx := context.Background()

	t.Run("redirect by id", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		shortener := shorty.New(memory.NewShortyStorage())
		handler := HandleRedirect(shortener)
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		input := models.ShortyInput{
			ID:  "test",
			URL: "https://example.com",
		}

		e := echo.New()
		c := e.NewContext(request, recorder)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(input.ID)

		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		_, err := shortener.Create(ctx, input)
		assert.NoError(t, err)

		require.NoError(t, handler(c))
		assert.Equal(t, recorder.Code, http.StatusMovedPermanently)
		assert.Equal(t, recorder.Header().Get("Location"), input.URL)
	})

	t.Run("redirect by not existing id", func(t *testing.T) {
		var err *echo.HTTPError

		recorder := httptest.NewRecorder()
		shortener := shorty.New(memory.NewShortyStorage())
		handler := HandleRedirect(shortener)
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		e := echo.New()
		c := e.NewContext(request, recorder)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("sdfsdf")

		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.ErrorAs(t, handler(c), &err)
		assert.Equal(t, err.Code, http.StatusNotFound)
	})
}
