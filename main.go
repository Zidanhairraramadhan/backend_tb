package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"musiclink-backend/config"
	"musiclink-backend/handler"
	"musiclink-backend/repository"
	"musiclink-backend/router"

	_ "musiclink-backend/docs" // Import generated swagger docs
)

// @title           MusicLink API
// @version         1.0
// @description     This is the smart music profile link aggregator API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5000
// @BasePath  /

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer " followed by your JWT token.
func main() {
	// Load environment variables
	if err := godotenv.Overload(); err != nil {
		log.Println("⚠️ Warning: No .env file found, relying on system environment variables")
	}

	// Initialize Database
	config.ConnectDB()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(config.DB)
	linkRepo := repository.NewLinkRepository(config.DB)

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(userRepo)
	userHandler := handler.NewUserHandler(userRepo)
	linkHandler := handler.NewLinkHandler(linkRepo, userRepo)

	// Create Fiber App
	app := fiber.New(fiber.Config{
		AppName: "MusicLink API Platform v1.0",
	})

	// Setup Router
	router.SetupRoutes(app, authHandler, userHandler, linkHandler)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 MusicLink Fiber Server starting on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start Fiber server: %v", err)
	}
}
