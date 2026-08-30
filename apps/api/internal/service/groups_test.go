package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/repository"
)

func TestGroupServiceAddParticipantByEmail(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()
	targetID := uuid.New()

	store := newFakeGroupStore()
	store.usersByEmail["friend@example.com"] = &domain.User{ID: targetID, Email: "friend@example.com"}
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID:      uuid.New(),
		GroupID: groupID,
		UserID:  ownerID,
		Role:    domain.RoleOwner,
		Active:  true,
		User:    domain.User{ID: ownerID, Email: "owner@example.com"},
	}

	service := NewGroupService(store)
	participant, err := service.AddParticipantByEmail(ctx, groupID, ownerID, " FRIEND@example.com ")

	require.NoError(t, err)
	require.Equal(t, targetID, participant.UserID)
	require.Equal(t, domain.RoleParticipant, participant.Role)
	require.True(t, participant.Active)
	require.Equal(t, "friend@example.com", participant.User.Email)
}

func TestGroupServiceAddParticipantByEmailRejectsRegularParticipant(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	regularID := uuid.New()

	store := newFakeGroupStore()
	store.usersByEmail["friend@example.com"] = &domain.User{ID: uuid.New(), Email: "friend@example.com"}
	store.participants[participantKey(groupID, regularID)] = &domain.Participant{
		ID:      uuid.New(),
		GroupID: groupID,
		UserID:  regularID,
		Role:    domain.RoleParticipant,
		Active:  true,
	}

	service := NewGroupService(store)
	_, err := service.AddParticipantByEmail(ctx, groupID, regularID, "friend@example.com")

	require.ErrorIs(t, err, ErrForbidden)
}

func TestGroupServiceAddParticipantByEmailRejectsUnknownEmail(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID:      uuid.New(),
		GroupID: groupID,
		UserID:  ownerID,
		Role:    domain.RoleOwner,
		Active:  true,
	}

	service := NewGroupService(store)
	_, err := service.AddParticipantByEmail(ctx, groupID, ownerID, "missing@example.com")

	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestGroupServiceAddParticipantByEmailRejectsDuplicateParticipant(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()
	targetID := uuid.New()

	store := newFakeGroupStore()
	store.usersByEmail["friend@example.com"] = &domain.User{ID: targetID, Email: "friend@example.com"}
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID:      uuid.New(),
		GroupID: groupID,
		UserID:  ownerID,
		Role:    domain.RoleOwner,
		Active:  true,
	}
	store.participants[participantKey(groupID, targetID)] = &domain.Participant{
		ID:      uuid.New(),
		GroupID: groupID,
		UserID:  targetID,
		Role:    domain.RoleParticipant,
		Active:  true,
	}

	service := NewGroupService(store)
	_, err := service.AddParticipantByEmail(ctx, groupID, ownerID, "friend@example.com")

	require.ErrorIs(t, err, ErrParticipantExists)
}

func TestGroupServiceListParticipantsRequiresMembership(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	strangerID := uuid.New()

	service := NewGroupService(newFakeGroupStore())
	_, err := service.ListParticipants(ctx, groupID, strangerID)

	require.ErrorIs(t, err, ErrForbidden)
}

type fakeGroupStore struct {
	usersByEmail map[string]*domain.User
	participants map[string]*domain.Participant
	expenses     []domain.Expense
	settlements  []domain.Settlement
	tags         map[uuid.UUID]*domain.Tag
}

func newFakeGroupStore() *fakeGroupStore {
	return &fakeGroupStore{
		usersByEmail: make(map[string]*domain.User),
		participants: make(map[string]*domain.Participant),
		tags:         make(map[uuid.UUID]*domain.Tag),
	}
}

func (s *fakeGroupStore) FindTagByID(_ context.Context, groupID, tagID uuid.UUID) (*domain.Tag, error) {
	tag, ok := s.tags[tagID]
	if !ok || tag.GroupID != groupID {
		return nil, repository.ErrNotFound
	}
	return tag, nil
}

