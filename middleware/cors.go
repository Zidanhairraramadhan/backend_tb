package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupCORS() fiber.Handler {
	allowedOrigins := os.Getenv("FRONTEND_URL")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000, http://127.0.0.1:3000"
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	})
}
