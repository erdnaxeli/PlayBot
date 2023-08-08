package iso8601_test

import (
	"testing"
	"time"

	"github.com/erdnaxeli/PlayBot/iso8601"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input  string
		output time.Duration
	}{
		{"PT", 0 * time.Second},
		{"PT0H", 0 * time.Second},
		{"PT0M", 0 * time.Second},
		{"PT0S", 0 * time.Second},
		{"PT0M0S", 0 * time.Second},
		{"PT0H0S", 0 * time.Second},
		{"PT0H0M", 0 * time.Second},
		{"PT1S", 1 * time.Second},
		{"PT2M", 2 * time.Minute},
		{"PT3H", 3 * time.Hour},
		{"PT4M5S", 4*time.Minute + 5*time.Second},
		{"PT6H7S", 6*time.Hour + 7*time.Second},
		{"PT8H9M", 8*time.Hour + 9*time.Minute},
		{"P", 0 * time.Second},
		{"P0H", 0 * time.Second},
		{"P0M", 0 * time.Second},
		{"P0S", 0 * time.Second},
		{"P0M0S", 0 * time.Second},
		{"P0H0S", 0 * time.Second},
		{"P0H0M", 0 * time.Second},
		{"P1S", 1 * time.Second},
		{"P2M", 2 * time.Minute},
		{"P3H", 3 * time.Hour},
		{"P4M5S", 4*time.Minute + 5*time.Second},
		{"P6H7S", 6*time.Hour + 7*time.Second},
		{"P8H9M", 8*time.Hour + 9*time.Minute},
	}

	for _, test := range tests {
		t.Run(
			test.input,
			func(t *testing.T) {
				result, err := iso8601.ParseDuration(test.input)
				require.Nil(t, err)

				assert.Equal(t, test.output, result)
			},
		)
	}
}
