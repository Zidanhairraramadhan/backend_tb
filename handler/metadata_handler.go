package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MetadataHandler handles Universal link metadata fetching
type MetadataHandler struct{}

// NewMetadataHandler creates a new MetadataHandler
func NewMetadataHandler() *MetadataHandler {
	return &MetadataHandler{}
}

// oEmbedResponse maps the typical fields returned by oEmbed APIs
type oEmbedResponse struct {
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// GetLinkMetadata fetches metadata via public oEmbed APIs based on the URL domain
// @Summary      Get link metadata
// @Description  Fetch title and thumbnail image of a link using public oEmbed APIs (no auth required). Supports Spotify, YouTube, SoundCloud, and TikTok.
// @Tags         metadata
// @Produce      json
// @Param        url  query     string  true  "URL (e.g. https://open.spotify.com/track/...)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      502  {object}  map[string]interface{}
// @Router       /api/link/metadata [get]
func (h *MetadataHandler) GetLinkMetadata(c *fiber.Ctx) error {
	linkURL := strings.TrimSpace(c.Query("url"))

	if linkURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Query parameter 'url' is required",
		})
	}

	var oembedURL string
	var platform string

	// Determine platform and oEmbed URL
	if strings.Contains(linkURL, "spotify.com") {
		oembedURL = fmt.Sprintf("https://open.spotify.com/oembed?url=%s", linkURL)
		platform = "Spotify"
	} else if strings.Contains(linkURL, "youtube.com") || strings.Contains(linkURL, "youtu.be") {
		oembedURL = fmt.Sprintf("https://www.youtube.com/oembed?url=%s&format=json", linkURL)
		platform = "YouTube"
	} else if strings.Contains(linkURL, "soundcloud.com") {
		oembedURL = fmt.Sprintf("https://soundcloud.com/oembed?format=json&url=%s", linkURL)
		platform = "SoundCloud"
	} else if strings.Contains(linkURL, "tiktok.com") {
		oembedURL = fmt.Sprintf("https://www.tiktok.com/oembed?url=%s", linkURL)
		platform = "TikTok"
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Platform tidak didukung untuk auto-fill",
		})
	}

	// Perform HTTP GET to the oEmbed API
	resp, err := http.Get(oembedURL) //nolint:gosec
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to reach %s oEmbed API", platform),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": fmt.Sprintf("%s oEmbed API returned status %d. Make sure the URL is public and valid.", platform, resp.StatusCode),
		})
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Failed to read oEmbed response",
		})
	}

	var oembed oEmbedResponse
	if err := json.Unmarshal(body, &oembed); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Failed to parse oEmbed response",
		})
	}

	return c.JSON(fiber.Map{
		"title": oembed.Title,
		"image": oembed.ThumbnailURL,
	})
}
