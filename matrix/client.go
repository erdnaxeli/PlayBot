package matrix

import (
	"fmt"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const MSG_MAX_SIZE = 10

type Client struct {
	startTime time.Time

	client *mautrix.Client
	config Config
}

type Config struct {
	HomeserverURL string
	UserID        string
	AccessToken   string
}

func NewClient(config Config) (Client, error) {
	matrixClient, err := mautrix.NewClient(
		config.HomeserverURL, id.UserID(config.UserID), config.AccessToken,
	)
	if err != nil {
		//nolint:errorlint
		return Client{}, fmt.Errorf("Error while creating client: %v", err)
	}

	client := Client{
		startTime: time.Now(),
		config:    config,
		client:    matrixClient,
	}
	syncer := matrixClient.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, client.onEventMessage)

	return client, nil
}

func (c Client) onEventMessage(source mautrix.EventSource, evt *event.Event) {
	evtTime := time.Unix(evt.Timestamp/1000, 0)
	if evtTime.Before(c.startTime) {
		return
	}

	content := evt.Content.AsMessage()
	fmt.Print(content)
}
