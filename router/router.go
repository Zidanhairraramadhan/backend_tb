package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"musiclink-backend/handler"
	"musiclink-backend/middleware"
)

func SetupRoutes(app *fiber.App, authH *handler.AuthHandler, userH *handler.UserHandler, linkH *handler.LinkHandler, metadataH *handler.MetadataHandler, adminH *handler.AdminHandler, publicH *handler.PublicHandler, analyticsH *handler.AnalyticsHandler, ogH *handler.OGHandler) {
	// ── Global Middleware ──
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(middleware.SetupHelmet())
	app.Use(middleware.SetupCORS())
	app.Use(middleware.RateLimiter()) // General: 30 req/min per IP

	// ── Swagger UI ──
	app.Get("/docs/*", swagger.HandlerDefault)

	// ── Health Check (for uptime monitors & deploy platforms) ──
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "MusicLink API",
		})
	})

	// ── Public Unprotected Auth Routes (with strict rate limiting) ──
	app.Post("/register", middleware.AuthRateLimiter(), authH.Register)
	app.Post("/login", middleware.AuthRateLimiter(), authH.Login)

	// ── Public Profile Routes ──
	app.Get("/public/:username", linkH.GetPublicProfile)
	app.Post("/api/links/:id/click", linkH.IncrementClickCounts)

	// ── Public Profile Endpoint (Preloaded Links) ──
	app.Get("/api/public/u/:username", publicH.GetPublicProfile)
	app.Post("/api/public/link/:id/click", publicH.TrackLinkClick)

	// ── Public Universal oEmbed Metadata (no auth required) ──
	app.Get("/api/link/metadata", metadataH.GetLinkMetadata)

	// ── Dynamic OpenGraph Meta Tags for Social Media Crawlers ──
	app.Get("/p/:username", ogH.ServeProfileOG)

	// ── Protected Routes Group ──
	api := app.Group("/api", middleware.JWTProtected())

	// ── Admin Area ──
	api.Get("/admin/users", middleware.RequireAdmin(), adminH.GetAllUsers)
	api.Get("/admin/stats", middleware.RequireAdmin(), adminH.GetGlobalStats)
	api.Put("/admin/users/:id/password", middleware.RequireAdmin(), adminH.ForceChangePassword)
	api.Delete("/admin/users/:id", middleware.RequireAdmin(), adminH.DeleteUser)

	// ── Auth Settings ──
	api.Put("/change-password", authH.ChangePassword)

	// ── Profile ──
	api.Get("/profile", userH.GetProfile)
	api.Put("/profile", userH.UpdateProfile)

	// ── Links CRUD ──
	api.Get("/links", linkH.GetMyLinks)
	api.Get("/links/:id", linkH.GetLinkByID)
	api.Post("/links", linkH.CreateLink)
	api.Put("/links/reorder", linkH.ReorderLinks) // Harus sebelum /links/:id agar tidak bentrok
	api.Put("/links/:id", linkH.UpdateLink)
	api.Delete("/links/:id", linkH.DeleteLink)

	// ── Analytics ──
	api.Get("/analytics/daily", analyticsH.GetDailyClicks)
	api.Get("/analytics/monthly", analyticsH.GetMonthlyClicks)
	api.Get("/analytics/sources", analyticsH.GetTrafficSources)
	api.Get("/analytics/summary", analyticsH.GetAnalyticsSummary)
}
