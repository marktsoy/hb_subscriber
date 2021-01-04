package application

import (
	"fmt"
	"log"

	"github.com/hamba/avro"
	"github.com/marktsoy/tg_subscriber/internal/models"
	"github.com/streadway/amqp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Subscriber ...
type Subscriber struct {
	bot              *tgbotapi.BotAPI
	rabbitConnection *amqp.Connection
	channel          *amqp.Channel
	queue            amqp.Queue
}

// New ...
func New(c *Config, bot *tgbotapi.BotAPI) *Subscriber {

	conn, err := amqp.Dial(c.RabbitAddr)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(c.SubscribersQueueName, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	return &Subscriber{
		bot:              bot,
		rabbitConnection: conn,
		queue:            q,
		channel:          ch,
	}
}

// Run Subscriber
func (s *Subscriber) Run() {

	s.bot.Debug = true

	log.Printf("Authorized on account %s", s.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err.Error())
	}
	for update := range updates {
		subscription := models.Subscription{}
		respText := ""
		switch update.Message.Text {
		default:
			continue
		case "/subscribe":
			subscription.IsSubscribed = true
			respText = "Subscribed"
		case "/unsubscribe":
			subscription.IsSubscribed = false
			respText = "Unsubscribed"
		}
		subscription.ChatID = update.Message.Chat.ID
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, respText)
		msg.ReplyToMessageID = update.Message.MessageID
		s.bot.Send(msg)
		s.pub(subscription)
	}

}

func (s *Subscriber) pub(subscription models.Subscription) {
	schema := models.SchemaSubscription()

	data, err := avro.Marshal(schema, subscription)
	if err != nil {
		fmt.Println(err.Error())
	}

	s.channel.Publish(
		"",           // exchange
		s.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
}
