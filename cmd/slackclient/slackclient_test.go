package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/stretchr/testify/require"
)

var (
	errTestRPC   = errors.New("rpc failed")
	errTestSlack = errors.New("slack failed")
)

type testExecutor struct {
	message *rpc.TextMessage
	result  *rpc.Result
	err     error
	calls   int
	seen    chan *rpc.TextMessage
}

func (e *testExecutor) Execute(_ context.Context, message *rpc.TextMessage) (*rpc.Result, error) {
	e.calls++
	e.message = message
	if e.seen != nil {
		e.seen <- message
	}
	return e.result, e.err
}

type postedMessage struct {
	channel string
	text    string
}

type testPoster struct {
	channel string
	text    string
	err     error
	calls   int
	posted  chan postedMessage
}

func (p *testPoster) PostMessage(_ context.Context, channel, text string) error {
	p.calls++
	p.channel = channel
	p.text = text
	if p.posted != nil {
		p.posted <- postedMessage{channel: channel, text: text}
	}
	return p.err
}

func newTestSlackHandler(secret string, executor *testExecutor, poster *testPoster, now time.Time) *slackEventsHandler {
	return &slackEventsHandler{
		signingSecret: secret,
		executor:      executor,
		poster:        poster,
		now:           func() time.Time { return now },
		queue:         make(chan slackEnvelope, 10),
		seen:          make(map[string]time.Time),
	}
}

func signedSlackRequest(t *testing.T, handler http.Handler, secret string, now time.Time, payload string) *httptest.ResponseRecorder {
	t.Helper()
	timestamp := fmt.Sprintf("%d", now.Unix())
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := fmt.Fprintf(mac, "v0:%s:%s", timestamp, payload)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, slackEventsPath, strings.NewReader(payload))
	request.Header.Set("X-Slack-Request-Timestamp", timestamp)
	request.Header.Set("X-Slack-Signature", "v0="+hex.EncodeToString(mac.Sum(nil)))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestSlackEventsHandlerURLVerification(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	handler := newTestSlackHandler("secret", nil, nil, now)

	response := signedSlackRequest(t, handler, "secret", now,
		`{"type":"url_verification","challenge":"challenge-token"}`)

	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `{"challenge":"challenge-token"}`, response.Body.String())
}

func TestSlackEventsHandlerForwardsMessageAndPostsReply(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	executor := &testExecutor{result: &rpc.Result{Msg: []*rpc.IrcMessage{
		{To: "C123", Msg: "\x0308[42] \x0303Artist\x0f"},
	}}}
	poster := &testPoster{}
	handler := newTestSlackHandler("secret", executor, poster, now)

	response := signedSlackRequest(t, handler, "secret", now,
		`{"type":"event_callback","event_id":"Ev1","event":{"type":"message","user":"U123","channel":"C123","text":"!get"}}`)

	require.Equal(t, http.StatusOK, response.Code)
	envelope := <-handler.queue
	require.NoError(t, handler.handleEvent(context.Background(), envelope))
	require.Equal(t, 1, executor.calls)
	require.Equal(t, &rpc.TextMessage{ChannelName: "C123", PersonName: "U123", Msg: "!get"}, executor.message)
	require.Equal(t, 1, poster.calls)
	require.Equal(t, "C123", poster.channel)
	require.Equal(t, "[42] Artist", poster.text)
}

func TestSlackEventsHandlerProcessesQueuedEvents(t *testing.T) {
	now := time.Now()
	executor := &testExecutor{
		result: &rpc.Result{Msg: []*rpc.IrcMessage{{To: "C123", Msg: "reply"}}},
		seen:   make(chan *rpc.TextMessage, 1),
	}
	poster := &testPoster{posted: make(chan postedMessage, 1)}
	handler := newSlackEventsHandler("secret", executor, poster)
	handler.now = func() time.Time { return now }
	defer close(handler.queue)

	response := signedSlackRequest(t, handler, "secret", now,
		`{"type":"event_callback","event_id":"Ev1","event":{"type":"message","user":"U123","channel":"C123","text":"hello"}}`)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, &rpc.TextMessage{ChannelName: "C123", PersonName: "U123", Msg: "hello"}, <-executor.seen)
	require.Equal(t, postedMessage{channel: "C123", text: "reply"}, <-poster.posted)
}

func TestSlackEventsHandlerIgnoresNonUserMessages(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	for _, event := range []string{
		`{"type":"message","user":"U123","channel":"C123","text":"hello","bot_id":"B123"}`,
		`{"type":"message","user":"U123","channel":"C123","text":"hello","subtype":"message_changed"}`,
		`{"type":"reaction_added","user":"U123","channel":"C123","text":"hello"}`,
		`{"type":"message","user":"U123","channel":"C123","text":"   "}`,
	} {
		t.Run(event, func(t *testing.T) {
			executor := &testExecutor{result: &rpc.Result{}}
			poster := &testPoster{}
			handler := newTestSlackHandler("secret", executor, poster, now)

			event = strings.TrimSuffix(strings.TrimPrefix(event, "{"), "}")
			response := signedSlackRequest(t, handler, "secret", now,
				fmt.Sprintf(`{"type":"event_callback","event":{%s}}`, event))

			require.Equal(t, http.StatusOK, response.Code)
			envelope := <-handler.queue
			require.NoError(t, handler.handleEvent(context.Background(), envelope))
			require.Zero(t, executor.calls)
			require.Zero(t, poster.calls)
		})
	}
}

