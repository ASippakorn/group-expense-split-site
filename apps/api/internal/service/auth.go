package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/repository"
	"github.com/your-org/splitr/apps/api/internal/security"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrValidation         = errors.New("validation failed")
	ErrForbidden          = errors.New("forbidden")
	ErrUserNotFound       = errors.New("user not found")
	ErrParticipantExists  = errors.New("participant already exists")
)

type AuthService struct {
	store          *repository.Store
	sessionTTL     time.Duration
	passwordPepper string
}

func NewAuthService(store *repository.Store, sessionTTL time.Duration, passwordPepper string) *AuthService {
	return &AuthService{store: store, sessionTTL: sessionTTL, passwordPepper: passwordPepper}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, *domain.Session, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return nil, nil, ErrValidation
	}
	if err := security.ValidatePassword(password); err != nil {
		return nil, nil, ErrValidation
	}

	passwordHash, err := security.HashPassword(password, s.passwordPepper)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{Email: normalizedEmail, PasswordHash: passwordHash}
	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, nil, err
	}

	session, err := s.store.CreateSession(ctx, user.ID, s.sessionTTL)
	return user, session, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, *domain.Session, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	user, err := s.store.FindUserByEmail(ctx, normalizedEmail)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, nil, err
	}
	if !security.VerifyPassword(user.PasswordHash, password, s.passwordPepper) {
		return nil, nil, ErrInvalidCredentials
	}

	session, err := s.store.CreateSession(ctx, user.ID, s.sessionTTL)
	return user, session, err
}

func (s *AuthService) CurrentUser(ctx context.Context, sessionID uuid.UUID) (*domain.User, error) {
	return s.store.FindSessionUser(ctx, sessionID)
}

func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return s.store.DeleteSession(ctx, sessionID)
}

func normalizeEmail(value string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if _, err := mail.ParseAddress(trimmed); err != nil {
		return "", err
	}
	return trimmed, nil
}
