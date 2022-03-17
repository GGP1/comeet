package microsoft

// Account represents a Microsoft account.
type Account struct {
	CalendarID string `yaml:"calendarID,omitempty"`
	// https://docs.microsoft.com/en-us/onedrive/find-your-office-365-tenant-id
	TenantID     string `yaml:"tenantID,omitempty"`
	ClientID     string `yaml:"clientID,omitempty"`
	ClientSecret string `yaml:"clientSecret,omitempty"`
	// domainWhitelist is inherited from the main configuration
	domainWhitelist map[string]struct{} `yaml:"-,omitempty"`
	fields          []string            `yaml:"-,omitempty"`
	maxResults      int32               `yaml:"-,omitempty"`
}

// Config represents Microsoft's configuration.
type Config struct {
	Accounts []*Account `yaml:"accounts,omitempty"`
	Enabled  bool       `yaml:"enabled,omitempty"`
}

// SetDefaultValues populates the configuration with pre-defined values.
func (c *Config) SetDefaultValues(domainWhitelist map[string]struct{}) {
	for _, account := range c.Accounts {
		account.fields = []string{
			"id", "start", "end", "subject", "bodyPreview",
			"location", "onlineMeetingURL", "onlineMeeting",
		}
		account.maxResults = 60
		account.domainWhitelist = domainWhitelist
	}
}
