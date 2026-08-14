package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/repository"
)

type GroupRepository interface {
	CreateGroupWithOwner(ctx context.Context, group *domain.Group, ownerID uuid.UUID) error
	ListGroupsForUser(ctx context.Context, userID uuid.UUID) ([]domain.Group, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindParticipant(ctx context.Context, groupID, userID uuid.UUID) (*domain.Participant, error)
	CreateParticipant(ctx context.Context, participant *domain.Participant) error
	ListParticipantsForGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Participant, error)
}

type GroupService struct {
	store GroupRepository
}

func NewGroupService(store GroupRepository) *GroupService {
	return &GroupService{store: store}
}

func (s *GroupService) CreateGroup(ctx context.Context, ownerID uuid.UUID, name, defaultCurrency, description string) (*domain.Group, error) {
	name = strings.TrimSpace(name)
	defaultCurrency = strings.ToUpper(strings.TrimSpace(defaultCurrency))
	description = strings.TrimSpace(description)

	if name == "" || len(defaultCurrency) != 3 {
		return nil, ErrValidation
	}

	group := &domain.Group{
		Name:            name,
		DefaultCurrency: defaultCurrency,
		Description:     description,
		OwnerID:         ownerID,
	}

	return group, s.store.CreateGroupWithOwner(ctx, group, ownerID)
}

func (s *GroupService) ListGroups(ctx context.Context, userID uuid.UUID) ([]domain.Group, error) {
	return s.store.ListGroupsForUser(ctx, userID)
}

func (s *GroupService) AddParticipantByEmail(ctx context.Context, groupID, actorID uuid.UUID, email string) (*domain.Participant, error) {
	actor, err := s.store.FindParticipant(ctx, groupID, actorID)
	if errors.Is(err, repository.ErrNotFound) || err == nil && (!actor.Active || actor.Role != domain.RoleOwner) {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}

	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return nil, ErrValidation
	}

	user, err := s.store.FindUserByEmail(ctx, normalizedEmail)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if _, err := s.store.FindParticipant(ctx, groupID, user.ID); err == nil {
		return nil, ErrParticipantExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	participant := &domain.Participant{
		GroupID: groupID,
		UserID:  user.ID,
		User:    *user,
		Role:    domain.RoleParticipant,
		Active:  true,
	}
	if err := s.store.CreateParticipant(ctx, participant); err != nil {
		return nil, err
	}
	return participant, nil
}

func (s *GroupService) ListParticipants(ctx context.Context, groupID, actorID uuid.UUID) ([]domain.Participant, error) {
	actor, err := s.store.FindParticipant(ctx, groupID, actorID)
	if errors.Is(err, repository.ErrNotFound) || err == nil && !actor.Active {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}
	return s.store.ListParticipantsForGroup(ctx, groupID)
}
