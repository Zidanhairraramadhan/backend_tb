package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"musiclink-backend/repository"
)

// AdminHandler handles admin-only endpoints
type AdminHandler struct {
	userRepo  *repository.UserRepository
	adminRepo *repository.AdminRepository
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(userRepo *repository.UserRepository, adminRepo *repository.AdminRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:  userRepo,
		adminRepo: adminRepo,
	}
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

// GetGlobalStats returns global statistics
// @Summary      Get global stats (Admin)
// @Description  Returns total users, links, clicks, and 5 recent users
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/admin/stats [get]
func (h *AdminHandler) GetGlobalStats(c *fiber.Ctx) error {
	stats, err := h.adminRepo.GetGlobalStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch global stats",
		})
	}

	type RecentUserResp struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	var recentUsers []RecentUserResp
	for _, u := range stats.RecentUsers {
		recentUsers = append(recentUsers, RecentUserResp{
			ID:        u.ID,
			Username:  u.Username,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		})
	}
	if recentUsers == nil {
		recentUsers = []RecentUserResp{}
	}

	return c.JSON(fiber.Map{
		"total_users":  stats.TotalUsers,
		"total_links":  stats.TotalLinks,
		"total_clicks": stats.TotalClicks,
		"recent_users": recentUsers,
	})
}

// ForceChangePassword allows admin to reset any user's password without needing the current password
// @Summary      Force change user password (Admin)
// @Description  Admin can forcefully reset any user's password by user ID. Does not require the current password.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Target User ID (UUID)"
// @Param        request  body      object  true  "New password body"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/admin/users/{id}/password [put]
func (h *AdminHandler) ForceChangePassword(c *fiber.Ctx) error {
	targetUserID := c.Params("id")
	if strings.TrimSpace(targetUserID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User ID is required"})
	}

	type Req struct {
		NewPassword string `json:"new_password"`
	}

	req := new(Req)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON request"})
	}

	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "New password must be at least 6 characters"})
	}

	user, err := h.userRepo.GetByID(targetUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to encrypt password"})
	}

	user.Password = string(hashedPassword)
	if err := h.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update password"})
	}

	return c.JSON(fiber.Map{
		"message":  "Password for user @" + user.Username + " has been reset successfully",
		"username": user.Username,
	})
}

// DeleteUser allows admin to delete any user by ID
// @Summary      Delete a user (Admin)
// @Description  Admin can delete any user by user ID.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Target User ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	targetUserID := c.Params("id")
	if strings.TrimSpace(targetUserID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User ID is required"})
	}

	user, err := h.userRepo.GetByID(targetUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Jangan izinkan admin menghapus dirinya sendiri
	currentUserID := c.Locals("user_id").(string)
	if user.ID == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Admin cannot delete their own account"})
	}

	if err := h.userRepo.Delete(targetUserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{
		"message": "User @" + user.Username + " has been deleted successfully",
	})
}
