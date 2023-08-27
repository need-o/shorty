package shorty

import (
	"context"
	"shorty/internal/models"
	"shorty/internal/storage/memory"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewID(t *testing.T) {
	type test struct {
		name string
		in   uint32
		out  string
	}

	tests := []test{
		{
			name: "valid 0",
			in:   0,
			out:  "",
		},
		{
			name: "valid numbers",
			in:   123456789,
			out:  "HUawi",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := NewID(test.in)

			for i := 0; i < 100; i++ {
				out2 := NewID(test.in)
				assert.Equal(t, out, out2)

			}

			assert.Equal(t, out, test.out)
		})
	}
}

func TestShorty(t *testing.T) {
	type test struct {
		name  string
		input models.ShortyInput
		run   func(test, *Shorty)
	}

	ctx := context.Background()

	tests := []test{
		{
			name: "create with URL",
			input: models.ShortyInput{
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				sh, err := shorty.Create(ctx, test.input)

				require.NoError(t, err)
				assert.NotEmpty(t, sh.ID)
				assert.Equal(t, sh.URL, test.input.URL)
				assert.NotNil(t, sh.CreatedAt)
				assert.NotNil(t, sh.CreatedAt)
			},
		},
		{
			name: "create with ID and URL",
			input: models.ShortyInput{
				ID:  "test",
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				sh, err := shorty.Create(ctx, test.input)

				require.NoError(t, err)
				assert.Equal(t, sh.ID, test.input.ID)
				assert.Equal(t, sh.URL, test.input.URL)
				assert.NotZero(t, sh.CreatedAt)
				assert.NotZero(t, sh.CreatedAt)
			},
		},
		{
			name: "create with existing ID",
			input: models.ShortyInput{
				ID:  "test",
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				_, err := shorty.Create(ctx, test.input)
				_, errExisting := shorty.Create(ctx, test.input)

				require.NoError(t, err)
				assert.ErrorIs(t, errExisting, models.ErrShortyExists)
			},
		},
		{
			name: "redirect with ID",
			input: models.ShortyInput{
				ID:  "test",
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				_, err := shorty.Create(ctx, test.input)
				require.NoError(t, err)

				url, err := shorty.Redirect(ctx, models.VisitInput{
					ShortyID:  test.input.ID,
					Referer:   "https://example.com",
					UserIP:    "127.0.0.1",
					UserAgent: "",
				})
				require.NoError(t, err)

				sh, err := shorty.Get(ctx, test.input.ID)
				require.NoError(t, err)

				assert.Equal(t, url.String(), test.input.URL)
				assert.True(t, len(sh.Visits) == 1)
			},
		},
		{
			name: "redirect not found",
			input: models.ShortyInput{
				ID: "not_found",
			},
			run: func(test test, shorty *Shorty) {
				_, err := shorty.Redirect(ctx, models.VisitInput{
					ShortyID:  test.input.ID,
					Referer:   "https://example.com",
					UserIP:    "127.0.0.1",
					UserAgent: "",
				})

				require.ErrorIs(t, err, models.ErrShortyNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shorty := New(memory.NewShortyStorage())

			test.run(test, shorty)
		})
	}
}
