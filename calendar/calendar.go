package calendar

import (
	"hash/crc32"
	"log"
	"sync"
	"time"

	"github.com/GGP1/comeet/event"
	"github.com/GGP1/comeet/scheduler"

	"github.com/pkg/errors"
)

const interval = time.Minute * 15

// Service represents a third party calendar service.
type Service interface {
	GetEvents() ([]*event.Event, error)
}

// Poller represents an entity capable of polling events from third party services.
type Poller interface {
	Start() error
}

type poller struct {
	scheduler       scheduler.Scheduler
	finishedEvents  <-chan string
	scheduledEvents map[string]uint32
	mu              *sync.Mutex
	services        []Service
}

// NewPoller returns an object that looks for and sends events to the scheduler.
func NewPoller(scheduler scheduler.Scheduler, finishedEvents <-chan string, services ...Service) Poller {
	return &poller{
		scheduler:       scheduler,
		services:        services,
		finishedEvents:  finishedEvents,
		scheduledEvents: make(map[string]uint32),
	}
}

// Start fetches and schedules events from third party services every x minutes.
func (c *poller) Start() error {
	c.fetchEvents()

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			c.fetchEvents()

		case eventID := <-c.finishedEvents:
			c.mu.Lock()
			delete(c.scheduledEvents, eventID)
			c.mu.Unlock()
		}
	}
}

// fetchEvents sets each one of the services to fetch and schedule events.
func (c *poller) fetchEvents() {
	for _, service := range c.services {
		go c.scheduleEvents(service)
	}
}

// scheduleEvents gets the events from a service and sends non-scheduled or updated ones to the scheduler.
func (c *poller) scheduleEvents(service Service) {
	events, err := service.GetEvents()
	if err != nil {
		log.Println(errors.Wrap(err, "failed fetching events"))
	}

	c.mu.Lock()
	for _, event := range events {
		checksum := crc32.ChecksumIEEE([]byte(event.String()))

		v, ok := c.scheduledEvents[event.ID]
		if ok && checksum == v {
			continue
		}

		c.scheduledEvents[event.ID] = checksum
		c.scheduler.Schedule(event)
	}
	c.mu.Unlock()
}
