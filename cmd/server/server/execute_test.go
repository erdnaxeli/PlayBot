package server_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/cmd/server/server"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestTextBot struct {
	mock.Mock
}

func (t *TestTextBot) Execute(
	channelName string,
	personName string,
	msg string,
	user string,
) (textbot.Result, bool, error) {
	args := t.Called(channelName, personName, msg, user)
	return args.Get(0).(textbot.Result), args.Bool(1), args.Error(2)
}

type TestUserNickAssociationRepository struct {
	mock.Mock
}

func (t *TestUserNickAssociationRepository) GetUserFromNick(nick string) (string, error) {
	args := t.Called(nick)
	return args.String(0), args.Error(1)
}

func (t *TestUserNickAssociationRepository) GetUserFromCode(code string) (string, error) {
	args := t.Called(code)
	return args.String(0), args.Error(1)
}

func (t *TestUserNickAssociationRepository) SaveAssociation(user string, nick string) error {
	args := t.Called(user, nick)
	return args.Error(0)
}

type TestMusicRecordPrinter struct {
	mock.Mock
}

func (t *TestMusicRecordPrinter) PrintResult(result textbot.Result) string {
	args := t.Called(result)
	return args.String(0)
}

type TestStatsPrinter struct {
	mock.Mock
}

func (t *TestStatsPrinter) PrintStatistics(stats playbot.MusicRecordStatistics) string {
	args := t.Called(stats)
	return args.String(0)
}

func getResult() textbot.Result {
	var record types.MusicRecord
	_ = gofakeit.Struct(&record)
	record.Duration = time.Second * time.Duration(rand.Int64N(300))

	tags := make([]string, rand.IntN(10))
	gofakeit.Slice(&tags)

	return textbot.Result{
		ID:          gofakeit.Int64(),
		MusicRecord: record,
		Tags:        tags,
		IsNew:       gofakeit.Bool(),
		Count:       gofakeit.Int64(),
	}
}

