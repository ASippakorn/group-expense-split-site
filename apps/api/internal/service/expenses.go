package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/repository"
)

type CreateEqualExpenseInput struct {
	Description        string
	AmountMinor        int64
	Currency           string
	ExpenseDate        string
	PayerParticipantID uuid.UUID
	ParticipantIDs     []uuid.UUID
}

func (s *GroupService) CreateEqualExpense(ctx context.Context, groupID, actorID uuid.UUID, input CreateEqualExpenseInput) (*domain.Expense, error) {
	if _, err := s.requireActiveParticipant(ctx, groupID, actorID); err != nil {
		return nil, err
	}

	description := strings.TrimSpace(input.Description)
	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	expenseDate, err := parseExpenseDate(input.ExpenseDate)
	if currency == "" {
		currency = domain.DefaultCurrency
	}
	if err != nil || description == "" || input.AmountMinor <= 0 || len(currency) != 3 || len(input.ParticipantIDs) == 0 {
		return nil, ErrValidation
	}

	participantIDs := uniqueParticipantIDs(input.ParticipantIDs)
	selected := make(map[uuid.UUID]domain.Participant, len(participantIDs))
	for _, participantID := range participantIDs {
		participant, err := s.store.FindParticipantByID(ctx, groupID, participantID)
		if errors.Is(err, repository.ErrNotFound) || err == nil && !participant.Active {
			return nil, ErrValidation
		}
		if err != nil {
			return nil, err
		}
		selected[participantID] = *participant
	}

	payer, ok := selected[input.PayerParticipantID]
	if !ok || !payer.Active {
		return nil, ErrValidation
	}

	sort.Slice(participantIDs, func(i, j int) bool {
		return participantIDs[i].String() < participantIDs[j].String()
	})
	splits := equalSplits(input.AmountMinor, participantIDs)
	for index := range splits {
		splits[index].Participant = selected[splits[index].ParticipantID]
	}

	expense := &domain.Expense{
		GroupID:            groupID,
		PayerParticipantID: input.PayerParticipantID,
		PayerParticipant:   payer,
		Description:        description,
		AmountMinor:        input.AmountMinor,
		Currency:           currency,
		ExpenseDate:        expenseDate,
		SplitType:          domain.SplitTypeEqual,
		Splits:             splits,
	}

	if err := s.store.CreateExpenseWithSplits(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *GroupService) ListExpenses(ctx context.Context, groupID, actorID uuid.UUID) ([]domain.Expense, error) {
	if _, err := s.requireActiveParticipant(ctx, groupID, actorID); err != nil {
		return nil, err
	}
	return s.store.ListExpensesForGroup(ctx, groupID)
}

func (s *GroupService) requireActiveParticipant(ctx context.Context, groupID, userID uuid.UUID) (*domain.Participant, error) {
	actor, err := s.store.FindParticipant(ctx, groupID, userID)
	if errors.Is(err, repository.ErrNotFound) || err == nil && !actor.Active {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}
	return actor, nil
}

func parseExpenseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
}

func uniqueParticipantIDs(participantIDs []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]bool, len(participantIDs))
	unique := make([]uuid.UUID, 0, len(participantIDs))
	for _, participantID := range participantIDs {
		if participantID == uuid.Nil || seen[participantID] {
			continue
		}
		seen[participantID] = true
		unique = append(unique, participantID)
	}
	return unique
}

func equalSplits(amountMinor int64, participantIDs []uuid.UUID) []domain.ExpenseSplit {
	count := int64(len(participantIDs))
	base := amountMinor / count
	remainder := amountMinor % count
	splits := make([]domain.ExpenseSplit, 0, len(participantIDs))

	for index, participantID := range participantIDs {
		amount := base
		if int64(index) < remainder {
			amount++
		}
		splits = append(splits, domain.ExpenseSplit{
			ParticipantID: participantID,
			AmountMinor:   amount,
		})
	}
	return splits
}
