package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/your-org/splitr/apps/api/internal/domain"
)

func TestGroupServiceCreateEqualExpense(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()
	friendID := uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	friendParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID:      ownerParticipantID,
		GroupID: groupID,
		UserID:  ownerID,
		Role:    domain.RoleOwner,
		Active:  true,
		User:    domain.User{ID: ownerID, Email: "owner@example.com"},
	}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{
		ID:      friendParticipantID,
		GroupID: groupID,
		UserID:  friendID,
		Role:    domain.RoleParticipant,
		Active:  true,
		User:    domain.User{ID: friendID, Email: "friend@example.com"},
	}

	service := NewGroupService(store)
	expense, err := service.CreateEqualExpense(ctx, groupID, ownerID, CreateEqualExpenseInput{
		Description:        "Noodles",
		AmountMinor:        10001,
		Currency:           "thb",
		ExpenseDate:        "2026-08-14",
		PayerParticipantID: ownerParticipantID,
		ParticipantIDs:     []uuid.UUID{friendParticipantID, ownerParticipantID},
	})

	require.NoError(t, err)
	require.Equal(t, "Noodles", expense.Description)
	require.Equal(t, int64(10001), expense.AmountMinor)
	require.Equal(t, domain.DefaultCurrency, expense.Currency)
	require.Equal(t, domain.SplitTypeEqual, expense.SplitType)
	require.Equal(t, ownerParticipantID, expense.PayerParticipantID)
	require.Len(t, expense.Splits, 2)
	require.Equal(t, ownerParticipantID, expense.Splits[0].ParticipantID)
	require.Equal(t, int64(5001), expense.Splits[0].AmountMinor)
	require.Equal(t, friendParticipantID, expense.Splits[1].ParticipantID)
	require.Equal(t, int64(5000), expense.Splits[1].AmountMinor)
}

func TestGroupServiceCreateEqualExpenseDefaultsCurrencyToTHB(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000031")

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID:      ownerParticipantID,
		GroupID: groupID,
		UserID:  ownerID,
		Role:    domain.RoleOwner,
		Active:  true,
		User:    domain.User{ID: ownerID, Email: "owner@example.com"},
	}

	service := NewGroupService(store)
	expense, err := service.CreateEqualExpense(ctx, groupID, ownerID, CreateEqualExpenseInput{
		Description:        "Coffee",
		AmountMinor:        4500,
		ExpenseDate:        "2026-08-14",
		PayerParticipantID: ownerParticipantID,
		ParticipantIDs:     []uuid.UUID{ownerParticipantID},
	})

	require.NoError(t, err)
	require.Equal(t, domain.DefaultCurrency, expense.Currency)
}

func TestGroupServiceCreateEqualExpenseRejectsPayerOutsideSplit(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()
	friendID := uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000011")
	friendParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000012")

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID:      ownerParticipantID,
		GroupID: groupID,
		UserID:  ownerID,
		Role:    domain.RoleOwner,
		Active:  true,
	}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{
		ID:      friendParticipantID,
		GroupID: groupID,
		UserID:  friendID,
		Role:    domain.RoleParticipant,
		Active:  true,
	}

	service := NewGroupService(store)
	_, err := service.CreateEqualExpense(ctx, groupID, ownerID, CreateEqualExpenseInput{
		Description:        "Taxi",
		AmountMinor:        4000,
		Currency:           "THB",
		ExpenseDate:        "2026-08-14",
		PayerParticipantID: ownerParticipantID,
		ParticipantIDs:     []uuid.UUID{friendParticipantID},
	})

	require.ErrorIs(t, err, ErrValidation)
}

func TestGroupServiceCreateEqualExpenseRejectsParticipantOutsideGroup(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	otherGroupID := uuid.New()
	ownerID := uuid.New()
	friendID := uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000021")
	otherParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000022")

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID:      ownerParticipantID,
		GroupID: groupID,
		UserID:  ownerID,
		Role:    domain.RoleOwner,
		Active:  true,
	}
	store.participants[participantKey(otherGroupID, friendID)] = &domain.Participant{
		ID:      otherParticipantID,
		GroupID: otherGroupID,
		UserID:  friendID,
		Role:    domain.RoleParticipant,
		Active:  true,
	}

	service := NewGroupService(store)
	_, err := service.CreateEqualExpense(ctx, groupID, ownerID, CreateEqualExpenseInput{
		Description:        "Coffee",
		AmountMinor:        3000,
		Currency:           "THB",
		ExpenseDate:        "2026-08-14",
		PayerParticipantID: ownerParticipantID,
		ParticipantIDs:     []uuid.UUID{ownerParticipantID, otherParticipantID},
	})

	require.ErrorIs(t, err, ErrValidation)
}
