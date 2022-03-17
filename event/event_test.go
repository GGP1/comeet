package event_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/GGP1/comeet/event"

	"github.com/stretchr/testify/assert"
)

func TestDates(t *testing.T) {
	t.Run("Recurrence", func(t *testing.T) {
		now := time.Now()
		event := &event.Event{
			Recurrence: "FREQ=MINUTELY;INTERVAL=15;COUNT=5",
			StartDate:  now,
			EndDate:    now.Add(time.Hour),
		}
		// Note: varies depending on event.StartGap
		expectedLen := 4

		actual, err := event.Dates()
		assert.NoError(t, err)

		assert.Equal(t, expectedLen, len(actual))
	})

	t.Run("No recurrence", func(t *testing.T) {
		date := time.Date(2022, 02, 04, 00, 00, 00, 00, time.UTC)
		event := &event.Event{
			Recurrence: "",
			StartDate:  date,
			EndDate:    date.Add(24 * time.Hour),
		}
		expected := []time.Time{date}

		actual, err := event.Dates()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestString(t *testing.T) {
	date := time.Date(2022, 06, 04, 17, 30, 00, 00, time.UTC)
	e := &event.Event{
		ID:          "123456789",
		Title:       "Test",
		Description: "description",
		Recurrence:  "FREQ=YEARLY;COUNT=2",
		URL: &url.URL{
			Host:   "meet.com",
			Scheme: "https",
		},
		StartDate: date,
		EndDate:   date.Add(24 * time.Hour),
	}

	expected := "1234567892022-06-04T17:30:00Z2022-06-05T17:30:00Zhttps://meet.comTestdescriptionFREQ=YEARLY;COUNT=2"
	actual := e.String()

	assert.Equal(t, expected, actual)
}
