package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter returns a general-purpose rate limiter middleware.
// Allows max 30 requests per IP per 1 minute.
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests. Please slow down and try again in a moment.",
				"code":    429,
			})
		},
	})
}

// AuthRateLimiter returns a strict rate limiter for authentication endpoints.
// Allows max 10 requests per IP per 1 minute to prevent brute force attacks.
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many authentication attempts. Please wait 1 minute before trying again.",
				"code":    429,
			})
		},
	})
}
