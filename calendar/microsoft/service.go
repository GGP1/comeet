package microsoft

import (
	"net/url"
	"time"

	"github.com/GGP1/comeet/calendar"
	"github.com/GGP1/comeet/event"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	eventss "github.com/microsoftgraph/msgraph-sdk-go/me/calendars/item/events"
	"github.com/microsoftgraph/msgraph-sdk-go/models/microsoft/graph"
	"github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar/events"
	"github.com/pkg/errors"
)

type service struct {
	client  *msgraphsdk.GraphServiceClient
	account *Account
}

// NewServices returns n microsoft calendar services.
func NewServices(config Config, domainWhitelist map[string]struct{}) ([]calendar.Service, error) {
	config.SetDefaultValues(domainWhitelist)

	services := make([]calendar.Service, 0, len(config.Accounts))

	for _, account := range config.Accounts {
		service, err := NewService(account)
		if err != nil {
			return nil, errors.Wrapf(err, "microsoft client %q", account.ClientID)
		}

		services = append(services, service)
	}

	return services, nil
}

// NewService returns a new Microsoft calendar service.
func NewService(account *Account) (calendar.Service, error) {
	client, err := getClient(account)
	if err != nil {
		return nil, errors.Wrap(err, "creating microsoft calendar client")
	}

	service := &service{
		client:  client,
		account: account,
	}
	return service, nil
}

// GetEvents ..
func (s *service) GetEvents() ([]*event.Event, error) {
	mEvents := make([]graph.Event, 0)

	// The package used to request the events changes depending on where a calendar
	// ID is used or not
	if s.account.CalendarID == "" {
		options := &events.EventsRequestBuilderGetOptions{
			Q: &events.EventsRequestBuilderGetQueryParameters{
				Top:    &s.account.maxResults,
				Select: s.account.fields,
			},
		}
		result, err := s.client.UsersById("ggpalomeque@hotmail.com").Calendar().Events().Get(options)
		if err != nil {
			return nil, err
		}
		mEvents = result.GetValue()
	} else {
		options := &eventss.EventsRequestBuilderGetOptions{
			Q: &eventss.EventsRequestBuilderGetQueryParameters{
				Top:    &s.account.maxResults,
				Select: s.account.fields,
			},
		}
		result, err := s.client.Me().CalendarsById(s.account.CalendarID).Events().Get(options)
		if err != nil {
			return nil, err
		}
		mEvents = result.GetValue()
	}

	events := make([]*event.Event, 0, len(mEvents))
	for _, mEvent := range mEvents {
		event, err := s.toComeetEvent(mEvent)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *service) toComeetEvent(e graph.Event) (*event.Event, error) {
	// TODO: is this relevant?
	// if !*e.GetIsOnlineMeeting() {
	// 	return nil, errors.New("the event is not an online meeting")
	// }

	startDate, err := time.Parse(time.RFC3339, *e.GetStart().GetDateTime())
	if err != nil {
		return nil, err
	}
	// Microsoft's SDK does not allow filtering by start time
	if startDate.Before(time.Now()) || startDate.After(time.Now().Add(event.Window)) {
		return nil, errors.New("the event is outside the stipulated window")
	}

	endDate, err := time.Parse(time.RFC3339, *e.GetEnd().GetDateTime())
	if err != nil {
		return nil, err
	}

	url, err := parseEventURL(e, s.account.domainWhitelist)
	if err != nil {
		return nil, err
	}

	event := &event.Event{
		ID:          *e.GetId(),
		Title:       *e.GetSubject(),
		Description: *e.GetBodyPreview(),
		StartDate:   startDate,
		EndDate:     endDate,
		URL:         url,
		// Microsoft's calendar events can't be repeated within a day
		// and we are just getting the next 60 minutes
		Recurrence: "",
	}
	return event, nil
}

func parseEventURL(e graph.Event, domainWhitelist map[string]struct{}) (*url.URL, error) {
	uri := firstNonEmptyURL(
		e.GetOnlineMeetingUrl(),
		e.GetOnlineMeeting().GetJoinUrl(),
		e.GetLocation().GetLocationUri(),
	)
	if uri == "" {
		return &url.URL{}, nil
	}

	url, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "event %q has an invalid url: %q", *e.GetId(), uri)
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

func firstNonEmptyURL(urls ...*string) string {
	for _, url := range urls {
		if *url != "" {
			return *url
		}
	}

	return ""
}
