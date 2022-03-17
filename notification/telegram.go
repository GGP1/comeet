package notification

import (
	"github.com/GGP1/comeet/event"
	"github.com/pkg/errors"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramNotifier struct {
	botAPI  *telegram.BotAPI
	chatIDs []int64
}

// NewTelegramNotifier returns a notifier that sends telegram messages.
func NewTelegramNotifier(config TelegramConfig) (Notifier, error) {
	botAPI, err := telegram.NewBotAPI(config.BotAPIToken)
	if err != nil {
		return nil, errors.Wrap(err, "creating telegram bot API")
	}
	notifier := &telegramNotifier{
		botAPI:  botAPI,
		chatIDs: config.ChatIDs,
	}
	return notifier, nil
}

func (t *telegramNotifier) Notify(event *event.Event) error {
	message := telegram.NewMessage(0, event.Message())
	message.ParseMode = telegram.ModeMarkdownV2
	message.ChannelUsername = "Comeet"

	for _, chatID := range t.chatIDs {
		message.ChatID = chatID

		if _, err := t.botAPI.Send(message); err != nil {
			return errors.Wrapf(err, "sending message to chat %d", chatID)
		}
	}

	return nil
}
