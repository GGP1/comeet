package notification

import (
	"strings"

	"github.com/GGP1/comeet/event"

	"github.com/gen2brain/beeep"
)

type desktopNotifier struct{}

// NewDesktopNotifier returns a notifier for desktop notifications.
func NewDesktopNotifier(config Config) (Notifier, error) {
	return &desktopNotifier{}, nil
}

func (d *desktopNotifier) Notify(event *event.Event) error {
	sep := "\n\n"
	description := strings.SplitAfter(event.Message(), sep)
	// TODO: take icon from the configuration file?
	return beeep.Notify(event.Title, strings.Join(description[1:], sep), "")
}