func TestSlackEventsHandlerRejectsInvalidAndStaleSignatures(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	handler := newTestSlackHandler("secret", nil, nil, now)

	request := httptest.NewRequest(http.MethodPost, slackEventsPath, strings.NewReader(`{}`))
	request.Header.Set("X-Slack-Request-Timestamp", fmt.Sprintf("%d", now.Unix()))
	request.Header.Set("X-Slack-Signature", "v0=invalid")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	require.Equal(t, http.StatusUnauthorized, response.Code)

	response = signedSlackRequest(t, handler, "secret", now.Add(-slackSignatureAge-time.Second), `{}`)
	require.Equal(t, http.StatusUnauthorized, response.Code)
}

func TestSlackEventsHandlerRejectsInvalidMethodAndPayload(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	handler := newTestSlackHandler("secret", nil, nil, now)

	request := httptest.NewRequest(http.MethodGet, slackEventsPath, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	require.Equal(t, http.StatusMethodNotAllowed, response.Code)

	response = signedSlackRequest(t, handler, "secret", now, `{`)
	require.Equal(t, http.StatusBadRequest, response.Code)
}

func TestSlackEventsHandlerReturnsErrorsFromRPCAndPoster(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	for name, testCase := range map[string]struct {
		executor *testExecutor
		poster   *testPoster
	}{
		"rpc": {
			executor: &testExecutor{err: errTestRPC},
			poster:   &testPoster{},
		},
		"poster": {
			executor: &testExecutor{result: &rpc.Result{Msg: []*rpc.IrcMessage{{To: "C123", Msg: "reply"}}}},
			poster:   &testPoster{err: errTestSlack},
		},
	} {
		t.Run(name, func(t *testing.T) {
			handler := newTestSlackHandler("secret", testCase.executor, testCase.poster, now)
			response := signedSlackRequest(t, handler, "secret", now,
				`{"type":"event_callback","event_id":"Ev1","event":{"type":"message","user":"U123","channel":"C123","text":"hello"}}`)
			require.Equal(t, http.StatusOK, response.Code)
			envelope := <-handler.queue
			require.Error(t, handler.handleEvent(context.Background(), envelope))
		})
	}
}

func TestSlackEventsHandlerDeduplicatesQueuedEvents(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	handler := newTestSlackHandler("secret", &testExecutor{}, &testPoster{}, now)
	payload := `{"type":"event_callback","event_id":"Ev1","event":{"type":"message","user":"U123","channel":"C123","text":"hello"}}`

	firstResponse := signedSlackRequest(t, handler, "secret", now, payload)
	secondResponse := signedSlackRequest(t, handler, "secret", now, payload)

	require.Equal(t, http.StatusOK, firstResponse.Code)
	require.Equal(t, http.StatusOK, secondResponse.Code)
	require.Len(t, handler.queue, 1)
}

func TestSlackEventsHandlerRejectsFullQueue(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	handler := newTestSlackHandler("secret", &testExecutor{}, &testPoster{}, now)
	handler.queue = make(chan slackEnvelope, 1)
	handler.queue <- slackEnvelope{}

	payload := `{"type":"event_callback","event_id":"Ev1","event":{"type":"message","user":"U123","channel":"C123","text":"hello"}}`
	response := signedSlackRequest(t, handler, "secret", now, payload)

	require.Equal(t, http.StatusServiceUnavailable, response.Code)
	require.Empty(t, handler.seen)
}

func TestSlackAPIClientPostsMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		require.Equal(t, "application/json; charset=utf-8", r.Header.Get("Content-Type"))
		var request slackPostMessageRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, slackPostMessageRequest{Channel: "C123", Text: "hello"}, request)
		_, err := w.Write([]byte(`{"ok":true}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := newSlackAPIClient("token")
	client.apiURL = server.URL
	require.NoError(t, client.PostMessage(context.Background(), "C123", "hello"))
}

func TestSlackAPIClientReportsAPIRejection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"ok":false,"error":"channel_not_found"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := newSlackAPIClient("token")
	client.apiURL = server.URL
	require.ErrorContains(t, client.PostMessage(context.Background(), "C123", "hello"), "channel_not_found")
}

func TestStripIRCFormatting(t *testing.T) {
	require.Equal(t, "[42] Artist", stripIRCFormatting("\x0308[42] \x0303Artist\x0f"))
}
