package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
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