func (s *fakeGroupStore) FindTagByName(_ context.Context, groupID uuid.UUID, name string) (*domain.Tag, error) {
	for _, tag := range s.tags {
		if tag.GroupID == groupID && tag.Name == name {
			return tag, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (s *fakeGroupStore) CreateTag(_ context.Context, tag *domain.Tag) error {
	tag.ID = uuid.New()
	s.tags[tag.ID] = tag
	return nil
}

func (s *fakeGroupStore) ListTagsForGroup(_ context.Context, groupID uuid.UUID) ([]domain.Tag, error) {
	tags := make([]domain.Tag, 0)
	for _, tag := range s.tags {
		if tag.GroupID == groupID {
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}

func participantKey(groupID, userID uuid.UUID) string {
	return groupID.String() + ":" + userID.String()
}

func (s *fakeGroupStore) CreateGroupWithOwner(context.Context, *domain.Group, uuid.UUID) error {
	return nil
}

func (s *fakeGroupStore) ListGroupsForUser(context.Context, uuid.UUID) ([]domain.Group, error) {
	return nil, nil
}

func (s *fakeGroupStore) FindUserByEmail(_ context.Context, email string) (*domain.User, error) {
	user, ok := s.usersByEmail[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (s *fakeGroupStore) FindParticipant(_ context.Context, groupID, userID uuid.UUID) (*domain.Participant, error) {
	participant, ok := s.participants[participantKey(groupID, userID)]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return participant, nil
}

func (s *fakeGroupStore) FindParticipantByID(_ context.Context, groupID, participantID uuid.UUID) (*domain.Participant, error) {
	for _, participant := range s.participants {
		if participant.GroupID == groupID && participant.ID == participantID {
			return participant, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (s *fakeGroupStore) CreateParticipant(_ context.Context, participant *domain.Participant) error {
	key := participantKey(participant.GroupID, participant.UserID)
	if _, ok := s.participants[key]; ok {
		return errors.New("duplicate participant")
	}
	participant.ID = uuid.New()
	participant.Active = true
	s.participants[key] = participant
	return nil
}

func (s *fakeGroupStore) ListParticipantsForGroup(_ context.Context, groupID uuid.UUID) ([]domain.Participant, error) {
	participants := make([]domain.Participant, 0)
	for _, participant := range s.participants {
		if participant.GroupID == groupID {
			participants = append(participants, *participant)
		}
	}
	return participants, nil
}

func (s *fakeGroupStore) CreateExpenseWithSplits(_ context.Context, expense *domain.Expense) error {
	expense.ID = uuid.New()
	for index := range expense.Splits {
		expense.Splits[index].ID = uuid.New()
		expense.Splits[index].ExpenseID = expense.ID
	}
	s.expenses = append(s.expenses, *expense)
	return nil
}

func (s *fakeGroupStore) ListExpensesForGroup(_ context.Context, groupID uuid.UUID) ([]domain.Expense, error) {
	expenses := make([]domain.Expense, 0)
	for _, expense := range s.expenses {
		if expense.GroupID == groupID {
			expenses = append(expenses, expense)
		}
	}
	return expenses, nil
}

func (s *fakeGroupStore) CreateSettlement(_ context.Context, settlement *domain.Settlement) error {
	settlement.ID = uuid.New()
	s.settlements = append(s.settlements, *settlement)
	return nil
}

func (s *fakeGroupStore) ListSettlementsForGroup(_ context.Context, groupID uuid.UUID) ([]domain.Settlement, error) {
	settlements := make([]domain.Settlement, 0)
	for _, settlement := range s.settlements {
		if settlement.GroupID == groupID {
			settlements = append(settlements, settlement)
		}
	}
	return settlements, nil
}

func (s *fakeGroupStore) DeleteSettlement(_ context.Context, groupID, settlementID uuid.UUID) error {
	for index, settlement := range s.settlements {
		if settlement.GroupID == groupID && settlement.ID == settlementID {
			s.settlements = append(s.settlements[:index], s.settlements[index+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}
