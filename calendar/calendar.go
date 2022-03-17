package calendar

import (
	"hash/crc32"
	"log"
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

// Poller ..
type Poller interface {
	Start() error
}

type poller struct {
	scheduler       scheduler.Scheduler
	finishedEvents  <-chan string
	scheduledEvents map[string]uint32
	services        []Service
}

// NewPoller returns a new calendar client.
func NewPoller(scheduler scheduler.Scheduler, finishedEvents <-chan string, services ...Service) Poller {
	return &poller{
		scheduler:       scheduler,
		services:        services,
		finishedEvents:  finishedEvents,
		scheduledEvents: make(map[string]uint32),
	}
}

// Run fetches and schedules events from third party services every 15 minutes.
func (c *poller) Start() error {
	if err := c.scheduleEvents(); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			if err := c.scheduleEvents(); err != nil {
				log.Println(err)
			}

		case eventID := <-c.finishedEvents:
			delete(c.scheduledEvents, eventID)
		}
	}
}

// scheduleEvents sends non-scheduled or updated events to the scheduler.
//
// Concurrency is not used as it provides more drawbacks than benefits.
func (c *poller) scheduleEvents() error {
	for _, service := range c.services {
		events, err := service.GetEvents()
		if err != nil {
			return errors.Wrap(err, "failed fetching events")
		}

		for _, event := range events {
			checksum := crc32.ChecksumIEEE([]byte(event.String()))

			v, ok := c.scheduledEvents[event.ID]
			if ok && checksum == v {
				continue
			}

			c.scheduledEvents[event.ID] = checksum
			c.scheduler.Schedule(event)
		}
	}

	return nil
}
