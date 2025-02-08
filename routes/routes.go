package routes

import (
	"ecommerce/handlers"
	"ecommerce/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(payment *handlers.Payment, user *handlers.User) *mux.Router {
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/users", user.CreateUser).Methods("POST")
	router.HandleFunc("/users/login", handlers.Login).Methods("POST")
	router.HandleFunc("/users/{id}", middleware.AuthMiddleware(handlers.GetUser)).Methods("GET")

	// // Product routes
	router.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/products", handlers.GetProducts).Methods("GET")

	// // Cart routes
	router.HandleFunc("/cart", middleware.AuthMiddleware(handlers.AddToCart)).Methods("POST")
	router.HandleFunc("/cart/{userID}", middleware.AuthMiddleware(handlers.GetCartItems)).Methods("GET")

	// // Order routes
	router.HandleFunc("/orders", middleware.AuthMiddleware(handlers.CreateOrder)).Methods("POST")
	// router.HandleFunc("/orders", middleware.AuthMiddleware(handlers.GetOrders)).Methods("GET")

	// Payment routes
	router.HandleFunc("/payments/initiate", middleware.AuthMiddleware(handlers.InitiatePayment)).Methods("POST")
	router.HandleFunc("/payments/webhook", payment.PaymentWebhook).Methods("POST")

	// Payment gateway (This is equivalent to the payment gateway notifying your server about the payment status)
	router.HandleFunc("/paymentGateway/notify", handlers.Notify).Methods("POST")
	return router
}
