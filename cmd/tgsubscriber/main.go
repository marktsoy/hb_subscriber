package main

import (
	"flag"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/BurntSushi/toml"
	"github.com/marktsoy/tg_subscriber/internal/application"
)

var (
	configPath string
)

func init() {
	fmt.Println("Initializing config")

	flag.StringVar(&configPath, "config-path", "configs/subscriber.toml", "Configuration path")
}

func main() {
	flag.Parse()
	c := application.NewConfig()

	_, err := toml.DecodeFile(configPath, &c)

	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(c.BotKey)
	if err != nil {
		panic(err)
	}
	subscriber := application.New(c, bot)

	subscriber.Run()
}
