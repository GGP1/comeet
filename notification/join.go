package notification

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/GGP1/comeet/event"

	"github.com/gen2brain/dlgs"
)

type joinNotifier struct{}

// NewJoinNotifier returns a notifier that displays a question dialog asking to join an online event or not.
func NewJoinNotifier(config Config) (Notifier, error) {
	return &joinNotifier{}, nil
}

func (s *joinNotifier) Notify(event *event.Event) error {
	if event.URL == nil {
		log.Printf("Join notifier skipped: event %q doesn't have a URL\n", event.Title)
		return nil
	}

	title := event.Title
	if time.Until(event.StartDate) < 1 {
		title += " is starting"
	}

	text := fmt.Sprintf(
		"%s\n\n%s - %s\n\n%s\n\nWould you like to connect?",
		event.URL,
		event.StartDate.Format(time.Kitchen),
		event.EndDate.Format(time.Kitchen),
		event.Description,
	)
	ok, err := dlgs.Question(title, text, true)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// For now we are just opening the browser but eventually
	// we could take other instructions from the user
	return openBrowser(event.URL)
}

func openBrowser(url *url.URL) error {
	var cmd *exec.Cmd
	urlStr := url.String()

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", urlStr)

	case "windows":
		cmd = exec.Command("cmd", "/c", "start", urlStr)

	case "darwin":
		cmd = exec.Command("open", urlStr)

	default:
		log.Println("Couldn't open the web browser: unsupported platform")
	}

	return cmd.Start()
}
