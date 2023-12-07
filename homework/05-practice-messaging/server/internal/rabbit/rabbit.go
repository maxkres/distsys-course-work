package rabbit

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"server/config"
	"time"
)

type Rabbit struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	q        amqp.Queue
	messages <-chan amqp.Delivery

	delivery chan string

	cfg config.Config
}

func New(cfg config.Config, del chan string) (Rabbit, error) {
	conn, err := amqp.Dial(cfg.GetRabbitUrl()) // Создаем подключение к RabbitMQ
	if err != nil {
		return Rabbit{}, fmt.Errorf("dial failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return Rabbit{}, fmt.Errorf("channel failed: %w", err)
	}

	q, err := ch.QueueDeclare(
		cfg.RabbitQueueName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return Rabbit{}, fmt.Errorf("queue declare failed: %w", err)
	}

	q, err = ch.QueueDeclare(
		cfg.RabbitQueueName+"ids", // name
		false,                     // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return Rabbit{}, fmt.Errorf("queue ids declare failed: %w", err)
	}

	messages, err := ch.Consume(
		cfg.RabbitQueueName+"ids", // queue
		"",                        // consumer
		true,                      // auto-ack
		false,                     // exclusive
		false,                     // no-local
		false,                     // no-wait
		nil,                       // args
	)
	if err != nil {
		return Rabbit{}, fmt.Errorf("consume failed: %w", err)
	}

	return Rabbit{
		conn:     conn,
		ch:       ch,
		q:        q,
		cfg:      cfg,
		delivery: del,
		messages: messages,
	}, nil
}

func (r Rabbit) Close() error {
	err := r.conn.Close()
	if err != nil {
		return fmt.Errorf("conn close failed: %w", err)
	}
	err = r.ch.Close()
	if err != nil {
		return fmt.Errorf("ch close failed: %w", err)
	}
	return nil
}

func (r Rabbit) Publish(ctx context.Context, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.ch.PublishWithContext(ctx,
		"",                    // exchange
		r.cfg.RabbitQueueName, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("publish with context failed: %w", err)
	}

	return nil
}

func (r Rabbit) Receive() {
	go func() {
		for message := range r.messages {
			r.delivery <- string(message.Body)
		}
	}()
}
