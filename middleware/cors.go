package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS configures the CORS middleware.
// FRONTEND_URL can be a comma-separated list of allowed origins, e.g.:
// FRONTEND_URL=https://musiclink.vercel.app,https://www.musiclink.app
// Defaults to common localhost ports for development.
func SetupCORS() fiber.Handler {
	allowedOrigins := os.Getenv("FRONTEND_URL")
	if allowedOrigins == "" {
		// Default: allow common local development ports
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173"
	}

	// Normalize: trim spaces around each origin
	origins := strings.Split(allowedOrigins, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}
	allowedOrigins = strings.Join(origins, ",")

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
		MaxAge:           300, // Cache preflight for 5 minutes
	})
}
