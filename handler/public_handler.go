package handler

import (
	"github.com/gofiber/fiber/v2"
	"musiclink-backend/repository"
)

// PublicHandler handles public-facing endpoints
type PublicHandler struct {
	userRepo *repository.UserRepository
}

// NewPublicHandler creates a new PublicHandler
func NewPublicHandler(userRepo *repository.UserRepository) *PublicHandler {
	return &PublicHandler{userRepo: userRepo}
}

// GetPublicProfile Fetches artist profile with their links preloaded
// @Summary      Get public profile
// @Description  Get public bio details and links of an artist by username
// @Tags         public
// @Produce      json
// @Param        username  path      string  true  "Artist username"
// @Success      200       {object}  model.User
// @Failure      404       {object}  map[string]interface{}
// @Router       /api/public/users/{username} [get]
func (h *PublicHandler) GetPublicProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.userRepo.GetPublicProfileByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User profile not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// TrackLinkClick increments the click count of a link
func (h *PublicHandler) TrackLinkClick(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.userRepo.IncrementLinkClick(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to track click",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "click tracked",
	})
}
