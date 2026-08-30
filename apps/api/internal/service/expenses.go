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

type CreateSplitInput struct {
	ParticipantID         uuid.UUID
	AmountMinor           int64
	PercentageBasisPoints int64
}

type CreateExpenseInput struct {
	Description        string
	AmountMinor        int64
	Currency           string
	ExpenseDate        string
	PayerParticipantID uuid.UUID
	SplitType          string
	ParticipantIDs     []uuid.UUID
	Splits             []CreateSplitInput
}

type Balance struct {
	ParticipantID   uuid.UUID
	Participant     domain.Participant
	PaidAmountMinor int64
	OwedAmountMinor int64
	AmountMinor     int64
}

func (s *GroupService) CreateEqualExpense(ctx context.Context, groupID, actorID uuid.UUID, input CreateEqualExpenseInput) (*domain.Expense, error) {
	return s.CreateExpense(ctx, groupID, actorID, CreateExpenseInput{
		Description: input.Description, AmountMinor: input.AmountMinor, Currency: input.Currency, ExpenseDate: input.ExpenseDate,
		PayerParticipantID: input.PayerParticipantID, SplitType: domain.SplitTypeEqual, ParticipantIDs: input.ParticipantIDs,
	})
}

func (s *GroupService) CreateExpense(ctx context.Context, groupID, actorID uuid.UUID, input CreateExpenseInput) (*domain.Expense, error) {
	if _, err := s.requireActiveParticipant(ctx, groupID, actorID); err != nil {
		return nil, err
	}

	description := strings.TrimSpace(input.Description)
	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	expenseDate, err := parseExpenseDate(input.ExpenseDate)
	if currency == "" {
		currency = domain.DefaultCurrency
	}
	if err != nil || description == "" || input.AmountMinor <= 0 || len(currency) != 3 {
		return nil, ErrValidation
	}

	participantIDs := uniqueParticipantIDs(input.ParticipantIDs)
	if input.SplitType == domain.SplitTypeManualAmount || input.SplitType == domain.SplitTypePercentage {
		participantIDs = splitParticipantIDs(input.Splits)
	}
	if len(participantIDs) == 0 {
		return nil, ErrValidation
	}
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
	var splits []domain.ExpenseSplit
	switch input.SplitType {
	case domain.SplitTypeEqual:
		splits = equalSplits(input.AmountMinor, participantIDs)
	case domain.SplitTypeManualAmount:
		splits, err = manualAmountSplits(input.AmountMinor, participantIDs, input.Splits)
	case domain.SplitTypePercentage:
		splits, err = percentageSplits(input.AmountMinor, participantIDs, input.Splits)
	default:
		return nil, ErrValidation
	}
	if err != nil {
		return nil, ErrValidation
	}
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
		SplitType:          input.SplitType,
		Splits:             splits,
	}

	if err := s.store.CreateExpenseWithSplits(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func splitParticipantIDs(splits []CreateSplitInput) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(splits))
	for _, split := range splits {
		ids = append(ids, split.ParticipantID)
	}
	return uniqueParticipantIDs(ids)
}

func manualAmountSplits(amountMinor int64, participantIDs []uuid.UUID, inputs []CreateSplitInput) ([]domain.ExpenseSplit, error) {
	if len(participantIDs) != len(inputs) {
		return nil, ErrValidation
	}
	amounts := make(map[uuid.UUID]int64, len(inputs))
	var total int64
	for _, input := range inputs {
		if input.AmountMinor < 0 {
			return nil, ErrValidation
		}
		if _, exists := amounts[input.ParticipantID]; exists {
			return nil, ErrValidation
		}
		amounts[input.ParticipantID] = input.AmountMinor
		total += input.AmountMinor
	}
	if total != amountMinor {
		return nil, ErrValidation
	}
	splits := make([]domain.ExpenseSplit, 0, len(participantIDs))
	for _, participantID := range participantIDs {
		splits = append(splits, domain.ExpenseSplit{ParticipantID: participantID, AmountMinor: amounts[participantID]})
	}
	return splits, nil
}

func percentageSplits(amountMinor int64, participantIDs []uuid.UUID, inputs []CreateSplitInput) ([]domain.ExpenseSplit, error) {
	if len(participantIDs) != len(inputs) {
		return nil, ErrValidation
	}
	percentages := make(map[uuid.UUID]int64, len(inputs))
	var total int64
	for _, input := range inputs {
		if input.PercentageBasisPoints <= 0 {
			return nil, ErrValidation
		}
		if _, exists := percentages[input.ParticipantID]; exists {
			return nil, ErrValidation
		}
		percentages[input.ParticipantID] = input.PercentageBasisPoints
		total += input.PercentageBasisPoints
	}
	if total != 10000 {
		return nil, ErrValidation
	}
	splits := make([]domain.ExpenseSplit, 0, len(participantIDs))
	var allocated int64
	for _, participantID := range participantIDs {
		amount := amountMinor * percentages[participantID] / 10000
		allocated += amount
		splits = append(splits, domain.ExpenseSplit{ParticipantID: participantID, AmountMinor: amount, PercentageBasisPoints: percentages[participantID]})
	}
	sort.SliceStable(splits, func(i, j int) bool {
		leftRemainder := (amountMinor * splits[i].PercentageBasisPoints) % 10000
		rightRemainder := (amountMinor * splits[j].PercentageBasisPoints) % 10000
		return leftRemainder > rightRemainder
	})
	for index := int64(0); index < amountMinor-allocated; index++ {
		splits[index].AmountMinor++
	}
	sort.Slice(splits, func(i, j int) bool { return splits[i].ParticipantID.String() < splits[j].ParticipantID.String() })
	return splits, nil
}

func (s *GroupService) ListExpenses(ctx context.Context, groupID, actorID uuid.UUID) ([]domain.Expense, error) {
	if _, err := s.requireActiveParticipant(ctx, groupID, actorID); err != nil {
		return nil, err
	}
	return s.store.ListExpensesForGroup(ctx, groupID)
}

func (s *GroupService) ListBalances(ctx context.Context, groupID, actorID uuid.UUID) ([]Balance, error) {
	if _, err := s.requireActiveParticipant(ctx, groupID, actorID); err != nil {
		return nil, err
	}

	participants, err := s.store.ListParticipantsForGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	expenses, err := s.store.ListExpensesForGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	balancesByParticipantID := make(map[uuid.UUID]Balance, len(participants))
	for _, participant := range participants {
		balancesByParticipantID[participant.ID] = Balance{ParticipantID: participant.ID, Participant: participant}
	}
	for _, expense := range expenses {
		balance := balancesByParticipantID[expense.PayerParticipantID]
		balance.PaidAmountMinor += expense.AmountMinor
		balancesByParticipantID[expense.PayerParticipantID] = balance

		for _, split := range expense.Splits {
			balance := balancesByParticipantID[split.ParticipantID]
			balance.OwedAmountMinor += split.AmountMinor
			balancesByParticipantID[split.ParticipantID] = balance
		}
	}

	participantIDs := make([]uuid.UUID, 0, len(balancesByParticipantID))
	for participantID := range balancesByParticipantID {
		participantIDs = append(participantIDs, participantID)
	}
	sort.Slice(participantIDs, func(i, j int) bool {
		return participantIDs[i].String() < participantIDs[j].String()
	})

	balances := make([]Balance, 0, len(participantIDs))
	for _, participantID := range participantIDs {
		balance := balancesByParticipantID[participantID]
		balance.AmountMinor = balance.PaidAmountMinor - balance.OwedAmountMinor
		balances = append(balances, balance)
	}
	return balances, nil
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