func TestExecute_noRecord_noCmd(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		textbot.Result{}, false, nil,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(t, &rpc.Result{}, result)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test the behavior when the record is not the result of a command, which means it is
// the user that just posted it.
// It should print the record **without** the URL.
func TestExecute_record_noCmd(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	var record types.MusicRecord
	_ = gofakeit.Struct(&record)
	record.Duration = time.Second * time.Duration(rand.Int64N(300))
	tags := make([]string, rand.IntN(10))
	gofakeit.Slice(&tags)
	execResult := textbot.Result{
		ID:          gofakeit.Int64(),
		MusicRecord: record,
		Tags:        tags,
		IsNew:       gofakeit.Bool(),
		Count:       gofakeit.Int64(),
	}
	printResult := gofakeit.Phrase()

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, false, nil,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	// The record does not comes from a command, the URL must be removed before calling
	// Print().
	printArgs := execResult
	printArgs.URL = ""
	mrp := &TestMusicRecordPrinter{}
	mrp.On("PrintResult", printArgs).Return(printResult)

	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: printResult,
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test the behavior when the record is the result of a command.
// It should print the record **with** the URL.
func TestExecute_record_cmd(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := getResult()
	printResult := gofakeit.Phrase()

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, nil,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	mrp.On("PrintResult", execResult).Return(printResult)

	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: printResult,
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test the behavior when no record is returned, but statistics are.
// It should print them.
func TestExecute_statistics(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := textbot.Result{}
	_ = gofakeit.Struct(&execResult.Statistics)

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, nil,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}

	statsPrintResult := gofakeit.Phrase()
	sp := &TestStatsPrinter{}
	sp.On("PrintStatistics", execResult.Statistics).Return(statsPrintResult)

	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: statsPrintResult,
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that when the error NoRecordFoundError is returned with a count of 0, the
// "no result found" message is returned.
func TestExecute_NoRecordFoundError_noCount(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := textbot.Result{
		Count: 0,
	}

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, playbot.ErrNoRecordFound,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Je n'ai rien dans ce registre.",
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that when the error NoRecordFoundError is returned with a count > 0, the "end
// of search" message is returned.
func TestExecute_NoRecordFoundError_count(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := textbot.Result{
		Count: 10,
	}

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, playbot.ErrNoRecordFound,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Tu tournes en rond, Jack !",
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that when an InvalidOffsetError is returned, the corresponding error message
// is returned.
func TestExecute_InvalidOffsetError(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := textbot.Result{}

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, playbot.ErrInvalidOffset,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Offset invalide.",
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that when an OffsetToBigError is returned, the corresponding error message
// is returned.
func TestExecute_OffsetToBigError(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := textbot.Result{}

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, textbot.ErrOffsetToBig,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "T'as compté tout ça sans te tromper, srsly ?",
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that when an InvalidUsageError is returned, the corresponding error message
// is returned.
func TestExecute_InvalidUsageError(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := textbot.Result{}

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, textbot.ErrInvalidUsage,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Len(t, result.Msg, 1)
	assert.Equal(t, channelName, result.Msg[0].To)
	assert.Contains(
		t,
		[]string{
			"Ahahahah ! 23 à 0 !",
			"C'est la piquette, Jack !",
			"Tu sais pas jouer, Jack !",
			"T'es mauvais, Jack !",
		},
		result.Msg[0].Msg,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that when an AuthenticationRequired error is returned, a message with the link
// to auth is sent to the user.
func TestExecute_AuthenticationRequired_noResult(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := textbot.Result{}

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, textbot.ErrAuthenticationRequired,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Ce nick n'est associé à aucun login arise. Va sur http://nightiies.iiens.net/links/fav pour obtenir ton code personel.",
					To:  personName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that when an AuthenticationRequired error is returned, a message with the link
// to auth is sent to the user, and the returned result is sent to the channel.
func TestExecute_AuthenticationRequired_result(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	channelName := gofakeit.DomainName()
	personName := gofakeit.FirstName()
	msg := gofakeit.Phrase()
	user := gofakeit.Name()
	execResult := getResult()

	b := &TestTextBot{}
	b.On("Execute", channelName, personName, msg, user).Return(
		execResult, true, textbot.ErrAuthenticationRequired,
	)

	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromNick", personName).Return(user, nil)

	printResult := gofakeit.Phrase()
	mrp := &TestMusicRecordPrinter{}
	mrp.On("PrintResult", execResult).Return(printResult)

	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: channelName,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Ce nick n'est associé à aucun login arise. Va sur http://nightiies.iiens.net/links/fav pour obtenir ton code personel.",
					To:  personName,
				},
				{
					Msg: printResult,
					To:  channelName,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that a message sent to PlayTest and starting with "PB" triggers the auth
// process.
func TestExecute_authMessage(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	personName := gofakeit.FirstName()
	msg := "PB" + gofakeit.Phrase()

	b := &TestTextBot{}
	r := &TestUserNickAssociationRepository{}
	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Vérification en cours…",
					To:  personName,
				},
				{
					Msg: "info " + personName,
					To:  "NickServ",
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test an message from NickServ telling that the nick is not registered returns an
// error message.
func TestExecute_nickserv_notRegistered(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	personName := "NickServ"
	nick := gofakeit.FirstName()
	msg := fmt.Sprintf("Le pseudo %s%s%s n'est pas enregistré.", string(server.STX), nick, string(server.STX))

	b := &TestTextBot{}
	r := &TestUserNickAssociationRepository{}
	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// we need to send an auth message to start the process
	_, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  nick,
			Msg:         "PB" + gofakeit.Phrase(),
		},
	)
	require.Nil(t, err)

	// we also starts the process for another user to ensure it will select the correct
	// one
	_, err = server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  gofakeit.FirstName(),
			Msg:         "PB" + gofakeit.Phrase(),
		},
	)
	require.Nil(t, err)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Il faut que ton pseudo soit enregistré auprès de NickServ pour pouvoir t'authentifier.",
					To:  nick,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test when receiving a message from NickServ telling that the nick is registered it
// saves the user and nick association.
func TestExecute_nickserv_registered_codeOk(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	personName := "NickServ"
	nick := gofakeit.FirstName()
	msg := fmt.Sprintf("%s est actuellement connecté.", nick)
	user := gofakeit.Username()
	code := "PB" + gofakeit.Phrase()

	b := &TestTextBot{}
	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromCode", code).Return(user, nil)
	r.On("SaveAssociation", user, nick).Return(nil)
	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// we need to send an auth message to start the process
	_, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  nick,
			Msg:         code,
		},
	)
	require.Nil(t, err)

	// we also starts the process for another user to ensure it will select the correct
	// one
	_, err = server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  gofakeit.FirstName(),
			Msg:         "PB" + gofakeit.Phrase(),
		},
	)
	require.Nil(t, err)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Association effectuée. Utilise la commande !fav pour enregistrer un lien dans tes favoris.",
					To:  nick,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test when receiving a message from NickServ telling that the nick is registered, but
// it does found the user associated to the code, it returns an error message.
func TestExecute_nickserv_registered_codeUnknown(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	personName := "NickServ"
	nick := gofakeit.FirstName()
	msg := fmt.Sprintf("%s est actuellement connecté.", nick)
	code := "PB" + gofakeit.Phrase()

	b := &TestTextBot{}
	r := &TestUserNickAssociationRepository{}
	r.On("GetUserFromCode", code).Return("", nil)
	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// we need to send an auth message to start the process
	_, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  nick,
			Msg:         code,
		},
	)
	require.Nil(t, err)

	// we also starts the process for another user to ensure it will select the correct
	// one
	_, err = server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  gofakeit.FirstName(),
			Msg:         "PB" + gofakeit.Phrase(),
		},
	)
	require.Nil(t, err)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(
		t,
		&rpc.Result{
			Msg: []*rpc.IrcMessage{
				{
					Msg: "Code inconnu. Va sur http://nightiies.iiens.net/links/fav pour obtenir ton code personel.",
					To:  nick,
				},
			},
		},
		result,
	)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}

// Test that a message from NickServ for a user who didn't start the auth process is
// ignored.
func TestExecute_nickserv_unknownMessage(t *testing.T) {
	// setup

	botNick := gofakeit.Phrase()
	personName := "NickServ"
	msg := fmt.Sprintf("Le pseudo %s%s%s n'est pas enregistré.", string(server.STX), gofakeit.FirstName(), string(server.STX))

	b := &TestTextBot{}
	r := &TestUserNickAssociationRepository{}
	mrp := &TestMusicRecordPrinter{}
	sp := &TestStatsPrinter{}
	server := server.NewServer(botNick, b, r, mrp, sp)

	// we need to send an auth message to start the process
	_, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  gofakeit.FirstName(),
			Msg:         "PB" + gofakeit.Phrase(),
		},
	)
	require.Nil(t, err)

	// test

	result, err := server.Execute(
		context.Background(),
		&rpc.TextMessage{
			ChannelName: botNick,
			PersonName:  personName,
			Msg:         msg,
		},
	)
	require.Nil(t, err, nil)

	// assertions

	assert.Equal(t, &rpc.Result{}, result)
	b.AssertExpectations(t)
	r.AssertExpectations(t)
	mrp.AssertExpectations(t)
	sp.AssertExpectations(t)
}
