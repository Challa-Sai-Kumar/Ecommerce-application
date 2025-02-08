package main

import (
	"ecommerce/database"
	"ecommerce/handlers"
	"ecommerce/kafka"
	"ecommerce/notifications"
	"ecommerce/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Database
	db := database.ConnectDB(config.Database.DataSourceName, config.Database.DriverName)
	defer db.Close()

	emailConfig := notifications.NewEmailConfig(config.Email.SMTPPort, config.Email.SMTPHost, config.Email.Username, config.Email.Password, config.Email.FromAddress)

	// Initialize Kafka producer
	orderProducer := kafka.NewProducer(config.Kafka.BrokerList, config.Kafka.Topics["order_status"])
	defer orderProducer.Close()

	userProducer := kafka.NewProducer(config.Kafka.BrokerList, config.Kafka.Topics["user_notifications"])
	defer userProducer.Close()

	// Start Kafka consumers in a separate Goroutine
	go func() {
		err := kafka.StartOrderConsumer(emailConfig, config.Kafka.BrokerList, config.Kafka.Topics["order_status"], config.Kafka.ConsumerGroups["order_status_group"])
		if err != nil {
			log.Printf("Consumer error for topic 'order_status': %v", err)
		}
	}()

	go func() {
		err := kafka.StartUserConsumer(emailConfig, config.Kafka.BrokerList, config.Kafka.Topics["user_notifications"], config.Kafka.ConsumerGroups["user_notifications_group"])
		if err != nil {
			log.Printf("Consumer error for topic 'user_notifications': %v", err)
		}
	}()

	payment := handlers.NewPayment(orderProducer)
	// handler := handlers.NewHandle(payment)
	user := handlers.NewUser(userProducer)

	// Set up Routes
	router := routes.SetupRoutes(payment, user)

	port := config.Server.Port
	// Start Server
	log.Printf("Starting server on port %d...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	// Initialize Email Service
	// emailService = &EmailService{
	// 	SMTPHost: "smtp.gmail.com",
	// 	SMTPPort: 587,
	// 	Username: "your-email@gmail.com",
	// 	Password: "your-app-password",
	// }

}
