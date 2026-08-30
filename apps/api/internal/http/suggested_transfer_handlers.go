package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/splitr/apps/api/internal/service"
)

type suggestedTransferResponse struct {
	PayerParticipantID    string              `json:"payerParticipantId"`
	Payer                 participantResponse `json:"payer"`
	ReceiverParticipantID string              `json:"receiverParticipantId"`
	Receiver              participantResponse `json:"receiver"`
	AmountMinor           int64               `json:"amountMinor"`
}

func (s *Server) listSuggestedTransfers(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}

	transfers, err := s.groups.ListSuggestedTransfers(c.UserContext(), groupID, currentUser(c).ID)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "You cannot view Suggested Transfers for this Group.", nil)
	}
	if err != nil {
		return err
	}

	response := make([]suggestedTransferResponse, 0, len(transfers))
	for _, transfer := range transfers {
		response = append(response, suggestedTransferResponse{
			PayerParticipantID:    transfer.PayerParticipantID.String(),
			Payer:                 toParticipantResponse(&transfer.PayerParticipant),
			ReceiverParticipantID: transfer.ReceiverParticipantID.String(),
			Receiver:              toParticipantResponse(&transfer.ReceiverParticipant),
			AmountMinor:           transfer.AmountMinor,
		})
	}
	return c.JSON(fiber.Map{"suggestedTransfers": response})
}
