package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"wati-clone-backend/internal/infrastructure/db"
	"wati-clone-backend/internal/api/handlers"
	"wati-clone-backend/internal/api/middleware"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	// Connect DB
	db.Connect()

	// Init Router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Routes
	api := r.Group("/api")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/signup", handlers.Signup)
			auth.POST("/login", handlers.Login)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Tenant
			protected.POST("/tenants", handlers.CreateTenant)
			
			// Messaging
			protected.POST("/messages/send", handlers.SendMessage)
			protected.GET("/messages", handlers.GetMessages)
			
			// Contacts
			protected.GET("/contacts", handlers.GetContacts)
			protected.POST("/contacts", handlers.CreateContact)
		}
		
		// Webhooks (Public)
		api.POST("/webhooks/whatsapp", handlers.WhatsAppWebhook)
		api.GET("/webhooks/whatsapp", handlers.WhatsAppWebhookVerify)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
