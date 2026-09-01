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

func TestGroupServiceCreateManualAmountExpenseAndUpdatesBalances(t *testing.T) {
	ctx := context.Background()
	groupID, ownerID, friendID := uuid.New(), uuid.New(), uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000101")
	friendParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000102")
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: ownerParticipantID, GroupID: groupID, UserID: ownerID, Active: true}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{ID: friendParticipantID, GroupID: groupID, UserID: friendID, Active: true}

	expense, err := NewGroupService(store).CreateExpense(ctx, groupID, ownerID, CreateExpenseInput{
		Description: "Dinner", AmountMinor: 10000, Currency: "THB", ExpenseDate: "2026-08-14",
		PayerParticipantID: ownerParticipantID, SplitType: domain.SplitTypeManualAmount,
		Splits: []CreateSplitInput{{ParticipantID: ownerParticipantID, AmountMinor: 2500}, {ParticipantID: friendParticipantID, AmountMinor: 7500}},
	})

	require.NoError(t, err)
	require.Equal(t, domain.SplitTypeManualAmount, expense.SplitType)
	require.Equal(t, int64(2500), expense.Splits[0].AmountMinor)
	require.Equal(t, int64(7500), expense.Splits[1].AmountMinor)
	balances, err := NewGroupService(store).ListBalances(ctx, groupID, ownerID)
	require.NoError(t, err)
	require.Equal(t, int64(7500), balances[0].AmountMinor)
	require.Equal(t, int64(-7500), balances[1].AmountMinor)
}

func TestGroupServiceRejectsManualAmountsThatDoNotEqualExpense(t *testing.T) {
	ctx := context.Background()
	groupID, ownerID := uuid.New(), uuid.New()
	participantID := uuid.New()
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: participantID, GroupID: groupID, UserID: ownerID, Active: true}

	_, err := NewGroupService(store).CreateExpense(ctx, groupID, ownerID, CreateExpenseInput{
		Description: "Dinner", AmountMinor: 10000, Currency: "THB", ExpenseDate: "2026-08-14", PayerParticipantID: participantID,
		SplitType: domain.SplitTypeManualAmount, Splits: []CreateSplitInput{{ParticipantID: participantID, AmountMinor: 9999}},
	})

	require.ErrorIs(t, err, ErrValidation)
}

func TestGroupServiceCreatesPercentageExpenseWithDeterministicRounding(t *testing.T) {
	ctx := context.Background()
	groupID, ownerID, friendID := uuid.New(), uuid.New(), uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000201")
	friendParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000202")
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: ownerParticipantID, GroupID: groupID, UserID: ownerID, Active: true}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{ID: friendParticipantID, GroupID: groupID, UserID: friendID, Active: true}

	expense, err := NewGroupService(store).CreateExpense(ctx, groupID, ownerID, CreateExpenseInput{
		Description: "Taxi", AmountMinor: 10001, Currency: "THB", ExpenseDate: "2026-08-14", PayerParticipantID: ownerParticipantID,
		SplitType: domain.SplitTypePercentage,
		Splits:    []CreateSplitInput{{ParticipantID: friendParticipantID, PercentageBasisPoints: 5000}, {ParticipantID: ownerParticipantID, PercentageBasisPoints: 5000}},
	})

	require.NoError(t, err)
	require.Equal(t, domain.SplitTypePercentage, expense.SplitType)
	require.Equal(t, int64(5001), expense.Splits[0].AmountMinor)
	require.Equal(t, int64(5000), expense.Splits[1].AmountMinor)
	require.Equal(t, int64(5000), expense.Splits[0].PercentageBasisPoints)
}

