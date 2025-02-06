package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

const (
	queueName    = "events"
	exchangeName = "events_exchange"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(cfg config.RabbitConf) *Client {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to rabbit: %v", err))
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("failed to open channel: %v", err))
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to declare queue %s: %v", queueName, err))
	}

	if err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		panic(fmt.Sprintf("failed to declare exchange %s: %v", exchangeName, err))
	}

	if err = ch.QueueBind(
		queueName,
		"",
		exchangeName,
		false,
		nil,
	); err != nil {
		panic(fmt.Sprintf("failed to bind exchange %s to queue %s: %v", exchangeName, queueName, err))
	}

	return &Client{
		conn:    conn,
		channel: ch,
	}
}

func (c *Client) Close() error {
	err := c.channel.Close()
	if err != nil {
		return err
	}

	return c.conn.Close()
}

func (c *Client) Add(message app.Notification) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = c.channel.Publish(
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (c *Client) Get() <-chan amqp.Delivery {
	consume, err := c.channel.Consume(
		queueName,
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	return consume
}
