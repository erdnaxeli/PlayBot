package matrix

import (
	"golang.org/x/exp/slog"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func (c Client) SendText(text string, roomID string) error {
	err := c.sendText(id.RoomID(roomID), text)
	if err != nil {
		return ErrSendMessage
	}

	return nil
}

func (c Client) sendText(roomID id.RoomID, message string) error {
	err := c.sendMessage(roomID, message, "")
	return err
}

func (c Client) sendMessage(roomID id.RoomID, message string, formattedMsg string) error {
	if formattedMsg == "" {
		formattedMsg = message
	}

	_, err := c.client.SendMessageEvent(
		roomID,
		event.EventMessage,
		map[string]string{
			"body":           message,
			"format":         "org.matrix.custom.html",
			"formatted_body": formattedMsg,
			"msgtype":        "m.text",
		},
	)
	if err != nil {
		slog.Error("Error while sending message: %v", err)
		return err
	}

	return nil
}