func TestGroupServiceCreatesTagExpenseForTagMembersAndUpdatesBalances(t *testing.T) {
	ctx := context.Background()
	groupID, ownerID, friendID, excludedID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000401")
	friendParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000402")
	excludedParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000403")
	tagID := uuid.New()
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: ownerParticipantID, GroupID: groupID, UserID: ownerID, Active: true}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{ID: friendParticipantID, GroupID: groupID, UserID: friendID, Active: true}
	store.participants[participantKey(groupID, excludedID)] = &domain.Participant{ID: excludedParticipantID, GroupID: groupID, UserID: excludedID, Active: true}
	store.tags[tagID] = &domain.Tag{ID: tagID, GroupID: groupID, Name: "Alcohol", Participants: []domain.Participant{
		*store.participants[participantKey(groupID, ownerID)], *store.participants[participantKey(groupID, friendID)],
	}}

	expense, err := NewGroupService(store).CreateExpense(ctx, groupID, ownerID, CreateExpenseInput{
		Description: "Wine", AmountMinor: 9000, Currency: "THB", ExpenseDate: "2026-08-14",
		PayerParticipantID: ownerParticipantID, SplitType: domain.SplitTypeTag, TagID: tagID,
	})

	require.NoError(t, err)
	require.Equal(t, domain.SplitTypeTag, expense.SplitType)
	require.Len(t, expense.Splits, 2)
	require.Equal(t, int64(4500), expense.Splits[0].AmountMinor)
	require.Equal(t, int64(4500), expense.Splits[1].AmountMinor)
	balances, err := NewGroupService(store).ListBalances(ctx, groupID, ownerID)
	require.NoError(t, err)
	require.Equal(t, int64(4500), balances[0].AmountMinor)
	require.Equal(t, int64(-4500), balances[1].AmountMinor)
	require.Equal(t, int64(0), balances[2].AmountMinor)
}

func TestGroupServiceCreatesGroupScopedTagWithSelectedMembers(t *testing.T) {
	ctx := context.Background()
	groupID, ownerID, friendID := uuid.New(), uuid.New(), uuid.New()
	ownerParticipantID, friendParticipantID := uuid.New(), uuid.New()
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: ownerParticipantID, GroupID: groupID, UserID: ownerID, Active: true}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{ID: friendParticipantID, GroupID: groupID, UserID: friendID, Active: true}

	tag, err := NewGroupService(store).CreateTag(ctx, groupID, ownerID, " Alcohol ", []uuid.UUID{friendParticipantID})

	require.NoError(t, err)
	require.Equal(t, groupID, tag.GroupID)
	require.Equal(t, "Alcohol", tag.Name)
	require.Equal(t, []domain.Participant{*store.participants[participantKey(groupID, friendID)]}, tag.Participants)
	listed, err := NewGroupService(store).ListTags(ctx, groupID, ownerID)
	require.NoError(t, err)
	require.Len(t, listed, 1)
	_, err = NewGroupService(store).CreateTag(ctx, groupID, ownerID, "Alcohol", []uuid.UUID{friendParticipantID})
	require.ErrorIs(t, err, ErrValidation)
}

func TestGroupServiceRejectsTagFromAnotherGroup(t *testing.T) {
	ctx := context.Background()
	groupID, otherGroupID, ownerID := uuid.New(), uuid.New(), uuid.New()
	participantID, tagID := uuid.New(), uuid.New()
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: participantID, GroupID: groupID, UserID: ownerID, Active: true}
	store.tags[tagID] = &domain.Tag{ID: tagID, GroupID: otherGroupID, Name: "Other group"}

	_, err := NewGroupService(store).CreateExpense(ctx, groupID, ownerID, CreateExpenseInput{
		Description: "Wine", AmountMinor: 9000, Currency: "THB", ExpenseDate: "2026-08-14",
		PayerParticipantID: participantID, SplitType: domain.SplitTypeTag, TagID: tagID,
	})

	require.ErrorIs(t, err, ErrValidation)
}

func TestGroupServiceRejectsPercentagesThatDoNotEqualOneHundredPercent(t *testing.T) {
	ctx := context.Background()
	groupID, ownerID := uuid.New(), uuid.New()
	participantID := uuid.New()
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: participantID, GroupID: groupID, UserID: ownerID, Active: true}

	_, err := NewGroupService(store).CreateExpense(ctx, groupID, ownerID, CreateExpenseInput{
		Description: "Taxi", AmountMinor: 10000, Currency: "THB", ExpenseDate: "2026-08-14", PayerParticipantID: participantID,
		SplitType: domain.SplitTypePercentage, Splits: []CreateSplitInput{{ParticipantID: participantID, PercentageBasisPoints: 9999}},
	})

	require.ErrorIs(t, err, ErrValidation)
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

