package irc_test

import (
	"net"
	"testing"
	"testing/synctest"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/erdnaxeli/PlayBot/irc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dropEvent(t *testing.T, conn net.Conn) {
	dropEvents(t, 1, conn)
}

func dropEvents(t *testing.T, count int, conn net.Conn) {
	buffer := make([]byte, 1000)
	for range count {
		_, err := conn.Read(buffer)
		require.Nil(t, err)
	}
}

func disconnect(t *testing.T, irc *irc.Conn, conn net.Conn) {
	go dropEvent(t, conn)
	irc.Disconnect()
}

// connect create a new irc.Conn instance and return the irc object, the config, and the server socket.
func connect(t *testing.T) (*irc.Conn, irc.Config, net.Conn) {
	server, client := net.Pipe()
	nick := gofakeit.FirstName()

	// Drop the connection events sent
	// We need to read before creating the Conn object,
	// else it would block while writing to the Pipe.
	go dropEvents(t, 2, server)

	config := irc.Config{
		Host: "localhost",
		Port: 42,
		Nick: nick,

		SocketFactory: func(_ irc.Config) (net.Conn, error) {
			return client, nil
		},
	}
	irc, err := irc.New(config)
	require.Nil(t, err)

	return irc, config, server
}

func dispatch(t *testing.T, irc *irc.Conn) {
	err := irc.Dispatch()
	require.Nil(t, err)
}

func write(t *testing.T, conn net.Conn, msg string) {
	n, err := conn.Write([]byte(msg))
	require.Nil(t, err)
	require.Equal(t, len(msg), n)
}

func TestDispatchPing(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		conn, _, server := connect(t)
		go dispatch(t, conn)

		// Send a PING event
		write(t, server, "PING test\r\n")
		// Wait for the event to be handled
		synctest.Wait()

		expected := "PONG :test\r\n"
		buffer := make([]byte, len(expected))
		_, err := server.Read(buffer)
		require.Nil(t, err)

		assert.Equal(t, expected, string(buffer))

		disconnect(t, conn, server)
	})
}

func TestDispatchPrivMsg(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		conn, _, server := connect(t)
		go dispatch(t, conn)

		conn.OnPrivMsg(func(_ *irc.Conn, msg irc.Message) error {
			assert.Equal(
				t,
				irc.Message{
					Prefix:     "nick!username@host",
					Command:    "PRIVMSG",
					Parameters: []string{"#test", "this is a test"},
				},
				msg,
			)
			return nil
		})

		// Send a PRIVMSG event
		write(t, server, ":nick!username@host PRIVMSG #test :this is a test\r\n")

		disconnect(t, conn, server)
	})
}
