package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"musiclink-backend/model"
	"musiclink-backend/repository"
)

type LinkHandler struct {
	linkRepo *repository.LinkRepository
	userRepo *repository.UserRepository
}

func NewLinkHandler(linkRepo *repository.LinkRepository, userRepo *repository.UserRepository) *LinkHandler {
	return &LinkHandler{
		linkRepo: linkRepo,
		userRepo: userRepo,
	}
}

// GetMyLinks Fetches all links for the current user
// @Summary      Get user links
// @Description  Get all links belonging to the authenticated user
// @Tags         links
// @Produce      json
// @Security     BearerAuth
// @Success      200      {array}   model.Link
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/links [get]
func (h *LinkHandler) GetMyLinks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	links, err := h.linkRepo.GetAllByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch links"})
	}

	return c.JSON(links)
}

// GetLinkByID Fetches link detail
// @Summary      Get link detail
// @Description  Get link by ID, only accessible by link owner or admin
// @Tags         links
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Link ID"
// @Success      200      {object}  model.Link
// @Failure      401      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Router       /api/links/{id} [get]
func (h *LinkHandler) GetLinkByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid link ID"})
	}

	link, err := h.linkRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Link not found"})
	}

	// Verify ownership (Admin bypassed)
	if link.UserID != userID && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: You do not own this link"})
	}

	return c.JSON(link)
}

// CreateLink Adds a new music link
// @Summary      Create link
// @Description  Add a new link under the authenticated user
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        link  body      model.CreateLinkRequest  true  "Link creation body"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/links [post]
func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	req := new(model.CreateLinkRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON request"})
	}

	// Validation
	if strings.TrimSpace(req.Platform) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Platform is required"})
	}
	if strings.TrimSpace(req.Title) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Title is required"})
	}
	if strings.TrimSpace(req.URL) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "URL is required"})
	}
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "URL must start with http:// or https://"})
	}

	link := &model.Link{
		UserID:   userID,
		Platform: strings.ToLower(req.Platform),
		Title:    req.Title,
		URL:      req.URL,
		Active:   true,
		Clicks:   0,
	}

	if err := h.linkRepo.Create(link); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add link"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Link added successfully",
		"link":    link,
	})
}

// UpdateLink Modifies a link by ID
// @Summary      Update link
// @Description  Update fields of an existing link by ID
// @Tags         links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                      true  "Link ID"
// @Param        link  body      model.UpdateLinkRequest  true  "Link update fields"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Router       /api/links/{id} [put]
func (h *LinkHandler) UpdateLink(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid link ID"})
	}

	req := new(model.UpdateLinkRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON request"})
	}

	link, err := h.linkRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Link not found"})
	}

	// Verify ownership (Admin bypassed)
	if link.UserID != userID && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: You do not own this link"})
	}

	// Apply updates
	if req.Platform != "" {
		link.Platform = strings.ToLower(req.Platform)
	}
	if req.Title != "" {
		link.Title = req.Title
	}
	if req.URL != "" {
		if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "URL must start with http:// or https://"})
		}
		link.URL = req.URL
	}
	if req.Active != nil {
		link.Active = *req.Active
	}

	if err := h.linkRepo.Update(link); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update link"})
	}

	return c.JSON(fiber.Map{
		"message": "Link updated successfully",
		"link":    link,
	})
}

// DeleteLink Removes a link by ID
// @Summary      Delete link
// @Description  Remove link by ID
// @Tags         links
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Link ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Router       /api/links/{id} [delete]
func (h *LinkHandler) DeleteLink(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid link ID"})
	}

	link, err := h.linkRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Link not found"})
	}

	// Verify ownership (Admin bypassed)
	if link.UserID != userID && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: You do not own this link"})
	}

	if err := h.linkRepo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete link"})
	}

	return c.JSON(fiber.Map{
		"message": "Link deleted successfully",
	})
}

// GetPublicProfile Fetches artist and active links for public page
// @Summary      Get public profile
// @Description  Get public bio details and active links of an artist by username
// @Tags         public
// @Produce      json
// @Param        username  path      string  true  "Artist username"
// @Success      200       {object}  map[string]interface{}
// @Failure      404       {object}  map[string]interface{}
// @Router       /public/{username} [get]
func (h *LinkHandler) GetPublicProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Artist profile not found"})
	}

	links, err := h.linkRepo.GetActiveByUserID(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch links"})
	}

	return c.JSON(fiber.Map{
		"user":  user,
		"links": links,
	})
}

// IncrementClickCounts Registers link click redirect count
func (h *LinkHandler) IncrementClickCounts(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid link ID"})
	}

	err = h.linkRepo.IncrementClicks(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to increment clicks"})
	}

	return c.JSON(fiber.Map{
		"message": "Click registered successfully",
	})
}
