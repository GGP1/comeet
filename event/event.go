package event

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
)

// Window represents the size of the window that contains the events being considered.
const Window = time.Hour

// Event represents a comeet event.
type Event struct {
	ID          string
	StartDate   time.Time
	EndDate     time.Time
	URL         *url.URL
	Title       string
	Description string
	Recurrence  string
}

// Dates returns a event's dates based on its recurrence.
//
// In case it has none, it returns the start date.
func (e *Event) Dates() ([]time.Time, error) {
	if e.Recurrence == "" {
		return []time.Time{e.StartDate}, nil
	}

	r, err := rrule.StrToRRule(e.Recurrence)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dates := r.Between(now, now.Add(Window), true)

	return dates, nil
}

// Message returns a formatted text to put in a message.
func (e *Event) Message() string {
	var sb *strings.Builder
	sb.WriteString(e.Title)
	sb.WriteString("\n\n")

	untilStart := time.Until(e.StartDate).Round(time.Second)
	var startTime string
	if untilStart < time.Second*2 {
		startTime = "Starting now"
	} else {
		startTime = fmt.Sprintf("Starts at: %s (%s left)", e.StartDate.Format(time.Kitchen), untilStart)
	}
	sb.WriteString(startTime)
	sb.WriteString("\n\n")

	if e.URL != nil {
		sb.WriteString(startTime)
		sb.WriteString("\n\n")
	}

	sb.WriteString(e.Description)

	return sb.String()
}

// String returns the string representation of an event's content.
func (e *Event) String() string {
	return e.ID +
		e.StartDate.Format(time.RFC3339) +
		e.EndDate.Format(time.RFC3339) +
		e.URL.Redacted() +
		e.Title +
		e.Description +
		e.Recurrence
}
