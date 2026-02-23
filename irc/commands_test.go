package irc_test

import (
	"fmt"
	"testing"
	"testing/synctest"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrivsg(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		conn, _, server := connect(t)

		to := gofakeit.FirstName()
		msg := gofakeit.Sentence()
		go func() {
			err := conn.Privmsg(to, msg)
			require.Nil(t, err)
		}()

		expected := fmt.Sprintf("PRIVMSG %s :%s\r\n", to, msg)
		buffer := make([]byte, len(expected))
		n, err := server.Read(buffer)
		require.Equal(t, len(buffer), n)
		require.Nil(t, err)

		assert.Equal(t, expected, string(buffer))
	})
}
