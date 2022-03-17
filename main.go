package main

import (
	"log"
	"os"

	"github.com/GGP1/comeet/calendar"

	"github.com/GGP1/comeet/config"
	"github.com/GGP1/comeet/executor"
	"github.com/GGP1/comeet/scheduler"

	"github.com/emersion/go-autostart"
	"github.com/pkg/errors"
)

func main() {
	if err := setToRunOnStartup(); err != nil {
		log.Fatal(errors.Wrap(err, "failed setting service to run on startup"))
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed creating configuration"))
	}

	services, err := config.GetCalendarServices(cfg)
	if err != nil {
		log.Fatal(err)
	}

	finishedEvents := make(chan string)
	executor, err := executor.New(cfg.Notification, finishedEvents)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed creating executor"))
	}

	log.Println("Comeet is starting...")
	scheduler := scheduler.New(executor, finishedEvents)
	go scheduler.Run()

	poller := calendar.NewPoller(scheduler, finishedEvents, services...)
	if err := poller.Start(); err != nil {
		log.Fatal(errors.Wrap(err, "failed running client"))
	}
}

func setToRunOnStartup() error {
	name := "comeet"
	app := &autostart.App{
		DisplayName: name,
		Name:        name,
	}
	if app.IsEnabled() {
		return nil
	}

	binPath, err := os.Executable()
	if err != nil {
		return err
	}

	app.Exec = []string{binPath}
	return app.Enable()
}
