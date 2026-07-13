package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"musiclink-backend/handler"
	"musiclink-backend/middleware"
)

func SetupRoutes(app *fiber.App, authH *handler.AuthHandler, userH *handler.UserHandler, linkH *handler.LinkHandler, metadataH *handler.MetadataHandler, adminH *handler.AdminHandler, publicH *handler.PublicHandler, analyticsH *handler.AnalyticsHandler, ogH *handler.OGHandler) {
	// Logger & CORS Middleware
	app.Use(logger.New())
	app.Use(middleware.SetupCORS())

	// Swagger UI Route
	app.Get("/docs/*", swagger.HandlerDefault)

	// Public Unprotected Auth Routes
	app.Post("/register", authH.Register)
	app.Post("/login", authH.Login)

	// Public Profile Route (Legacy/Alternative)
	app.Get("/public/:username", linkH.GetPublicProfile)
	app.Post("/api/links/:id/click", linkH.IncrementClickCounts)

	// Public Profile Endpoint (Preloaded Links)
	app.Get("/api/public/u/:username", publicH.GetPublicProfile)
	app.Post("/api/public/link/:id/click", publicH.TrackLinkClick)

	// Public Universal oEmbed Metadata (no auth required)
	app.Get("/api/link/metadata", metadataH.GetLinkMetadata)

	// Dynamic OpenGraph Meta Tags for Social Media Crawlers
	app.Get("/p/:username", ogH.ServeProfileOG)

	// Protected Routes Group
	api := app.Group("/api", middleware.JWTProtected())

	// Admin Area
	api.Get("/admin/users", middleware.RequireAdmin(), adminH.GetAllUsers)
	api.Get("/admin/stats", middleware.RequireAdmin(), adminH.GetGlobalStats)

	// Auth Settings
	api.Put("/change-password", authH.ChangePassword)

	// Profile
	api.Get("/profile", userH.GetProfile)
	api.Put("/profile", userH.UpdateProfile)

	// Links CRUD
	api.Get("/links", linkH.GetMyLinks)
	api.Get("/links/:id", linkH.GetLinkByID)
	api.Post("/links", linkH.CreateLink)
	api.Put("/links/reorder", linkH.ReorderLinks) // Harus sebelum /links/:id agar tidak bentrok
	api.Put("/links/:id", linkH.UpdateLink)
	api.Delete("/links/:id", linkH.DeleteLink)

	// Analytics
	api.Get("/analytics/daily", analyticsH.GetDailyClicks)
	api.Get("/analytics/monthly", analyticsH.GetMonthlyClicks)
	api.Get("/analytics/sources", analyticsH.GetTrafficSources)
	api.Get("/analytics/summary", analyticsH.GetAnalyticsSummary)
}

