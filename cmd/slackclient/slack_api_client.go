package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	errSlackAPIHTTP     = errors.New("slack API returned an HTTP error")
	errSlackAPIRejected = errors.New("slack API rejected the message")
)

type slackAPIClient struct {
	token      string
	httpClient *http.Client
	apiURL     string
}

type slackPostMessageRequest struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type slackPostMessageResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func newSlackAPIClient(token string) *slackAPIClient {
	return &slackAPIClient{
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiURL:     "https://slack.com/api/chat.postMessage",
	}
}

func (c *slackAPIClient) PostMessage(ctx context.Context, channel, text string) error {
	body, err := json.Marshal(slackPostMessageRequest{Channel: channel, Text: text})
	if err != nil {
		return fmt.Errorf("marshal Slack message: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create Slack API request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send Slack API request: %w", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("Error closing Slack API response body: %v", err)
		}
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		responseBody, readErr := io.ReadAll(io.LimitReader(response.Body, 4096))
		if readErr != nil {
			return fmt.Errorf("%w %d (reading response: %w)", errSlackAPIHTTP, response.StatusCode, readErr)
		}
		return fmt.Errorf(
			"%w %d: %s",
			errSlackAPIHTTP,
			response.StatusCode,
			strings.TrimSpace(string(responseBody)),
		)
	}

	var result slackPostMessageResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode Slack API response: %w", err)
	}
	if !result.OK {
		return fmt.Errorf("%w: %s", errSlackAPIRejected, result.Error)
	}
	return nil
}
