package config

import (
	"net/mail"
	"net/url"
	"os"

	"github.com/GGP1/comeet/calendar"
	"github.com/GGP1/comeet/calendar/google"
	"github.com/GGP1/comeet/calendar/microsoft"
	"github.com/GGP1/comeet/event"
	"github.com/GGP1/comeet/notification"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var errNoConfig = errors.New("configuration file required, please specify its path in the \"COMEET_CONFIG\" environment variable")

// Config represents comeet's configuration.
type Config struct {
	Notification    notification.Config `yaml:"notification,omitempty"`
	DomainWhitelist []string            `yaml:"domainWhitelist,omitempty"`
	Google          google.Config       `yaml:"google,omitempty"`
	Microsoft       microsoft.Config    `yaml:"microsoft,omitempty"`
}

// New returns a new comeet's configuration.
func New() (Config, error) {
	config, err := load()
	if err != nil {
		return Config{}, err
	}

	if err := validate(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// GetCalendarServices returns third party calendar services configurations.
func GetCalendarServices(config Config) ([]calendar.Service, error) {
	domainWhitelist := sliceToMap(config.DomainWhitelist)
	accountsNum := len(config.Google.Accounts) + len(config.Microsoft.Accounts)
	services := make([]calendar.Service, 0, accountsNum)

	if config.Google.Enabled {
		gServices, err := google.NewServices(config.Google, domainWhitelist)
		if err != nil {
			return nil, err
		}
		services = append(services, gServices...)
	}
	if config.Microsoft.Enabled {
		mServices, err := microsoft.NewServices(config.Microsoft, domainWhitelist)
		if err != nil {
			return nil, err
		}
		services = append(services, mServices...)
	}

	return services, nil
}

func load() (Config, error) {
	configPath := os.Getenv("COMEET_CONFIG")
	if configPath == "" {
		return Config{}, errNoConfig
	}

	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "opening configuration file")
	}
	defer f.Close()

	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, errors.Wrap(err, "decoding configuration file")
	}

	return config, nil
}

func sliceToMap(slice []string) map[string]struct{} {
	mp := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		mp[s] = struct{}{}
	}

	return mp
}

func validate(c Config) error {
	if !c.Google.Enabled && !c.Microsoft.Enabled {
		return errors.New("at least one calendar service is required")
	}

	if c.Google.Enabled {
		for _, acc := range c.Google.Accounts {
			if acc.ClientID == "" || acc.ClientSecret == "" {
				return errors.Errorf("Google client %q has invalid credentials", acc.ClientID)
			}
		}
	}

	for _, notif := range c.Notification.Notifications {
		if len(notif.Services) == 0 {
			return errors.New("no notification services specified")
		}
		if notif.Delta > event.Window {
			return errors.Errorf(
				"invalid delta value, maximum allowed is %v since only events within that window are fetched",
				event.Window,
			)
		}
	}

	if c.Notification.Mail.SenderAddress != "" {
		if _, err := url.Parse(c.Notification.Mail.SMTPHostAddress); err != nil {
			return errors.Wrap(err, "invalid smtp host address")
		}

		if c.Notification.Mail.SMTPHostPort < "0" || c.Notification.Mail.SMTPHostPort > "65535" {
			return errors.New("invalid smtp host port number")
		}

		addresses := append(c.Notification.Mail.ReceiverAddresses, c.Notification.Mail.SenderAddress)
		for _, address := range addresses {
			if address == "" {
				continue
			}
			if _, err := mail.ParseAddress(address); err != nil {
				return errors.Wrapf(err, "invalid mail address: %q", address)
			}
		}
	}

	return nil
}
