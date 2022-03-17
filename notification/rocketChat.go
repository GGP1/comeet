package notification

import (
	"net/url"

	"github.com/GGP1/comeet/event"
	"github.com/pkg/errors"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/rest"
)

type rocketChatNotifier struct {
	client  *rest.Client
	roomIDs []string
}

// NewRocketChatNotifier returns a notifier that sends rocket chat messages.
func NewRocketChatNotifier(config RocketChatConfig) (Notifier, error) {
	serverURL, err := url.Parse(config.ServerURL)
	if err != nil {
		return nil, err
	}

	client := rest.NewClient(serverURL, false)
	creds := &models.UserCredentials{
		ID:    config.UserID,
		Token: config.Token,
	}
	if err := client.Login(creds); err != nil {
		return nil, err
	}

	notifier := &rocketChatNotifier{
		client:  client,
		roomIDs: config.RoomIDs,
	}
	return notifier, nil
}

func (r *rocketChatNotifier) Notify(event *event.Event) error {
	message := &models.PostMessage{
		Text:      event.Message(),
		ParseUrls: true,
		Alias:     "Comeet",
	}

	for _, roomID := range r.roomIDs {
		message.RoomID = roomID

		if _, err := r.client.PostMessage(message); err != nil {
			return errors.Wrapf(err, "rocketChat: failed sending message to room %q", roomID)
		}
	}

	return nil
}
