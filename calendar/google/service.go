package google

import (
	"context"
	"net/url"
	"time"

	"github.com/GGP1/comeet/calendar"
	"github.com/GGP1/comeet/event"

	"github.com/pkg/errors"
	gcalendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type service struct {
	eventsService *gcalendar.EventsService
	account       *Account
}

// NewServices returns n Google calendar services.
func NewServices(config Config, domainWhitelist map[string]struct{}) ([]calendar.Service, error) {
	if err := config.SetDefaultValues(domainWhitelist); err != nil {
		return nil, err
	}

	services := make([]calendar.Service, 0, len(config.Accounts))
	for _, account := range config.Accounts {
		service, err := NewService(account)
		if err != nil {
			return nil, errors.Wrapf(err, "client %q", account.ClientID)
		}

		services = append(services, service)
	}

	return services, nil
}

// NewService returns a new Google calendar service.
func NewService(account *Account) (calendar.Service, error) {
	ctx := context.Background()

	client, err := getClient(ctx, account)
	if err != nil {
		return nil, err
	}

	srv, err := gcalendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.Wrap(err, "creating google calendar service")
	}

	return &service{
		account:       account,
		eventsService: gcalendar.NewEventsService(srv),
	}, nil
}

// GetEvents returns the events retrieved from the Google Calendar API.
func (s *service) GetEvents() ([]*event.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	timeMin := time.Now()
	timeMax := timeMin.Add(event.Window)
	gEvents, err := s.eventsService.
		List(s.account.CalendarID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(timeMax.Format(time.RFC3339)).
		OrderBy("startTime").
		MaxResults(60).
		Context(ctx).
		Do()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving events")
	}

	events := make([]*event.Event, 0, len(gEvents.Items))
	for _, item := range gEvents.Items {
		event, err := s.toComeetEvent(item)
		if err != nil {
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *service) toComeetEvent(e *gcalendar.Event) (*event.Event, error) {
	url, err := parseEventURL(e, s.account.domainWhitelist)
	if err != nil {
		return nil, err
	}

	startDate, err := parseDateTime(e.Start)
	if err != nil {
		return nil, err
	}

	endDate, err := parseDateTime(e.End)
	if err != nil {
		return nil, err
	}

	return &event.Event{
		ID:          e.Id,
		Title:       e.Summary,
		Description: e.Description,
		URL:         url,
		StartDate:   startDate,
		EndDate:     endDate,
		// Google's calendar events can't be repeated within a day
		// and we are just getting the next 60 minutes
		Recurrence: "",
	}, nil
}

func parseDateTime(t *gcalendar.EventDateTime) (time.Time, error) {
	datetime, err := time.Parse(time.RFC3339, t.DateTime)
	if err != nil {
		datetime, err = time.Parse("2006-01-02", t.Date)
		if err != nil {
			return time.Time{}, err
		}
	}

	return datetime, nil
}

func parseEventURL(e *gcalendar.Event, domainWhitelist map[string]struct{}) (*url.URL, error) {
	var uri string
	if e.HangoutLink != "" {
		uri = e.HangoutLink
	} else if e.Location != "" {
		uri = e.Location
	} else {
		return &url.URL{}, nil
	}

	url, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrap(err, "invalid url")
	}

	if url.Scheme == "" {
		url.Scheme = "https"
	} else if url.Scheme != "https" {
		return nil, errors.New("invalid url scheme, only https is accepted")
	}

	if len(domainWhitelist) != 0 {
		if _, ok := domainWhitelist[url.Host]; !ok {
			return nil, errors.Errorf("domain %q is not in the whitelist", url.Host)
		}
	}

	return url, nil
}
