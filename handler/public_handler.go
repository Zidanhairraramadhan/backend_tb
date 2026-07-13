package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"musiclink-backend/model"
	"musiclink-backend/repository"
)

// PublicHandler handles public-facing endpoints
type PublicHandler struct {
	userRepo     *repository.UserRepository
	linkRepo     *repository.LinkRepository
	clickLogRepo *repository.ClickLogRepository
}

// NewPublicHandler creates a new PublicHandler
func NewPublicHandler(userRepo *repository.UserRepository, linkRepo *repository.LinkRepository, clickLogRepo *repository.ClickLogRepository) *PublicHandler {
	return &PublicHandler{
		userRepo:     userRepo,
		linkRepo:     linkRepo,
		clickLogRepo: clickLogRepo,
	}
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

// TrackLinkClick menyimpan log klik terperinci dan menambah counter klik
// @Summary      Track link click
// @Description  Logs a detailed click event (referrer, device) and increments click counter
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Link ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/public/link/{id}/click [post]
func (h *PublicHandler) TrackLinkClick(c *fiber.Ctx) error {
	id := c.Params("id")

	// 1. Increment counter di tabel links (backward compatible)
	err := h.userRepo.IncrementLinkClick(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to track click",
		})
	}

	// 2. Cari link untuk mendapatkan user_id (pemilik link)
	link, err := h.linkRepo.GetByID(id)
	if err != nil {
		// Jika link tidak ditemukan, tetap return sukses (counter sudah bertambah)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "click tracked",
		})
	}

	// 3. Deteksi referrer dari body JSON atau HTTP Referer header
	var req model.ClickTrackRequest
	_ = c.BodyParser(&req) // Opsional, boleh kosong

	referrer := req.Referrer
	if referrer == "" {
		referrer = detectReferrer(c.Get("Referer"))
	}
	if referrer == "" {
		referrer = "Direct"
	}

	// 4. Deteksi device dari User-Agent header
	device := detectDevice(c.Get("User-Agent"))

	// 5. Simpan log klik terperinci
	clickLog := &model.ClickLog{
		LinkID:   id,
		UserID:   link.UserID,
		Referrer: referrer,
		Device:   device,
	}

	// Simpan secara asinkron-like (jangan block response jika gagal)
	_ = h.clickLogRepo.Create(clickLog)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "click tracked",
	})
}

// detectReferrer mengekstrak nama platform dari URL Referer HTTP header
func detectReferrer(referer string) string {
	if referer == "" {
		return ""
	}
	r := strings.ToLower(referer)

	switch {
	case strings.Contains(r, "instagram.com") || strings.Contains(r, "instagram"):
		return "Instagram"
	case strings.Contains(r, "twitter.com") || strings.Contains(r, "x.com"):
		return "Twitter/X"
	case strings.Contains(r, "tiktok.com"):
		return "TikTok"
	case strings.Contains(r, "facebook.com") || strings.Contains(r, "fb.com"):
		return "Facebook"
	case strings.Contains(r, "youtube.com") || strings.Contains(r, "youtu.be"):
		return "YouTube"
	case strings.Contains(r, "whatsapp.com") || strings.Contains(r, "wa.me"):
		return "WhatsApp"
	case strings.Contains(r, "t.me") || strings.Contains(r, "telegram"):
		return "Telegram"
	case strings.Contains(r, "google.com") || strings.Contains(r, "google"):
		return "Google"
	case strings.Contains(r, "discord.com") || strings.Contains(r, "discord.gg"):
		return "Discord"
	default:
		return "Other"
	}
}

// detectDevice mendeteksi tipe perangkat (Mobile/Desktop) dari User-Agent
func detectDevice(userAgent string) string {
	if userAgent == "" {
		return "Unknown"
	}
	ua := strings.ToLower(userAgent)
	mobileKeywords := []string{"mobile", "android", "iphone", "ipad", "ipod", "blackberry", "opera mini", "opera mobi"}
	for _, kw := range mobileKeywords {
		if strings.Contains(ua, kw) {
			return "Mobile"
		}
	}
	return "Desktop"
}
