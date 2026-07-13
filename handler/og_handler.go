package handler

import (
	"fmt"
	"html"
	"strings"

	"github.com/gofiber/fiber/v2"
	"musiclink-backend/repository"
)

// OGHandler handles OpenGraph meta tag serving for social media crawlers
type OGHandler struct {
	userRepo *repository.UserRepository
}

// NewOGHandler creates a new OGHandler
func NewOGHandler(userRepo *repository.UserRepository) *OGHandler {
	return &OGHandler{userRepo: userRepo}
}

// socialCrawlerBots adalah daftar substring User-Agent yang digunakan oleh crawler media sosial
var socialCrawlerBots = []string{
	"facebookexternalhit",
	"Facebot",
	"Twitterbot",
	"LinkedInBot",
	"WhatsApp",
	"TelegramBot",
	"Slackbot",
	"Discordbot",
	"Pinterest",
	"Googlebot",
	"bingbot",
	"Applebot",
	"Pinterestbot",
}

// isSocialCrawler mendeteksi apakah request datang dari bot/crawler media sosial
func isSocialCrawler(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	for _, bot := range socialCrawlerBots {
		if strings.Contains(ua, strings.ToLower(bot)) {
			return true
		}
	}
	return false
}

// ServeProfileOG mendeteksi crawler dan menyajikan HTML dengan OpenGraph meta tags
// Jika bukan crawler, redirect ke halaman SPA frontend
// @Summary      Serve OpenGraph meta tags
// @Description  Returns minimal HTML with OG tags for social media crawlers, redirects normal users to SPA
// @Tags         public
// @Produce      html
// @Param        username  path      string  true  "Artist username"
// @Success      200       {string}  string  "HTML with OG meta tags"
// @Success      302       {string}  string  "Redirect to SPA"
// @Failure      404       {object}  map[string]interface{}
// @Router       /p/{username} [get]
func (h *OGHandler) ServeProfileOG(c *fiber.Ctx) error {
	username := c.Params("username")
	userAgent := c.Get("User-Agent")

	// Jika bukan crawler, redirect ke SPA frontend
	if !isSocialCrawler(userAgent) {
		// Redirect ke halaman profil publik SPA
		frontendURL := c.Query("frontend", "")
		if frontendURL == "" {
			// Default: asumsikan frontend berada di origin yang sama
			frontendURL = fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
		}
		return c.Redirect(fmt.Sprintf("%s/#/public/%s", frontendURL, username), fiber.StatusFound)
	}

	// Crawler terdeteksi — ambil data user dari database
	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Profile not found")
	}

	// Siapkan data untuk OG tags
	displayName := html.EscapeString(user.Username)
	if user.Name != "" {
		displayName = html.EscapeString(user.Name)
	}

	bio := html.EscapeString(user.Bio)
	if bio == "" {
		bio = fmt.Sprintf("Dengarkan musik dari %s di MusicLink", displayName)
	}

	avatarURL := user.AvatarURL
	if avatarURL == "" {
		// Placeholder jika tidak ada avatar
		avatarURL = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=1DB954&color=fff&size=400&bold=true", username)
	}

	profileURL := fmt.Sprintf("%s://%s/p/%s", c.Protocol(), c.Hostname(), username)

	genre := user.Genre
	if genre == "" {
		genre = "Musician"
	}

	// Render HTML minimal dengan OG tags
	ogHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <!-- Primary Meta Tags -->
    <title>%s — MusicLink</title>
    <meta name="title" content="%s — MusicLink">
    <meta name="description" content="%s">

    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="profile">
    <meta property="og:url" content="%s">
    <meta property="og:title" content="%s — MusicLink">
    <meta property="og:description" content="%s">
    <meta property="og:image" content="%s">
    <meta property="og:site_name" content="MusicLink">

    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image">
    <meta property="twitter:url" content="%s">
    <meta property="twitter:title" content="%s — MusicLink">
    <meta property="twitter:description" content="%s">
    <meta property="twitter:image" content="%s">

    <!-- Additional -->
    <meta property="og:locale" content="id_ID">
    <meta name="robots" content="index, follow">
    <link rel="canonical" href="%s">
</head>
<body>
    <h1>%s</h1>
    <p>%s</p>
    <p>%s</p>
    <a href="%s">Kunjungi profil di MusicLink</a>
</body>
</html>`,
		displayName, displayName, bio,
		profileURL, displayName, bio, avatarURL,
		profileURL, displayName, bio, avatarURL,
		profileURL,
		displayName, bio, genre, profileURL,
	)

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(ogHTML)
}
