package notification

import (
	"github.com/GGP1/comeet/event"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type slackNotifier struct {
	client     *slack.Client
	channelIDs []string
}

// NewSlackNotifier returns a notifier that sends messages thorugh slack channels.
func NewSlackNotifier(config SlackConfig) (Notifier, error) {
	notifier := &slackNotifier{
		client:     slack.New(config.BotAccessToken),
		channelIDs: config.ChannelIDs,
	}

	return notifier, nil
}

func (s *slackNotifier) Notify(event *event.Event) error {
	text := event.Message()

	for _, channelID := range s.channelIDs {
		id, timestamp, err := s.client.PostMessage(
			channelID,
			slack.MsgOptionText(text, false),
			slack.MsgOptionUsername("Comeet"),
		)
		if err != nil {
			return errors.Wrapf(err, "slack: failed sending message to channel %q at %v", id, timestamp)
		}
	}

	return nil
}
