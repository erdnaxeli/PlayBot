package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/erdnaxeli/PlayBot/cmd/server/rpc"
)

const (
	maxSlackRequestBody   = 1_048_576
	slackSignatureVersion = "v0"
	slackSignatureAge     = 5 * time.Minute
	slackEventRetryWindow = 15 * time.Minute
	slackEventQueueSize   = 100
)

var (
	errMissingSlackCredentials = errors.New("missing Slack signature credentials")
	errStaleSlackRequest       = errors.New("slack request timestamp is outside the allowed window")
	errInvalidSlackSignature   = errors.New("slack request signature does not match")
	errNilRPCResult            = errors.New("RPC returned a nil result")
	errNilRPCReply             = errors.New("RPC returned a nil reply")
	errMissingReplyDestination = errors.New("RPC reply has no destination")
)

type playbotExecutor interface {
	Execute(context.Context, *rpc.TextMessage) (*rpc.Result, error)
}

type slackMessagePoster interface {
	PostMessage(context.Context, string, string) error
}

type slackEventsHandler struct {
	signingSecret string
	executor      playbotExecutor
	poster        slackMessagePoster
	now           func() time.Time
	queue         chan slackEnvelope
	seenMu        sync.Mutex
	seen          map[string]time.Time
}

type slackEnvelope struct {
	Type      string     `json:"type"`
	Challenge string     `json:"challenge"`
	EventID   string     `json:"event_id"`
	Event     slackEvent `json:"event"`
}

type slackEvent struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	BotID   string `json:"bot_id"`
}

func newSlackEventsHandler(
	signingSecret string,
	executor playbotExecutor,
	poster slackMessagePoster,
) *slackEventsHandler {
	handler := &slackEventsHandler{
		signingSecret: signingSecret,
		executor:      executor,
		poster:        poster,
		now:           time.Now,
		queue:         make(chan slackEnvelope, slackEventQueueSize),
		seen:          make(map[string]time.Time),
	}
	go handler.runEvents()
	return handler
}

func (h *slackEventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxSlackRequestBody))
	if err != nil {
		http.Error(w, "request body is invalid or too large", http.StatusBadRequest)
		return
	}

	if err := verifySlackRequest(
		h.signingSecret,
		r.Header.Get("X-Slack-Request-Timestamp"),
		r.Header.Get("X-Slack-Signature"),
		body,
		h.now(),
	); err != nil {
		http.Error(w, "invalid Slack signature", http.StatusUnauthorized)
		return
	}

	var envelope slackEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		http.Error(w, "invalid Slack event payload", http.StatusBadRequest)
		return
	}

	switch envelope.Type {
	case "url_verification":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(struct {
			Challenge string `json:"challenge"`
		}{Challenge: envelope.Challenge}); err != nil {
			http.Error(w, "unable to write Slack URL verification response", http.StatusInternalServerError)
		}
	case "event_callback":
		if h.enqueue(envelope) {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Slack event queue is full", http.StatusServiceUnavailable)
		}
	default:
		http.Error(w, "unsupported Slack request type", http.StatusBadRequest)
	}
}

func (h *slackEventsHandler) enqueue(envelope slackEnvelope) bool {
	h.seenMu.Lock()
	defer h.seenMu.Unlock()

	now := h.now()
	for eventID, seenAt := range h.seen {
		if now.Sub(seenAt) >= slackEventRetryWindow {
			delete(h.seen, eventID)
		}
	}

	if envelope.EventID != "" {
		if _, exists := h.seen[envelope.EventID]; exists {
			return true
		}
	}

	select {
	case h.queue <- envelope:
		if envelope.EventID != "" {
			h.seen[envelope.EventID] = now
		}
		return true
	default:
		return false
	}
}

func (h *slackEventsHandler) runEvents() {
	for envelope := range h.queue {
		if err := h.handleEvent(context.Background(), envelope); err != nil {
			log.Printf("Error processing Slack event %q: %v", envelope.EventID, err)
		}
	}
}

func (h *slackEventsHandler) handleEvent(ctx context.Context, envelope slackEnvelope) error {
	event := envelope.Event
	if !isUserMessage(event) {
		return nil
	}

	result, err := h.executor.Execute(ctx, &rpc.TextMessage{
		ChannelName: event.Channel,
		PersonName:  event.User,
		Msg:         event.Text,
	})
	if err != nil {
		return fmt.Errorf("execute message through RPC: %w", err)
	}
	if result == nil {
		return errNilRPCResult
	}

	return h.postReplies(ctx, result)
}

func (h *slackEventsHandler) postReplies(ctx context.Context, result *rpc.Result) error {
	for _, message := range result.Msg {
		if message == nil {
			return errNilRPCReply
		}
		if message.To == "" {
			return errMissingReplyDestination
		}
		if err := h.poster.PostMessage(ctx, message.To, stripIRCFormatting(message.Msg)); err != nil {
			return fmt.Errorf("post Slack reply to %q: %w", message.To, err)
		}
	}

	return nil
}

func verifySlackRequest(secret, timestamp, signature string, body []byte, now time.Time) error {
	if secret == "" || timestamp == "" || signature == "" {
		return errMissingSlackCredentials
	}

	unixTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("parse Slack request timestamp: %w", err)
	}
	requestTime := time.Unix(unixTimestamp, 0)
	if age := now.Sub(requestTime); age > slackSignatureAge || age < -slackSignatureAge {
		return errStaleSlackRequest
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = fmt.Fprintf(mac, "%s:%s:", slackSignatureVersion, timestamp)
	_, _ = mac.Write(body)
	expected := slackSignatureVersion + "=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return errInvalidSlackSignature
	}

	return nil
}

var ircColorCode = regexp.MustCompile(`\x03(?:[0-9]{1,2}(?:,[0-9]{1,2})?)?`)

func stripIRCFormatting(message string) string {
	message = ircColorCode.ReplaceAllString(message, "")
	return strings.Map(func(char rune) rune {
		switch char {
		case '\x02', '\x0f', '\x11', '\x16', '\x1d', '\x1f':
			return -1
		default:
			return char
		}
	}, message)
}

func isUserMessage(event slackEvent) bool {
	return event.Type == "message" &&
		event.Subtype == "" &&
		event.BotID == "" &&
		event.User != "" &&
		event.Channel != "" &&
		strings.TrimSpace(event.Text) != ""
}
