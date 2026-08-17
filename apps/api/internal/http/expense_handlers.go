package http

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/service"
)

type createExpenseRequest struct {
	Description        string   `json:"description"`
	AmountMinor        int64    `json:"amountMinor"`
	Currency           string   `json:"currency"`
	ExpenseDate        string   `json:"expenseDate"`
	PayerParticipantID string   `json:"payerParticipantId"`
	ParticipantIDs     []string `json:"participantIds"`
}

type expenseResponse struct {
	ID                 string                 `json:"id"`
	Description        string                 `json:"description"`
	AmountMinor        int64                  `json:"amountMinor"`
	Currency           string                 `json:"currency"`
	ExpenseDate        string                 `json:"expenseDate"`
	SplitType          string                 `json:"splitType"`
	PayerParticipantID string                 `json:"payerParticipantId"`
	Payer              participantResponse    `json:"payer"`
	Splits             []expenseSplitResponse `json:"splits"`
}

type expenseSplitResponse struct {
	ID            string              `json:"id"`
	ParticipantID string              `json:"participantId"`
	Participant   participantResponse `json:"participant"`
	AmountMinor   int64               `json:"amountMinor"`
}

func (s *Server) createExpense(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}

	var req createExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Request body is invalid.", nil)
	}

	payerParticipantID, participantIDs, err := parseExpenseParticipantIDs(req.PayerParticipantID, req.ParticipantIDs)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Expense Participants are invalid.", nil)
	}

	expense, err := s.groups.CreateEqualExpense(c.UserContext(), groupID, currentUser(c).ID, service.CreateEqualExpenseInput{
		Description:        req.Description,
		AmountMinor:        req.AmountMinor,
		Currency:           req.Currency,
		ExpenseDate:        req.ExpenseDate,
		PayerParticipantID: payerParticipantID,
		ParticipantIDs:     participantIDs,
	})
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "You cannot create Expenses for this Group.", nil)
	}
	if errors.Is(err, service.ErrValidation) {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Expense details are invalid.", nil)
	}
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"expense": toExpenseResponse(expense)})
}

func (s *Server) listExpenses(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}

	expenses, err := s.groups.ListExpenses(c.UserContext(), groupID, currentUser(c).ID)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "You cannot view Expenses for this Group.", nil)
	}
	if err != nil {
		return err
	}

	response := make([]expenseResponse, 0, len(expenses))
	for i := range expenses {
		response = append(response, toExpenseResponse(&expenses[i]))
	}
	return c.JSON(fiber.Map{"expenses": response})
}

func parseExpenseParticipantIDs(payerParticipantID string, participantIDValues []string) (uuid.UUID, []uuid.UUID, error) {
	payerID, err := uuid.Parse(payerParticipantID)
	if err != nil {
		return uuid.Nil, nil, err
	}

	participantIDs := make([]uuid.UUID, 0, len(participantIDValues))
	for _, value := range participantIDValues {
		participantID, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil, nil, err
		}
		participantIDs = append(participantIDs, participantID)
	}
	return payerID, participantIDs, nil
}

func toExpenseResponse(expense *domain.Expense) expenseResponse {
	splits := make([]expenseSplitResponse, 0, len(expense.Splits))
	for i := range expense.Splits {
		splits = append(splits, expenseSplitResponse{
			ID:            expense.Splits[i].ID.String(),
			ParticipantID: expense.Splits[i].ParticipantID.String(),
			Participant:   toParticipantResponse(&expense.Splits[i].Participant),
			AmountMinor:   expense.Splits[i].AmountMinor,
		})
	}

	return expenseResponse{
		ID:                 expense.ID.String(),
		Description:        expense.Description,
		AmountMinor:        expense.AmountMinor,
		Currency:           expense.Currency,
		ExpenseDate:        expense.ExpenseDate.In(time.UTC).Format("2006-01-02"),
		SplitType:          expense.SplitType,
		PayerParticipantID: expense.PayerParticipantID.String(),
		Payer:              toParticipantResponse(&expense.PayerParticipant),
		Splits:             splits,
	}
}
