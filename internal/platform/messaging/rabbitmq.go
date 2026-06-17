package messaging

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ExchangePayments = "payments.events"

type Event struct {
	Name     string          `json:"name"`
	ID       string          `json:"id"`
	Occurred time.Time       `json:"occurred_at"`
	Payload  json.RawMessage `json:"payload"`
}

type Publisher interface {
	Publish(ctx context.Context, routingKey string, payload any) error
}

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(context.Context, string, any) error { return nil }

func Connect(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := ch.ExchangeDeclare(ExchangePayments, "topic", true, false, false, false, nil); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &RabbitMQ{conn: conn, ch: ch}, nil
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, payload any) error {
	body, err := json.Marshal(Event{Name: routingKey, Occurred: time.Now().UTC(), Payload: mustRaw(payload)})
	if err != nil {
		return err
	}
	return r.ch.PublishWithContext(ctx, ExchangePayments, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		Body:         body,
	})
}

func (r *RabbitMQ) Consume(ctx context.Context, queue string, keys []string, handler func(context.Context, Event) error) error {
	if _, err := r.ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}
	for _, key := range keys {
		if err := r.ch.QueueBind(queue, key, ExchangePayments, false, nil); err != nil {
			return err
		}
	}
	msgs, err := r.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-msgs:
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				_ = msg.Nack(false, false)
				continue
			}
			if err := handler(ctx, event); err != nil {
				slog.Error("event handler failed", "event", event.Name, "error", err)
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
	}
}

func (r *RabbitMQ) Close() {
	if r.ch != nil {
		_ = r.ch.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}

func mustRaw(payload any) json.RawMessage {
	body, _ := json.Marshal(payload)
	return body
}
