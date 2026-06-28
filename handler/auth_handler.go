package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"musiclink-backend/middleware"
	"musiclink-backend/model"
	"musiclink-backend/repository"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// Register Handles user registration
// @Summary      Register user
// @Description  Create a new user account (role defaults to 'user' unless username contains 'admin')
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      model.RegisterRequest  true  "Register request body"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(model.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON request"})
	}

	// Validation
	if strings.TrimSpace(req.Username) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Username is required"})
	}
	if len(req.Username) < 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Username must be at least 4 characters"})
	}
	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Password must be at least 6 characters"})
	}
	if strings.TrimSpace(req.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Full name is required"})
	}

	// Check if username duplicate
	existing, _ := h.userRepo.GetByUsername(req.Username)
	if existing != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Username is already taken"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to encrypt password"})
	}

	// Determine Role
	role := "user"
	if strings.Contains(strings.ToLower(req.Username), "admin") {
		role = "admin"
	}

	// Derive avatar initial
	names := strings.Fields(req.Name)
	initial := ""
	if len(names) > 0 {
		initial += string(names[0][0])
		if len(names) > 1 {
			initial += string(names[1][0])
		}
	}
	initial = strings.ToUpper(initial)

	user := &model.User{
		Username:      req.Username,
		Password:      string(hashedPassword),
		Role:          role,
		Name:          req.Name,
		AvatarInitial: initial,
		Bio:           "Hello, I am using MusicLink! 🎧",
		Country:       "United States",
		Genre:         "Pop",
		Verified:      role == "admin", // Admin is auto-verified for cool look
	}

	if err := h.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login Authenticates user and returns JWT
// @Summary      Login user
// @Description  Authenticate user and return JWT bearer token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      model.LoginRequest  true  "Login request body"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(model.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON request"})
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Username and password are required"})
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	// Generate Token
	token, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

// ChangePassword Changes the password of currently logged in user
// @Summary      Change password
// @Description  Update password for current logged-in user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      model.ChangePasswordRequest  true  "Change password body"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/change-password [put]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	req := new(model.ChangePasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON request"})
	}

	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "New password must be at least 6 characters"})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Incorrect current password"})
	}

	// Hash new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to encrypt new password"})
	}

	user.Password = string(hashedNewPassword)
	if err := h.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update password"})
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}
