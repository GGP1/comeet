package executor

import (
	"log"
	"sort"
	"time"

	"github.com/GGP1/comeet/event"
	"github.com/GGP1/comeet/notification"
)

// Executor is the interface that wraps the Execute method to run the actions of an event.
type Executor interface {
	Execute(event *event.Event, date time.Time, abort <-chan struct{})
}

// Stage represents a point in time prior to the event start where one or many notifications are sent.
type Stage struct {
	Notifiers []notification.Notifier
	Delta     time.Duration
}

type executor struct {
	finishedEvents chan<- string
	stages         []Stage
}

// New returns an object satisfying the executor interface.
func New(nConfig notification.Config, finishedEvents chan<- string) (Executor, error) {
	stages, err := buildStages(nConfig)
	if err != nil {
		return nil, err
	}

	executor := &executor{
		stages:         stages,
		finishedEvents: finishedEvents,
	}

	return executor, nil
}

// Execute runs an event's stages when their target time is reached.
func (e *executor) Execute(event *event.Event, date time.Time, abort <-chan struct{}) {
	timer := time.NewTimer(0)

	for _, s := range e.stages {
		targetTime := date.Add(-s.Delta)
		now := time.Now()
		// Leave a small window just in case we got multiple stages executed
		// at the same delta, so not to skip the queued stage/s.
		if targetTime.Before(now.Add(1 * time.Second)) {
			continue
		}

		timer.Reset(targetTime.Sub(now))

		select {
		case <-timer.C:
			for _, notifier := range s.Notifiers {
				// Do not wait for their execution
				go func(notifier notification.Notifier) {
					if err := notifier.Notify(event); err != nil {
						log.Println(err)
					}
				}(notifier)
			}

		case <-abort:
			timer.Stop()
			return
		}
	}

	// Wait until the event ends to marked it as finishd so it's kept
	// in the client's cache and we avoid re-scheduling it
	endTimer := time.AfterFunc(event.EndDate.Sub(time.Now()), func() {
		e.finishedEvents <- event.ID
	})
	<-endTimer.C
}

func buildStages(nConfig notification.Config) ([]Stage, error) {
	stages := make([]Stage, 0, len(nConfig.Notifications))

	for _, notif := range nConfig.Notifications {
		notifiers, err := notification.GetNotifiers(nConfig, notif.Services)
		if err != nil {
			return nil, err
		}

		stage := Stage{
			Notifiers: notifiers,
			Delta:     notif.Delta,
		}
		stages = append(stages, stage)
	}

	// Sort stages from older to newer
	sort.SliceStable(stages, func(i, j int) bool {
		return stages[i].Delta > stages[j].Delta
	})

	return stages, nil
}
