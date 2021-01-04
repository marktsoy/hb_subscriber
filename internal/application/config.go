package application

type Config struct {
	BotKey string `toml:"tg_bot_public_key"`

	RabbitAddr           string `toml:"rabbit_addr"`
	SubscribersQueueName string `toml:"rabbit_subscriber_queue"`
}

func NewConfig() *Config {
	return &Config{}
}
