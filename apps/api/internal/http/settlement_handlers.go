package http

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/service"
)

type createSettlementRequest struct {
	PayerParticipantID    string `json:"payerParticipantId"`
	ReceiverParticipantID string `json:"receiverParticipantId"`
	AmountMinor           int64  `json:"amountMinor"`
	Currency              string `json:"currency"`
	SettlementDate        string `json:"settlementDate"`
	Note                  string `json:"note"`
}
type settlementResponse struct {
	ID                    string              `json:"id"`
	PayerParticipantID    string              `json:"payerParticipantId"`
	Payer                 participantResponse `json:"payer"`
	ReceiverParticipantID string              `json:"receiverParticipantId"`
	Receiver              participantResponse `json:"receiver"`
	AmountMinor           int64               `json:"amountMinor"`
	Currency              string              `json:"currency"`
	SettlementDate        string              `json:"settlementDate"`
	Note                  string              `json:"note"`
}

func (s *Server) createSettlement(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, 400, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}
	var req createSettlementRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, 400, "INVALID_JSON", "Request body is invalid.", nil)
	}
	payerID, payerErr := uuid.Parse(req.PayerParticipantID)
	receiverID, receiverErr := uuid.Parse(req.ReceiverParticipantID)
	if payerErr != nil || receiverErr != nil {
		return writeError(c, 400, "VALIDATION_FAILED", "Settlement Participants are invalid.", nil)
	}
	settlement, err := s.groups.CreateSettlement(c.UserContext(), groupID, currentUser(c).ID, service.CreateSettlementInput{PayerParticipantID: payerID, ReceiverParticipantID: receiverID, AmountMinor: req.AmountMinor, Currency: req.Currency, SettlementDate: req.SettlementDate, Note: req.Note})
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, 403, "FORBIDDEN", "You cannot record this Settlement.", nil)
	}
	if errors.Is(err, service.ErrValidation) {
		return writeError(c, 400, "VALIDATION_FAILED", "Settlement details are invalid.", nil)
	}
	if err != nil {
		return err
	}
	return c.Status(201).JSON(fiber.Map{"settlement": toSettlementResponse(settlement)})
}
func (s *Server) listSettlements(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, 400, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}
	settlements, err := s.groups.ListSettlements(c.UserContext(), groupID, currentUser(c).ID)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, 403, "FORBIDDEN", "You cannot view Settlements for this Group.", nil)
	}
	if err != nil {
		return err
	}
	response := make([]settlementResponse, 0, len(settlements))
	for i := range settlements {
		response = append(response, toSettlementResponse(&settlements[i]))
	}
	return c.JSON(fiber.Map{"settlements": response})
}
func (s *Server) deleteSettlement(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, 400, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}
	settlementID, err := uuid.Parse(c.Params("settlementID"))
	if err != nil {
		return writeError(c, 400, "VALIDATION_FAILED", "Settlement ID is invalid.", nil)
	}
	err = s.groups.DeleteSettlement(c.UserContext(), groupID, currentUser(c).ID, settlementID)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, 403, "FORBIDDEN", "You cannot delete this Settlement.", nil)
	}
	if err != nil {
		return writeError(c, 404, "NOT_FOUND", "Settlement was not found.", nil)
	}
	return c.SendStatus(204)
}
func toSettlementResponse(settlement *domain.Settlement) settlementResponse {
	return settlementResponse{ID: settlement.ID.String(), PayerParticipantID: settlement.PayerParticipantID.String(), Payer: toParticipantResponse(&settlement.PayerParticipant), ReceiverParticipantID: settlement.ReceiverParticipantID.String(), Receiver: toParticipantResponse(&settlement.ReceiverParticipant), AmountMinor: settlement.AmountMinor, Currency: settlement.Currency, SettlementDate: settlement.SettlementDate.In(time.UTC).Format("2006-01-02"), Note: strings.TrimSpace(settlement.Note)}
}
