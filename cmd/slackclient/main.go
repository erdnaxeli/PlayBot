// Package main runs the Slack Events API client for PlayBot.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/config"
)

const slackEventsPath = "/slack/events"

func main() {
	applicationConfig, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error while reading config file: %s", err)
	}
	if applicationConfig.Slack.SigningSecret == "" {
		log.Fatal("SLACK_SIGNING_SECRET must be configured")
	}
	if applicationConfig.Slack.BotToken == "" {
		log.Fatal("SLACK_BOT_TOKEN must be configured")
	}

	executor := rpc.NewPlaybotCliProtobufClient(
		fmt.Sprintf("http://%s", applicationConfig.ServerAddress),
		&http.Client{Timeout: 10 * time.Second},
	)
	handler := newSlackEventsHandler(
		applicationConfig.Slack.SigningSecret,
		executor,
		newSlackAPIClient(applicationConfig.Slack.BotToken),
	)
	mux := http.NewServeMux()
	mux.Handle(slackEventsPath, handler)

	log.Printf("Starting Slack Events API endpoint on %s%s", applicationConfig.Slack.ListenAddress, slackEventsPath)
	server := &http.Server{
		Addr:              applicationConfig.Slack.ListenAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
