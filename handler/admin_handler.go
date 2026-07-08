package handler

import (
	"github.com/gofiber/fiber/v2"
	"musiclink-backend/repository"
)

// AdminHandler handles admin-only endpoints
type AdminHandler struct {
	userRepo *repository.UserRepository
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

// GetAllUsers returns all users with their associated links
// @Summary      Get all users (Admin)
// @Description  Returns a list of every registered user along with their music links embedded in each user object
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   model.User
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/admin/users [get]
func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userRepo.GetAllWithLinks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"total": len(users),
		"users": users,
	})
}
