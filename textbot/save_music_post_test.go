package textbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTags(t *testing.T) {
	tests := []struct {
		msg  string
		tags []string
	}{
		{"", nil},
		{"  ", nil},
		{"t", nil},
		{"t  #t", []string{"t"}},
		{"  #t  ", []string{"t"}},
		{"#t u #v ", []string{"t", "v"}},
	}

	for _, test := range tests {
		t.Run(
			test.msg,
			func(t *testing.T) {
				tags := extractTagsFromMsg(test.msg)

				assert.Equal(t, test.tags, tags)
			},
		)
	}
}
