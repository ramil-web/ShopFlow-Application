package services

import (
	"encoding/json"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationService struct {
	MQConn *amqp.Connection
}

type ApplicationCreatedMessage struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`
	File   string `json:"file_url"`
	Email  string `json:"email"`
}

func NewEventPublisher(conn *amqp.Connection) *NotificationService {
	return &NotificationService{MQConn: conn}
}

var ErrNoMQConnection = errors.New("no rabbitmq connection")

// PublishApplicationCreated публикует событие о созданной заявке в очередь "application_created"
func (s *NotificationService) PublishApplicationCreated(msg ApplicationCreatedMessage) error {
	if s == nil || s.MQConn == nil {
		return ErrNoMQConnection
	}

	ch, err := s.MQConn.Channel()
	if err != nil {
		log.Println("[notification] failed to open channel:", err)
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(msg)
	if err != nil {
		log.Println("[notification] failed to marshal message:", err)
		return err
	}

	exchangeName := "shopflow.events" // или os.Getenv("EXCHANGE_NAME")
	routingKey := "application_created"

	if err := ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		log.Println("[notification] failed to publish message:", err)
		return err
	}

	log.Printf("[notification] Published %s id=%d user_id=%d\n", routingKey, msg.ID, msg.UserID)
	return nil
}
