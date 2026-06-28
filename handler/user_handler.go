package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"musiclink-backend/model"
	"musiclink-backend/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetProfile Fetches current user profile
// @Summary      Get user profile
// @Description  Get profile of currently authenticated user
// @Tags         profile
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  model.User
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User profile not found"})
	}

	return c.JSON(user)
}

// UpdateProfile Updates current user profile
// @Summary      Update user profile
// @Description  Update profile fields for the currently authenticated user
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        profile  body      model.ProfileRequest  true  "Profile update fields"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	req := new(model.ProfileRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON request"})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Update fields
	if strings.TrimSpace(req.Name) != "" {
		user.Name = req.Name
		// Re-derive initials
		names := strings.Fields(req.Name)
		initial := ""
		if len(names) > 0 {
			initial += string(names[0][0])
			if len(names) > 1 {
				initial += string(names[1][0])
			}
		}
		user.AvatarInitial = strings.ToUpper(initial)
	}

	user.Bio = req.Bio
	user.Genre = req.Genre
	user.Country = req.Country

	if err := h.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update profile"})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}
