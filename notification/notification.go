package notification

import (
	"time"

	"github.com/GGP1/comeet/event"

	"github.com/pkg/errors"
)

var notifiers = map[string]newNotifier{
	"desktop":    func(c Config) (Notifier, error) { return NewDesktopNotifier(c) },
	"join":       func(c Config) (Notifier, error) { return NewJoinNotifier(c) },
	"mail":       func(c Config) (Notifier, error) { return NewMailNotifier(c.Mail) },
	"rocketchat": func(c Config) (Notifier, error) { return NewRocketChatNotifier(c.RocketChat) },
	"slack":      func(c Config) (Notifier, error) { return NewSlackNotifier(c.Slack) },
	"telegram":   func(c Config) (Notifier, error) { return NewTelegramNotifier(c.Telegram) },
}

type newNotifier func(c Config) (Notifier, error)

// Notifier is the interface that wraps the Notify method used to send notifications.
type Notifier interface {
	Notify(event *event.Event) error
}

// Notification ..
type Notification struct {
	// TODO: change to Notifiers?
	Services []string      `yaml:"services,omitempty"`
	Delta    time.Duration `yaml:"delta,omitempty"`
}

// GetNotifiers returns the notifiers the user had set up.
func GetNotifiers(config Config, services []string) ([]Notifier, error) {
	list := make([]Notifier, 0, len(services))

	for _, service := range services {
		newNotifier, ok := notifiers[service]
		if !ok {
			return nil, errors.Errorf("service %q is not supported", service)
		}

		notifier, err := newNotifier(config)
		if err != nil {
			return nil, err
		}

		list = append(list, notifier)
	}

	return list, nil
}
