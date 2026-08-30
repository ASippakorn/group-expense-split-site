package http

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/service"
)

type createExpenseRequest struct {
	Description        string                      `json:"description"`
	AmountMinor        int64                       `json:"amountMinor"`
	Currency           string                      `json:"currency"`
	ExpenseDate        string                      `json:"expenseDate"`
	PayerParticipantID string                      `json:"payerParticipantId"`
	ParticipantIDs     []string                    `json:"participantIds"`
	SplitType          string                      `json:"splitType"`
	Splits             []createExpenseSplitRequest `json:"splits"`
	TagID              string                      `json:"tagId"`
}

type createExpenseSplitRequest struct {
	ParticipantID string `json:"participantId"`
	AmountMinor   int64  `json:"amountMinor"`
	Percentage    string `json:"percentage"`
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
	TagID              string                 `json:"tagId,omitempty"`
	Splits             []expenseSplitResponse `json:"splits"`
}

type expenseSplitResponse struct {
	ID            string              `json:"id"`
	ParticipantID string              `json:"participantId"`
	Participant   participantResponse `json:"participant"`
	AmountMinor   int64               `json:"amountMinor"`
	Percentage    string              `json:"percentage,omitempty"`
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

	splits, err := parseExpenseSplits(req.Splits)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Split values are invalid.", nil)
	}
	splitType := req.SplitType
	if splitType == "" {
		splitType = domain.SplitTypeEqual
	}
	tagID := uuid.Nil
	if strings.TrimSpace(req.TagID) != "" {
		tagID, err = uuid.Parse(req.TagID)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Tag details are invalid.", nil)
		}
	}
	expense, err := s.groups.CreateExpense(c.UserContext(), groupID, currentUser(c).ID, service.CreateExpenseInput{
		Description:        req.Description,
		AmountMinor:        req.AmountMinor,
		Currency:           req.Currency,
		ExpenseDate:        req.ExpenseDate,
		PayerParticipantID: payerParticipantID,
		ParticipantIDs:     participantIDs,
		SplitType:          splitType,
		Splits:             splits,
		TagID:              tagID,
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

func parseExpenseSplits(values []createExpenseSplitRequest) ([]service.CreateSplitInput, error) {
	splits := make([]service.CreateSplitInput, 0, len(values))
	for _, value := range values {
		participantID, err := uuid.Parse(value.ParticipantID)
		if err != nil {
			return nil, err
		}
		percentageBasisPoints := int64(0)
		if strings.TrimSpace(value.Percentage) != "" {
			percentageBasisPoints, err = parsePercentageBasisPoints(value.Percentage)
			if err != nil {
				return nil, err
			}
		}
		splits = append(splits, service.CreateSplitInput{ParticipantID: participantID, AmountMinor: value.AmountMinor, PercentageBasisPoints: percentageBasisPoints})
	}
	return splits, nil
}

func parsePercentageBasisPoints(value string) (int64, error) {
	trimmed := strings.TrimSpace(value)
	if !strings.Contains(trimmed, ".") {
		trimmed += ".0"
	}
	parts := strings.Split(trimmed, ".")
	if len(parts) != 2 || len(parts[1]) > 2 || parts[0] == "" {
		return 0, errors.New("invalid percentage")
	}
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 0 {
		return 0, errors.New("invalid percentage")
	}
	fraction := parts[1]
	if len(fraction) == 1 {
		fraction += "0"
	}
	fractionValue, err := strconv.ParseInt(fraction, 10, 64)
	if err != nil {
		return 0, err
	}
	return whole*100 + fractionValue, nil
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
			Percentage:    formatPercentage(expense.Splits[i].PercentageBasisPoints),
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
		TagID:              optionalUUIDString(expense.TagID),
		Splits:             splits,
	}
}

func optionalUUIDString(value *uuid.UUID) string {
	if value == nil {
		return ""
	}
	return value.String()
}

func formatPercentage(basisPoints int64) string {
	if basisPoints == 0 {
		return ""
	}
	fraction := strconv.FormatInt(basisPoints%100, 10)
	if len(fraction) == 1 {
		fraction = "0" + fraction
	}
	fraction = strings.TrimRight(fraction, "0")
	if fraction == "" {
		return strconv.FormatInt(basisPoints/100, 10)
	}
	return strconv.FormatInt(basisPoints/100, 10) + "." + fraction
}
