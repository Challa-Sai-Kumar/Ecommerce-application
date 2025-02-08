package kafka

import (
	"context"
	"ecommerce/models"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	topic  string
	writer *kafka.Writer
}

// NewProducer initializes and returns a Kafka producer.
func NewProducer(broker []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: broker,
		Topic:   topic,
	})

	return &Producer{
		topic:  topic,
		writer: writer,
	}
}

// PublishOrderStatus sends an order status update to Kafka.
func (p *Producer) PublishOrderStatus(orderID, status string) error {
	value, err := json.Marshal(map[string]string{
		"order_id": orderID,
		"status":   status,
	})
	if err != nil {
		log.Printf("Failed to marshal order status: %v", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte(orderID),
		Value: value,
	}

	if err := p.writer.WriteMessages(context.Background(), message); err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}

	log.Printf("Published order status: OrderID=%s, Status=%s", orderID, status)
	return nil
}

func (p *Producer) Close() {
	p.writer.Close()
}

// Publish user created
func (p *Producer) PublishUserAccountCreated(userInfo *models.User) error {
	value, err := json.Marshal(map[string]string{
		"id":         userInfo.ID,
		"first_name": userInfo.FirstName,
		"email":      userInfo.Email,
	})
	if err != nil {
		log.Printf("Failed to marshal order status: %v", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte(userInfo.ID),
		Value: []byte(value),
	}

	err = p.writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}

	log.Printf("Published user created: userID=%s, email=%s", userInfo.ID, userInfo.Email)
	return nil
}
