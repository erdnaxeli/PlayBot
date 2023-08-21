package mariadb

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/erdnaxeli/PlayBot/cmd/server/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAccountAssociationRepository(t *testing.T) server.UserNickAssociationRepository {
	r, err := New("test", "test", "localhost", "test")
	require.Nil(
		t,
		err,
		"A MariaDB server must be listening on localhost with a user 'test', a password 'test' and a database 'test' initialized with the test-db.sql file.",
	)
	return r
}

func TestGetUserFromNick_unknownNick(t *testing.T) {
	r := getAccountAssociationRepository(t)
	nick := gofakeit.Name()

	user, err := r.GetUserFromNick(nick)
	require.Nil(t, err)

	assert.Equal(t, user, "")
}

func TestGetUserFromNick_ok(t *testing.T) {
	user := gofakeit.Name()
	nick := gofakeit.Name()
	code := gofakeit.Password(true, true, true, true, true, 10)
	r, err := New("test", "test", "localhost", "test")
	require.Nil(
		t,
		err,
		"A MariaDB server must be listening on localhost with a user 'test', a password 'test' and a database 'test' initialized with the test-db.sql file.",
	)
	_, err = r.db.Exec(
		"insert into playbot_codes (user, nick, code) values (?, ?, ?)",
		user,
		nick,
		code,
	)
	require.Nil(t, err)

	result, err := r.GetUserFromNick(nick)
	require.Nil(t, err)

	assert.Equal(t, user, result)
}

func TestGetUserFromCode_unknownCode(t *testing.T) {
	r := getAccountAssociationRepository(t)
	code := gofakeit.BitcoinAddress()

	user, err := r.GetUserFromCode(code)
	require.Nil(t, err)

	assert.Equal(t, user, "")
}

func TestGetUserFromCode_ok(t *testing.T) {
	code := gofakeit.Password(true, true, true, true, true, 10)
	user := gofakeit.Name()
	r, err := New("test", "test", "localhost", "test")
	require.Nil(
		t,
		err,
		"A MariaDB server must be listening on localhost with a user 'test', a password 'test' and a database 'test' initialized with the test-db.sql file.",
	)
	_, err = r.db.Exec("insert into playbot_codes (user, code) values (?, ?)", user, code)
	require.Nil(t, err)

	result, err := r.GetUserFromCode(code)
	require.Nil(t, err)

	assert.Equal(t, user, result)
}

func TestSaveUserAssociation_unknownUser(t *testing.T) {
	r := getAccountAssociationRepository(t)
	user := gofakeit.Name()
	nick := gofakeit.Name()

	err := r.SaveAssociation(user, nick)

	assert.ErrorIs(t, err, UserNotFoundErr)
}

func TestSaveUserAssociation_existingUser_sameNick(t *testing.T) {
	user := gofakeit.Name()
	nick := gofakeit.Name()
	code := gofakeit.Password(true, true, true, true, true, 10)
	r, err := New("test", "test", "localhost", "test")
	require.Nil(
		t,
		err,
		"A MariaDB server must be listening on localhost with a user 'test', a password 'test' and a database 'test' initialized with the test-db.sql file.",
	)
	_, err = r.db.Exec(
		"insert into playbot_codes (user, nick, code) values (?, ?, ?)",
		user,
		nick,
		code,
	)
	require.Nil(t, err)

	err = r.SaveAssociation(user, nick)
	require.Nil(t, err)

	result, err := r.GetUserFromNick(nick)
	require.Nil(t, err)
	assert.Equal(t, user, result)
}

func TestSaveUserAssociation_existingUser_newNick(t *testing.T) {
	user := gofakeit.Name()
	nick := gofakeit.Name()
	code := gofakeit.Password(true, true, true, true, true, 10)
	r, err := New("test", "test", "localhost", "test")
	require.Nil(
		t,
		err,
		"A MariaDB server must be listening on localhost with a user 'test', a password 'test' and a database 'test' initialized with the test-db.sql file.",
	)
	_, err = r.db.Exec(
		"insert into playbot_codes (user, code) values (?, ?)",
		user,
		code,
	)
	require.Nil(t, err)

	err = r.SaveAssociation(user, nick)
	require.Nil(t, err)

	result, err := r.GetUserFromNick(nick)
	require.Nil(t, err)
	assert.Equal(t, user, result)
}
