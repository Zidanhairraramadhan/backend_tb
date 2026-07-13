package handler

import (
	"github.com/gofiber/fiber/v2"
	"musiclink-backend/repository"
)

// AnalyticsHandler handles analytics-related API endpoints
type AnalyticsHandler struct {
	clickLogRepo *repository.ClickLogRepository
}

// NewAnalyticsHandler creates a new AnalyticsHandler
func NewAnalyticsHandler(clickLogRepo *repository.ClickLogRepository) *AnalyticsHandler {
	return &AnalyticsHandler{clickLogRepo: clickLogRepo}
}

// GetDailyClicks returns click counts grouped by day (last 7 days)
// @Summary      Get daily click analytics
// @Description  Get click counts per day for the last 7 days
// @Tags         analytics
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/analytics/daily [get]
func (h *AnalyticsHandler) GetDailyClicks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	stats, err := h.clickLogRepo.GetDailyClicks(userID, 7)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch daily analytics",
		})
	}

	return c.JSON(fiber.Map{
		"daily": stats,
	})
}

// GetMonthlyClicks returns click counts grouped by month (last 12 months)
// @Summary      Get monthly click analytics
// @Description  Get click counts per month for the last 12 months
// @Tags         analytics
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/analytics/monthly [get]
func (h *AnalyticsHandler) GetMonthlyClicks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	stats, err := h.clickLogRepo.GetMonthlyClicks(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch monthly analytics",
		})
	}

	return c.JSON(fiber.Map{
		"monthly": stats,
	})
}

// GetTrafficSources returns click counts grouped by referrer source
// @Summary      Get traffic source analytics
// @Description  Get click distribution by traffic source (referrer)
// @Tags         analytics
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/analytics/sources [get]
func (h *AnalyticsHandler) GetTrafficSources(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	sources, err := h.clickLogRepo.GetTrafficSources(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch traffic sources",
		})
	}

	return c.JSON(fiber.Map{
		"sources": sources,
	})
}

// GetAnalyticsSummary returns a comprehensive analytics summary
// @Summary      Get analytics summary
// @Description  Get total clicks, weekly growth, and most active day
// @Tags         analytics
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/analytics/summary [get]
func (h *AnalyticsHandler) GetAnalyticsSummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	totalClicks, _ := h.clickLogRepo.GetTotalClicksByUser(userID)
	thisWeek, _ := h.clickLogRepo.GetClicksThisWeek(userID)
	lastWeek, _ := h.clickLogRepo.GetClicksLastWeek(userID)

	// Hitung growth percentage
	var growthPercentage float64
	if lastWeek > 0 {
		growthPercentage = float64(thisWeek-lastWeek) / float64(lastWeek) * 100
	} else if thisWeek > 0 {
		growthPercentage = 100
	}

	// Ambil daily stats untuk menentukan hari paling aktif
	dailyStats, _ := h.clickLogRepo.GetDailyClicks(userID, 30)
	mostActiveDay := "N/A"
	maxDayClicks := 0
	for _, d := range dailyStats {
		if d.Clicks > maxDayClicks {
			maxDayClicks = d.Clicks
			mostActiveDay = d.Date
		}
	}

	return c.JSON(fiber.Map{
		"total_clicks":      totalClicks,
		"clicks_this_week":  thisWeek,
		"clicks_last_week":  lastWeek,
		"growth_percentage": growthPercentage,
		"most_active_day":   mostActiveDay,
	})
}
