package google

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcalendar "google.golang.org/api/calendar/v3"
)

const redirectURL = "urn:ietf:wg:oauth:2.0:oob"

// Account represents a Google account.
type Account struct {
	// domainWhitelist is inherited from the main configuration
	domainWhitelist map[string]struct{} `yaml:"-,omitempty"`
	oauth           *oauth2.Config      `yaml:"-,omitempty"`
	ClientID        string              `yaml:"clientID,omitempty"`
	ClientSecret    string              `yaml:"clientSecret,omitempty"`
	TokenPath       string              `yaml:"tokenPath,omitempty"`
	CalendarID      string              `yaml:"calendarID,omitempty"`
}

// Config represents Google's configuration.
type Config struct {
	Accounts []*Account `yaml:"accounts,omitempty"`
	Enabled  bool       `yaml:"enabled,omitempty"`
}

// SetDefaultValues populates the configuration with pre-defined values.
func (c *Config) SetDefaultValues(domainWhitelist map[string]struct{}) error {
	for _, account := range c.Accounts {
		if account.TokenPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return errors.Wrap(err, "finding home directory")
			}

			account.TokenPath = filepath.Join(home, fmt.Sprintf(".comeet_token_%s.json", account.ClientID))
		}

		if account.CalendarID == "" {
			account.CalendarID = "primary"
		}

		account.oauth = &oauth2.Config{
			ClientID:     account.ClientID,
			ClientSecret: account.ClientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{gcalendar.CalendarEventsReadonlyScope},
		}
		account.domainWhitelist = domainWhitelist
	}

	return nil
}
