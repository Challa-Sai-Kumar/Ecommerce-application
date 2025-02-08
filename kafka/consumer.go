package kafka

import (
	"context"
	"ecommerce/database/dao"
	"ecommerce/models"
	"ecommerce/notifications"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderStatusEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type UserInfo struct {
	UserID        string
	UserFirstName string
	UserEmail     string
}

// StartKafkaConsumer initializes a Kafka consumer and processes messages.
func StartOrderConsumer(emailConfig *notifications.EmailConfig, broker []string, topic, groupID string) error {
	reader := newKafkaReader(broker, topic, groupID)
	defer reader.Close()

	log.Printf("Starting Kafka consumer for topic: %s, groupID: %s", topic, groupID)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		message, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message from topic %s: %v", topic, err)
			continue
		}

		var event OrderStatusEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		orderDetails, err := dao.GetOrderDetails(event.OrderID)
		if err != nil {
			log.Printf("Unable to fetch order by ID %s: %v", event.OrderID, err)
			continue
		}

		if err := emailConfig.NotifyOrderStatus(orderDetails, event.Status); err != nil {
			log.Printf("Failed to send order status email: %v", err)
		}
	}
}

func StartUserConsumer(emailConfig *notifications.EmailConfig, broker []string, topic, groupID string) error {
	reader := newKafkaReader(broker, topic, groupID)
	defer reader.Close()

	log.Printf("Starting Kafka consumer for topic: %s, groupID: %s", topic, groupID)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		message, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message from topic %s: %v", topic, err)
			continue
		}

		var user models.User
		if err := json.Unmarshal(message.Value, &user); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		if err := emailConfig.NotifyUserCreated(&user); err != nil {
			log.Printf("Failed to send user created email: %v", err)
		}
	}
}

func newKafkaReader(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})
}
