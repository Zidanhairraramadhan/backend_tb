package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
// @description     Smart Music Profile Link Aggregator — One Link for All Your Music.
// @termsOfService  http://swagger.io/terms/

// @contact.name   MusicLink Support
// @contact.email  support@musiclink.app

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5000
// @BasePath  /

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Enter your JWT token in the format: Bearer <token>
func main() {
	// Load environment variables from .env file (ignore error in production)
	if err := godotenv.Overload(); err != nil {
		log.Println("⚠️  No .env file found — relying on system environment variables")
	}

	// Initialize Database
	config.ConnectDB()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(config.DB)
	linkRepo := repository.NewLinkRepository(config.DB)
	adminRepo := repository.NewAdminRepository(config.DB)
	clickLogRepo := repository.NewClickLogRepository(config.DB)

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(userRepo)
	userHandler := handler.NewUserHandler(userRepo)
	linkHandler := handler.NewLinkHandler(linkRepo, userRepo)
	metadataHandler := handler.NewMetadataHandler()
	adminHandler := handler.NewAdminHandler(userRepo, adminRepo)
	publicHandler := handler.NewPublicHandler(userRepo, linkRepo, clickLogRepo)
	analyticsHandler := handler.NewAnalyticsHandler(clickLogRepo)
	ogHandler := handler.NewOGHandler(userRepo)

	// Create Fiber App with custom error handler
	app := fiber.New(fiber.Config{
		AppName: "MusicLink API Platform v1.0",
		// Custom global error handler — ensures consistent JSON error responses
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"message": err.Error(),
				"code":    code,
			})
		},
	})

	// Setup all routes
	router.SetupRoutes(app, authHandler, userHandler, linkHandler, metadataHandler, adminHandler, publicHandler, analyticsHandler, ogHandler)

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// ── Graceful Shutdown ──
	// Listen for OS signals (Ctrl+C, kill) and shut down cleanly
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 MusicLink API starting on port %s...\n", port)
		log.Printf("📚 Swagger UI available at: http://localhost:%s/docs\n", port)
		log.Printf("💚 Health check at: http://localhost:%s/health\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Block until signal received
	<-quit
	log.Println("⏳ Shutting down server gracefully...")

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Fatalf("❌ Forced shutdown: %v", err)
	}

	log.Println("✅ Server shutdown complete.")
}
