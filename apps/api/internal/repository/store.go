package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, user *domain.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, err
}

func (s *Store) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, err
}

func (s *Store) CreateSession(ctx context.Context, userID uuid.UUID, ttl time.Duration) (*domain.Session, error) {
	session := &domain.Session{
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(ttl),
	}
	return session, s.db.WithContext(ctx).Create(session).Error
}

func (s *Store) FindSessionUser(ctx context.Context, sessionID uuid.UUID) (*domain.User, error) {
	var session domain.Session
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND expires_at > ?", sessionID, time.Now().UTC()).
		First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session.User, nil
}

func (s *Store) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&domain.Session{}, "id = ?", sessionID).Error
}

func (s *Store) CreateGroupWithOwner(ctx context.Context, group *domain.Group, ownerID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		participant := &domain.Participant{
			GroupID: group.ID,
			UserID:  ownerID,
			Role:    domain.RoleOwner,
			Active:  true,
		}
		return tx.Create(participant).Error
	})
}

func (s *Store) ListGroupsForUser(ctx context.Context, userID uuid.UUID) ([]domain.Group, error) {
	var groups []domain.Group
	err := s.db.WithContext(ctx).
		Joins("JOIN participants ON participants.group_id = groups.id").
		Where("participants.user_id = ? AND participants.active = true", userID).
		Order("groups.created_at DESC").
		Find(&groups).Error
	return groups, err
}

func (s *Store) FindParticipant(ctx context.Context, groupID, userID uuid.UUID) (*domain.Participant, error) {
	var participant domain.Participant
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &participant, err
}

func (s *Store) FindParticipantByID(ctx context.Context, groupID, participantID uuid.UUID) (*domain.Participant, error) {
	var participant domain.Participant
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("group_id = ? AND id = ?", groupID, participantID).
		First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &participant, err
}

func (s *Store) CreateParticipant(ctx context.Context, participant *domain.Participant) error {
	return s.db.WithContext(ctx).Omit("User").Create(participant).Error
}

func (s *Store) ListParticipantsForGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Participant, error) {
	var participants []domain.Participant
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("group_id = ? AND active = true", groupID).
		Order("created_at ASC").
		Find(&participants).Error
	return participants, err
}

func (s *Store) CreateExpenseWithSplits(ctx context.Context, expense *domain.Expense) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := &domain.Expense{
			GroupID:            expense.GroupID,
			PayerParticipantID: expense.PayerParticipantID,
			Description:        expense.Description,
			AmountMinor:        expense.AmountMinor,
			Currency:           expense.Currency,
			ExpenseDate:        expense.ExpenseDate,
			SplitType:          expense.SplitType,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}

		expense.ID = record.ID
		for index := range expense.Splits {
			expense.Splits[index].ExpenseID = record.ID
		}
		if len(expense.Splits) > 0 {
			if err := tx.Omit("Participant").Create(&expense.Splits).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) ListExpensesForGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Expense, error) {
	var expenses []domain.Expense
	err := s.db.WithContext(ctx).
		Preload("PayerParticipant.User").
		Preload("Splits", func(db *gorm.DB) *gorm.DB {
			return db.Order("expense_splits.participant_id ASC")
		}).
		Preload("Splits.Participant.User").
		Where("group_id = ?", groupID).
		Order("expense_date DESC, created_at DESC").
		Find(&expenses).Error
	return expenses, err
}
