package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// SetupHelmet returns the Helmet middleware for HTTP security headers.
// Adds headers like X-Content-Type-Options, X-Frame-Options, HSTS, etc.
func SetupHelmet() fiber.Handler {
	return helmet.New(helmet.Config{
		// Prevent MIME type sniffing
		ContentTypeNosniff: "nosniff",
		// Prevent clickjacking
		XFrameOptions: "DENY",
		// Disable legacy XSS filter (modern approach — let CSP handle it)
		XSSProtection: "0",
		// Force HTTPS for 1 year (effective in production)
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    false,
		// Control referrer information
		ReferrerPolicy: "strict-origin-when-cross-origin",
		// Restrict browser features
		PermissionPolicy: "camera=(), microphone=(), geolocation=()",
		// Cross-Origin policies
		CrossOriginEmbedderPolicy: "unsafe-none", // Allow embedding for public profiles
		CrossOriginOpenerPolicy:   "same-origin-allow-popups",
		CrossOriginResourcePolicy: "cross-origin",
	})
}