func TestGroupServiceListBalancesCalculatesPaidAndOwedAmounts(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()
	friendID := uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	friendParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{
		ID: ownerParticipantID, GroupID: groupID, UserID: ownerID, Role: domain.RoleOwner, Active: true,
	}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{
		ID: friendParticipantID, GroupID: groupID, UserID: friendID, Role: domain.RoleParticipant, Active: true,
	}
	store.expenses = []domain.Expense{{
		GroupID: groupID, PayerParticipantID: ownerParticipantID, AmountMinor: 10001, Currency: domain.DefaultCurrency,
		Splits: []domain.ExpenseSplit{
			{ParticipantID: ownerParticipantID, AmountMinor: 5001},
			{ParticipantID: friendParticipantID, AmountMinor: 5000},
		},
	}}

	balances, err := NewGroupService(store).ListBalances(ctx, groupID, ownerID)

	require.NoError(t, err)
	require.Equal(t, []Balance{
		{ParticipantID: ownerParticipantID, Participant: *store.participants[participantKey(groupID, ownerID)], PaidAmountMinor: 10001, OwedAmountMinor: 5001, AmountMinor: 5000},
		{ParticipantID: friendParticipantID, Participant: *store.participants[participantKey(groupID, friendID)], PaidAmountMinor: 0, OwedAmountMinor: 5000, AmountMinor: -5000},
	}, balances)
}

func TestGroupServiceListBalancesAccumulatesMultipleExpensesWithUnevenSplits(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	ownerID := uuid.New()
	friendOneID := uuid.New()
	friendTwoID := uuid.New()
	ownerParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	friendOneParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	friendTwoParticipantID := uuid.MustParse("00000000-0000-0000-0000-000000000003")

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: ownerParticipantID, GroupID: groupID, UserID: ownerID, Active: true}
	store.participants[participantKey(groupID, friendOneID)] = &domain.Participant{ID: friendOneParticipantID, GroupID: groupID, UserID: friendOneID, Active: true}
	store.participants[participantKey(groupID, friendTwoID)] = &domain.Participant{ID: friendTwoParticipantID, GroupID: groupID, UserID: friendTwoID, Active: true}
	store.expenses = []domain.Expense{
		{
			GroupID: groupID, PayerParticipantID: ownerParticipantID, AmountMinor: 10001,
			Splits: []domain.ExpenseSplit{
				{ParticipantID: ownerParticipantID, AmountMinor: 3334},
				{ParticipantID: friendOneParticipantID, AmountMinor: 3334},
				{ParticipantID: friendTwoParticipantID, AmountMinor: 3333},
			},
		},
		{
			GroupID: groupID, PayerParticipantID: friendOneParticipantID, AmountMinor: 9000,
			Splits: []domain.ExpenseSplit{
				{ParticipantID: ownerParticipantID, AmountMinor: 3000},
				{ParticipantID: friendOneParticipantID, AmountMinor: 3000},
				{ParticipantID: friendTwoParticipantID, AmountMinor: 3000},
			},
		},
	}

	balances, err := NewGroupService(store).ListBalances(ctx, groupID, ownerID)

	require.NoError(t, err)
	require.Equal(t, []Balance{
		{ParticipantID: ownerParticipantID, Participant: *store.participants[participantKey(groupID, ownerID)], PaidAmountMinor: 10001, OwedAmountMinor: 6334, AmountMinor: 3667},
		{ParticipantID: friendOneParticipantID, Participant: *store.participants[participantKey(groupID, friendOneID)], PaidAmountMinor: 9000, OwedAmountMinor: 6334, AmountMinor: 2666},
		{ParticipantID: friendTwoParticipantID, Participant: *store.participants[participantKey(groupID, friendTwoID)], PaidAmountMinor: 0, OwedAmountMinor: 6333, AmountMinor: -6333},
	}, balances)
}

func TestGroupServiceRecordsAndDeletesSettlementThatUpdatesBalances(t *testing.T) {
	ctx := context.Background()
	groupID, ownerID, friendID := uuid.New(), uuid.New(), uuid.New()
	ownerParticipantID, friendParticipantID := uuid.New(), uuid.New()
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, ownerID)] = &domain.Participant{ID: ownerParticipantID, GroupID: groupID, UserID: ownerID, Active: true}
	store.participants[participantKey(groupID, friendID)] = &domain.Participant{ID: friendParticipantID, GroupID: groupID, UserID: friendID, Active: true}
	store.expenses = []domain.Expense{{GroupID: groupID, PayerParticipantID: ownerParticipantID, AmountMinor: 10000, Splits: []domain.ExpenseSplit{{ParticipantID: ownerParticipantID, AmountMinor: 5000}, {ParticipantID: friendParticipantID, AmountMinor: 5000}}}}
	service := NewGroupService(store)

	settlement, err := service.CreateSettlement(ctx, groupID, friendID, CreateSettlementInput{PayerParticipantID: friendParticipantID, ReceiverParticipantID: ownerParticipantID, AmountMinor: 7500, Currency: "thb", SettlementDate: "2026-08-30", Note: "Paid back"})
	require.NoError(t, err)
	require.Equal(t, int64(7500), settlement.AmountMinor)
	require.Equal(t, domain.DefaultCurrency, settlement.Currency)

	balances, err := service.ListBalances(ctx, groupID, ownerID)
	require.NoError(t, err)
	require.Equal(t, int64(-2500), balanceForParticipant(balances, ownerParticipantID).AmountMinor)
	require.Equal(t, int64(2500), balanceForParticipant(balances, friendParticipantID).AmountMinor)

	require.NoError(t, service.DeleteSettlement(ctx, groupID, ownerID, settlement.ID))
	balances, err = service.ListBalances(ctx, groupID, ownerID)
	require.NoError(t, err)
	require.Equal(t, int64(5000), balanceForParticipant(balances, ownerParticipantID).AmountMinor)
	require.Equal(t, int64(-5000), balanceForParticipant(balances, friendParticipantID).AmountMinor)
}

