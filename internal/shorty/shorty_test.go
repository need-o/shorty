package shorty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	type test struct {
		in  uint32
		out string
	}

	tests := []test{
		{
			in:  123456789,
			out: "HUawi",
		},
		{
			in:  987654321,
			out: "zag0eb",
		},
	}

	for _, test := range tests {
		out1 := NewID(test.in)
		out2 := NewID(test.in)

		assert.Equal(t, out1, test.out)
		assert.Equal(t, out1, out2)
	}
}
