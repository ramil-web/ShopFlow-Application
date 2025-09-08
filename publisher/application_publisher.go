package publisher

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ApplicationPublisher struct {
	MQConn *amqp.Connection
}

type ApplicationCreatedMessage struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Text   string `json:"text"`
	File   string `json:"file_url"`
}

func NewApplicationPublisher(conn *amqp.Connection) *ApplicationPublisher {
	return &ApplicationPublisher{MQConn: conn}
}

func (p *ApplicationPublisher) PublishApplicationCreated(msg ApplicationCreatedMessage) error {
	ch, err := p.MQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"application_created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(msg)

	if err := ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return err
	}

	log.Printf("[publisher] Published application_created: app_id=%d user_id=%d\n", msg.ID, msg.UserID)
	return nil
}