func TestGroupServiceListsDeterministicSuggestedTransfers(t *testing.T) {
	ctx := context.Background()
	groupID, actorID := uuid.New(), uuid.New()
	creditorTwoUserID, debtorOneUserID, debtorTwoUserID := uuid.New(), uuid.New(), uuid.New()
	creditorOneID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	creditorTwoID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	debtorOneID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	debtorTwoID := uuid.MustParse("00000000-0000-0000-0000-000000000004")

	store := newFakeGroupStore()
	store.participants[participantKey(groupID, actorID)] = &domain.Participant{ID: creditorOneID, GroupID: groupID, UserID: actorID, Active: true}
	store.participants[participantKey(groupID, creditorTwoUserID)] = &domain.Participant{ID: creditorTwoID, GroupID: groupID, UserID: creditorTwoUserID, Active: true}
	store.participants[participantKey(groupID, debtorOneUserID)] = &domain.Participant{ID: debtorOneID, GroupID: groupID, UserID: debtorOneUserID, Active: true}
	store.participants[participantKey(groupID, debtorTwoUserID)] = &domain.Participant{ID: debtorTwoID, GroupID: groupID, UserID: debtorTwoUserID, Active: true}
	store.expenses = []domain.Expense{
		{GroupID: groupID, PayerParticipantID: creditorOneID, AmountMinor: 7000, Splits: []domain.ExpenseSplit{{ParticipantID: debtorTwoID, AmountMinor: 6000}, {ParticipantID: debtorOneID, AmountMinor: 1000}}},
		{GroupID: groupID, PayerParticipantID: creditorTwoID, AmountMinor: 3000, Splits: []domain.ExpenseSplit{{ParticipantID: debtorOneID, AmountMinor: 3000}}},
	}

	transfers, err := NewGroupService(store).ListSuggestedTransfers(ctx, groupID, actorID)

	require.NoError(t, err)
	actual := make([]struct {
		payerID    uuid.UUID
		receiverID uuid.UUID
		amount     int64
	}, 0, len(transfers))
	for _, transfer := range transfers {
		actual = append(actual, struct {
			payerID    uuid.UUID
			receiverID uuid.UUID
			amount     int64
		}{transfer.PayerParticipantID, transfer.ReceiverParticipantID, transfer.AmountMinor})
	}
	require.Equal(t, []struct {
		payerID    uuid.UUID
		receiverID uuid.UUID
		amount     int64
	}{
		{debtorOneID, creditorOneID, 4000},
		{debtorTwoID, creditorOneID, 3000},
		{debtorTwoID, creditorTwoID, 3000},
	}, actual)
}

func TestSuggestedTransfersFindsTheMinimumNumberOfRepayments(t *testing.T) {
	participant := func(id string) domain.Participant { return domain.Participant{ID: uuid.MustParse(id)} }
	transfers := suggestedTransfers([]Balance{
		{Participant: participant("00000000-0000-0000-0000-000000000001"), AmountMinor: 800},
		{Participant: participant("00000000-0000-0000-0000-000000000002"), AmountMinor: 600},
		{Participant: participant("00000000-0000-0000-0000-000000000003"), AmountMinor: 400},
		{Participant: participant("00000000-0000-0000-0000-000000000004"), AmountMinor: -1000},
		{Participant: participant("00000000-0000-0000-0000-000000000005"), AmountMinor: -800},
	})

	require.Equal(t, []SuggestedTransfer{
		{PayerParticipantID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), ReceiverParticipantID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), AmountMinor: 600},
		{PayerParticipantID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), ReceiverParticipantID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), AmountMinor: 400},
		{PayerParticipantID: uuid.MustParse("00000000-0000-0000-0000-000000000005"), ReceiverParticipantID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), AmountMinor: 800},
	}, transferIdentity(transfers))
}

