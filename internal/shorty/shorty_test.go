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
		input models.ShorteningInput
		run   func(test, *Shorty)
	}

	tests := []test{
		{
			name: "create with URL",
			input: models.ShorteningInput{
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				sh, err := shorty.Create(context.Background(), test.input)

				require.NoError(t, err)
				assert.NotEmpty(t, sh.ID)
				assert.Equal(t, sh.URL, test.input.URL)
				assert.NotZero(t, sh.CreatedAt)
				assert.NotZero(t, sh.CreatedAt)
			},
		},
		{
			name: "create with ID and URL",
			input: models.ShorteningInput{
				ID:  "test",
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				sh, err := shorty.Create(context.Background(), test.input)

				require.NoError(t, err)
				assert.Equal(t, sh.ID, test.input.ID)
				assert.Equal(t, sh.URL, test.input.URL)
				assert.NotZero(t, sh.CreatedAt)
				assert.NotZero(t, sh.CreatedAt)
			},
		},
		{
			name: "create with existing ID",
			input: models.ShorteningInput{
				ID:  "test",
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				_, err := shorty.Create(context.Background(), test.input)
				_, errExisting := shorty.Create(context.Background(), test.input)

				require.NoError(t, err)
				assert.ErrorIs(t, errExisting, models.ErrShorteningExists)
			},
		},
		{
			name: "redirect with ID",
			input: models.ShorteningInput{
				ID:  "test",
				URL: "https://example.com",
			},
			run: func(test test, shorty *Shorty) {
				_, err := shorty.Create(context.Background(), test.input)
				require.NoError(t, err)

				rdr, err := shorty.Redirect(context.Background(), test.input.ID)
				require.NoError(t, err)

				sh2, err := shorty.Get(context.Background(), test.input.ID)
				require.NoError(t, err)

				assert.Equal(t, rdr.URL, test.input.URL)
				assert.True(t, sh2.Visits == 1)
			},
		},
		{
			name: "redirect not found",
			input: models.ShorteningInput{
				ID: "not_found",
			},
			run: func(test test, shorty *Shorty) {
				_, err := shorty.Redirect(context.Background(), test.input.ID)

				require.ErrorIs(t, err, models.ErrShorteningNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shorty := New(memory.NewShorteningStorage())

			test.run(test, shorty)
		})
	}
}
