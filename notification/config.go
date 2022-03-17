package notification

// Config contains notifications and messaging services configurations.
type Config struct {
	Mail          MailConfig       `yaml:"mail,omitempty"`
	Notifications []Notification   `yaml:"notifications,omitempty"`
	RocketChat    RocketChatConfig `yaml:"rocketchat,omitempty"`
	Slack         SlackConfig      `yaml:"slack,omitempty"`
	Telegram      TelegramConfig   `yaml:"telegram,omitempty"`
}

// MailConfig contains a mail's configuration.
type MailConfig struct {
	SenderAddress     string   `yaml:"senderAddress,omitempty"`
	SenderPassword    string   `yaml:"senderPassword,omitempty"`
	SMTPHostAddress   string   `yaml:"smtpHostAddress,omitempty"`
	SMTPHostPort      string   `yaml:"smtpHostPort,omitempty"`
	ReceiverAddresses []string `yaml:"receiverAddresses,omitempty"`
}

// RocketChatConfig contains rocket chat's configuration.
type RocketChatConfig struct {
	ServerURL string   `yaml:"serverURL,omitempty"`
	UserID    string   `yaml:"userID,omitempty"`
	Token     string   `yaml:"token,omitempty"`
	RoomIDs   []string `yaml:"roomIDs,omitempty"`
}

// SlackConfig contains slack's configuration.
type SlackConfig struct {
	BotAccessToken string   `yaml:"botAccessToken,omitempty"`
	ChannelIDs     []string `yaml:"channelIDs,omitempty"`
}

// TelegramConfig contains telegram's configuration.
type TelegramConfig struct {
	BotAPIToken string  `yaml:"botAPIToken"`
	ChatIDs     []int64 `yaml:"chatID"`
}