func TestGroupServiceListsSimpleSuggestedTransfer(t *testing.T) {
	ctx := context.Background()
	groupID, actorID, debtorUserID := uuid.New(), uuid.New(), uuid.New()
	creditorID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	debtorID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, actorID)] = &domain.Participant{ID: creditorID, GroupID: groupID, UserID: actorID, Active: true}
	store.participants[participantKey(groupID, debtorUserID)] = &domain.Participant{ID: debtorID, GroupID: groupID, UserID: debtorUserID, Active: true}
	store.expenses = []domain.Expense{{GroupID: groupID, PayerParticipantID: creditorID, AmountMinor: 5000, Splits: []domain.ExpenseSplit{{ParticipantID: debtorID, AmountMinor: 5000}}}}

	transfers, err := NewGroupService(store).ListSuggestedTransfers(ctx, groupID, actorID)

	require.NoError(t, err)
	require.Equal(t, []SuggestedTransfer{{PayerParticipantID: debtorID, ReceiverParticipantID: creditorID, AmountMinor: 5000}}, transferIdentity(transfers))
}

func TestSuggestedTransfersHandlesZeroBalancesOverpaymentAndTies(t *testing.T) {
	participant := func(id string) domain.Participant { return domain.Participant{ID: uuid.MustParse(id)} }
	zeroBalances := []Balance{{Participant: participant("00000000-0000-0000-0000-000000000001")}}
	require.Empty(t, suggestedTransfers(zeroBalances))

	transfers := suggestedTransfers([]Balance{
		{Participant: participant("00000000-0000-0000-0000-000000000004"), AmountMinor: 5000},
		{Participant: participant("00000000-0000-0000-0000-000000000002"), AmountMinor: -10000},
		{Participant: participant("00000000-0000-0000-0000-000000000003"), AmountMinor: 5000},
		{Participant: participant("00000000-0000-0000-0000-000000000001"), AmountMinor: 0},
	})
	require.Equal(t, []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
	}, []uuid.UUID{transfers[0].PayerParticipantID, transfers[1].PayerParticipantID})
	require.Equal(t, []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		uuid.MustParse("00000000-0000-0000-0000-000000000004"),
	}, []uuid.UUID{transfers[0].ReceiverParticipantID, transfers[1].ReceiverParticipantID})
	require.Equal(t, []int64{5000, 5000}, []int64{transfers[0].AmountMinor, transfers[1].AmountMinor})
}

func transferIdentity(transfers []SuggestedTransfer) []SuggestedTransfer {
	identities := make([]SuggestedTransfer, 0, len(transfers))
	for _, transfer := range transfers {
		identities = append(identities, SuggestedTransfer{PayerParticipantID: transfer.PayerParticipantID, ReceiverParticipantID: transfer.ReceiverParticipantID, AmountMinor: transfer.AmountMinor})
	}
	return identities
}

func balanceForParticipant(balances []Balance, participantID uuid.UUID) Balance {
	for _, balance := range balances {
		if balance.ParticipantID == participantID {
			return balance
		}
	}
	return Balance{}
}

func TestGroupServiceAllowsOnlyRepaymentSidesToRecordSettlements(t *testing.T) {
	ctx := context.Background()
	groupID, payerUserID, receiverUserID, outsiderID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	payerID, receiverID := uuid.New(), uuid.New()
	store := newFakeGroupStore()
	store.participants[participantKey(groupID, payerUserID)] = &domain.Participant{ID: payerID, GroupID: groupID, UserID: payerUserID, Active: true}
	store.participants[participantKey(groupID, receiverUserID)] = &domain.Participant{ID: receiverID, GroupID: groupID, UserID: receiverUserID, Active: true}
	store.participants[participantKey(groupID, outsiderID)] = &domain.Participant{ID: uuid.New(), GroupID: groupID, UserID: outsiderID, Active: true}
	service := NewGroupService(store)
	_, err := service.CreateSettlement(ctx, groupID, outsiderID, CreateSettlementInput{PayerParticipantID: payerID, ReceiverParticipantID: receiverID, AmountMinor: 100, Currency: "THB", SettlementDate: "2026-08-30"})
	require.ErrorIs(t, err, ErrForbidden)
	_, err = service.CreateSettlement(ctx, groupID, receiverUserID, CreateSettlementInput{PayerParticipantID: payerID, ReceiverParticipantID: receiverID, AmountMinor: 100, Currency: "THB", SettlementDate: "2026-08-30"})
	require.NoError(t, err)
}
