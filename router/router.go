package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"musiclink-backend/handler"
	"musiclink-backend/middleware"
)

func SetupRoutes(app *fiber.App, authH *handler.AuthHandler, userH *handler.UserHandler, linkH *handler.LinkHandler) {
	// Logger & CORS Middleware
	app.Use(logger.New())
	app.Use(middleware.SetupCORS())

	// Swagger UI Route
	app.Get("/docs/*", swagger.HandlerDefault)

	// Public Unprotected Auth Routes
	app.Post("/register", authH.Register)
	app.Post("/login", authH.Login)

	// Public Profile Route
	app.Get("/public/:username", linkH.GetPublicProfile)
	app.Post("/api/clicks/:id", linkH.IncrementClickCounts)

	// Protected Routes Group
	api := app.Group("/api", middleware.JWTProtected())

	// Auth Settings
	api.Put("/change-password", authH.ChangePassword)

	// Profile
	api.Get("/profile", userH.GetProfile)
	api.Put("/profile", userH.UpdateProfile)

	// Links CRUD
	api.Get("/links", linkH.GetMyLinks)
	api.Get("/links/:id", linkH.GetLinkByID)
	api.Post("/links", linkH.CreateLink)
	api.Put("/links/:id", linkH.UpdateLink)
	api.Delete("/links/:id", linkH.DeleteLink)
}
