package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/repository"
)

const currentUserKey = "currentUser"

func (s *Server) requireUser(c *fiber.Ctx) error {
	cookie := c.Cookies(s.cfg.SessionCookieName)
	if cookie == "" {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHENTICATED", "Please sign in.", nil)
	}

	sessionID, err := uuid.Parse(cookie)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHENTICATED", "Please sign in.", nil)
	}

	user, err := s.auth.CurrentUser(c.UserContext(), sessionID)
	if errors.Is(err, repository.ErrNotFound) {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHENTICATED", "Please sign in.", nil)
	}
	if err != nil {
		return err
	}

	c.Locals(currentUserKey, user)
	return c.Next()
}

func currentUser(c *fiber.Ctx) *domain.User {
	user, _ := c.Locals(currentUserKey).(*domain.User)
	return user
}
