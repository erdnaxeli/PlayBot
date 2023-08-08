package textbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		msg  string
		args []string
	}{
		{
			msg:  "",
			args: []string{},
		},
		{
			msg:  "   ",
			args: []string{},
		},
		{
			msg:  " test -f  --lol   ",
			args: []string{"test", "-f", "--lol"},
		},
		{
			msg:  "test",
			args: []string{"test"},
		},
		{
			msg:  "!",
			args: []string{"!"},
		},
		{
			msg:  "!  ",
			args: []string{"!"},
		},
		{
			msg:  "  ! ",
			args: []string{"!"},
		},
	}
	for _, test := range tests {
		t.Run(
			test.msg,
			func(t *testing.T) {
				args := parseArgs(test.msg)

				assert.Equal(t, test.args, args)
			},
		)
	}
}
