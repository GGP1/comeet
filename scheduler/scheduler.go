package scheduler

import (
	"github.com/GGP1/comeet/event"
	"github.com/GGP1/comeet/executor"
)

// Scheduler ..
type Scheduler interface {
	Run()
	Schedule(event *event.Event)
}

type scheduler struct {
	// scheduled stores an abort channel for each event
	// to exit an scheduled execution
	scheduled map[string]chan struct{}
	// events is used for receiving events from services
	events chan *event.Event
	// finishedEvents is used to send a signal to the service that
	// an event has ended
	finishedEvents <-chan string
	executor       executor.Executor
}

// New returns a new scheduler.
func New(executor executor.Executor, finishedEvents <-chan string) Scheduler {
	return &scheduler{
		events:         make(chan *event.Event),
		scheduled:      make(map[string]chan struct{}),
		finishedEvents: finishedEvents,
		executor:       executor,
	}
}

// Run listens for new events and schedules their execution.
func (s *scheduler) Run() {
	for {
		select {
		case event := <-s.events:
			abort := make(chan struct{}, 1)

			// If the event was updated, trigger the scheduled abort
			// channel and create a new one
			if abortCh, ok := s.scheduled[event.ID]; ok {
				abortCh <- struct{}{}
			}
			s.scheduled[event.ID] = abort

			// Just ignore if the rrule is invalid
			// Most of the time we are getting event.StartDate anyway
			dates, _ := event.Dates()
			for _, date := range dates {
				go s.executor.Execute(event, date, abort)
			}

		case eventID := <-s.finishedEvents:
			delete(s.scheduled, eventID)
		}
	}
}

// Schedule sends an event to the channel the scheduler is listening to.
func (s *scheduler) Schedule(event *event.Event) {
	s.events <- event
}
