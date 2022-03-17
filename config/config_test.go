package config

import (
	"os"
	"testing"
	"time"

	"github.com/GGP1/comeet/calendar/google"
	"github.com/GGP1/comeet/calendar/microsoft"
	"github.com/GGP1/comeet/notification"

	"github.com/stretchr/testify/assert"
)

const mockConfigPath = "testdata/mock_config.yaml"

var mockConfig = Config{
	DomainWhitelist: []string{"teams.microsoft.com"},
	Microsoft: microsoft.Config{
		Accounts: []*microsoft.Account{
			{
				TenantID:     "tenantID",
				ClientID:     "clientID",
				ClientSecret: "clientSecret",
			},
		},
	},
}

func TestNew(t *testing.T) {
	os.Setenv("COMEET_CONFIG", mockConfigPath)

	got, err := New()
	assert.NoError(t, err)
	assert.Equal(t, mockConfig, got)
}

func TestGetCalendarServices(t *testing.T) {
	// Expect google default configuration being used when its service is present
	cases := []struct {
		desc          string
		config        Config
		expectedCount int
	}{
		{
			desc: "All services",
			config: Config{
				Google: google.Config{
					Enabled: true,
				},
				Microsoft: microsoft.Config{
					Enabled: true,
				},
			},
			expectedCount: 1,
		},
		{
			desc:          "Google",
			config:        Config{Google: google.Config{Enabled: true}},
			expectedCount: 1,
		},
		{
			desc:          "Microsoft",
			config:        Config{Microsoft: microsoft.Config{Enabled: true}},
			expectedCount: 0,
		},
		{
			desc:          "None",
			config:        Config{},
			expectedCount: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			services, err := GetCalendarServices(tc.config)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCount, len(services))
		})
	}
}

func TestLoad(t *testing.T) {
	envKey := "COMEET_CONFIG"

	cases := []struct {
		setConfigPath func()
		desc          string
		expected      Config
		success       bool
	}{
		{
			desc: "Config path set",
			setConfigPath: func() {
				err := os.Setenv(envKey, mockConfigPath)
				assert.NoError(t, err)
			},
			expected: mockConfig,
			success:  true,
		},
		{
			desc: "No config path",
			setConfigPath: func() {
				err := os.Unsetenv(envKey)
				assert.NoError(t, err)
			},
			success: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.setConfigPath()
			got, err := load()
			if tc.success {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestSliceToMap(t *testing.T) {
	cases := []struct {
		expected map[string]struct{}
		desc     string
		slice    []string
	}{
		{
			desc:  "Unique",
			slice: []string{"1", "2"},
			expected: map[string]struct{}{
				"1": {},
				"2": {},
			},
		},
		{
			desc:  "Duplicates",
			slice: []string{"1", "1"},
			expected: map[string]struct{}{
				"1": {},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := sliceToMap(tc.slice)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestValidate(t *testing.T) {
	validConfig := Config{
		Google: google.Config{
			Enabled: true,
			Accounts: []*google.Account{
				{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
		},
		DomainWhitelist: []string{"https://github.com/GGP1/comeet"},
		Notification: notification.Config{
			Notifications: []notification.Notification{
				{Services: []string{"desktop"}, Delta: time.Minute},
			},
			Mail: notification.MailConfig{
				SenderAddress:   "valid@mail.com",
				SMTPHostAddress: "smtp.provider.com",
				SMTPHostPort:    "25",
			},
		},
	}

	cases := []struct {
		desc    string
		config  Config
		success bool
	}{
		{
			desc:    "Valid",
			config:  validConfig,
			success: true,
		},
		{
			desc:    "No calendar services",
			config:  Config{},
			success: false,
		},
		{
			desc: "Invalid client id",
			config: Config{
				Google: google.Config{
					Enabled:  true,
					Accounts: []*google.Account{{ClientSecret: "clientSecret"}},
				},
			},
			success: false,
		},
		{
			desc: "Invalid client secret",
			config: Config{
				Google: google.Config{
					Enabled:  true,
					Accounts: []*google.Account{{ClientID: "clientID"}},
				},
			},
			success: false,
		},
		{
			desc: "No notification services",
			config: Config{
				Notification: notification.Config{
					Notifications: []notification.Notification{
						{Services: []string{}, Delta: time.Minute},
					},
				},
			},
			success: false,
		},
		{
			desc: "Invalid notification delta",
			config: Config{
				Notification: notification.Config{
					Notifications: []notification.Notification{
						{Services: []string{"mail"}, Delta: 120 * time.Minute},
					},
				},
			},
			success: false,
		},
		{
			desc: "Invalid host address",
			config: Config{
				Notification: notification.Config{
					Mail: notification.MailConfig{
						SenderAddress:   "1",
						SMTPHostAddress: "invalidAddresscom",
					},
				},
			},
			success: false,
		},
		{
			desc: "Invalid host port",
			config: Config{
				Notification: notification.Config{
					Mail: notification.MailConfig{
						SenderAddress:   "1",
						SMTPHostAddress: validConfig.Notification.Mail.SMTPHostAddress,
						SMTPHostPort:    "75000",
					},
				},
			},
			success: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := validate(tc.config)
			if tc.success {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
