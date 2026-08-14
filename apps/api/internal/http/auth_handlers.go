package http

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/service"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (s *Server) register(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Request body is invalid.", nil)
	}

	user, session, err := s.auth.Register(c.UserContext(), req.Email, req.Password)
	if errors.Is(err, service.ErrValidation) {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Email or password is invalid.", nil)
	}
	if err != nil {
		return err
	}

	s.setSessionCookie(c, session.ID, session.ExpiresAt)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"user": toUserResponse(user)})
}

func (s *Server) login(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Request body is invalid.", nil)
	}

	user, session, err := s.auth.Login(c.UserContext(), req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		return writeError(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Email or password is incorrect.", nil)
	}
	if err != nil {
		return err
	}

	s.setSessionCookie(c, session.ID, session.ExpiresAt)
	return c.JSON(fiber.Map{"user": toUserResponse(user)})
}

func (s *Server) logout(c *fiber.Ctx) error {
	sessionID, _ := uuid.Parse(c.Cookies(s.cfg.SessionCookieName))
	if err := s.auth.Logout(c.UserContext(), sessionID); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     s.cfg.SessionCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"user": toUserResponse(currentUser(c))})
}

func (s *Server) setSessionCookie(c *fiber.Ctx, sessionID uuid.UUID, expiresAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     s.cfg.SessionCookieName,
		Value:    sessionID.String(),
		Expires:  expiresAt,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})
}

func toUserResponse(user *domain.User) userResponse {
	return userResponse{ID: user.ID.String(), Email: user.Email}
}
