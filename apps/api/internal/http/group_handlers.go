package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/your-org/splitr/apps/api/internal/domain"
	"github.com/your-org/splitr/apps/api/internal/service"
)

type createGroupRequest struct {
	Name            string `json:"name"`
	DefaultCurrency string `json:"defaultCurrency"`
	Description     string `json:"description"`
}

type groupResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	DefaultCurrency string `json:"defaultCurrency"`
	Description     string `json:"description"`
	OwnerID         string `json:"ownerId"`
}

type addParticipantRequest struct {
	Email string `json:"email"`
}

type participantResponse struct {
	ID     string       `json:"id"`
	User   userResponse `json:"user"`
	Role   string       `json:"role"`
	Active bool         `json:"active"`
}

type balanceResponse struct {
	Participant     participantResponse `json:"participant"`
	PaidAmountMinor int64               `json:"paidAmountMinor"`
	OwedAmountMinor int64               `json:"owedAmountMinor"`
	AmountMinor     int64               `json:"amountMinor"`
}

type createTagRequest struct {
	Name           string   `json:"name"`
	ParticipantIDs []string `json:"participantIds"`
}

type tagResponse struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Participants []participantResponse `json:"participants"`
}

func (s *Server) createGroup(c *fiber.Ctx) error {
	var req createGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Request body is invalid.", nil)
	}

	group, err := s.groups.CreateGroup(c.UserContext(), currentUser(c).ID, req.Name, req.DefaultCurrency, req.Description)
	if errors.Is(err, service.ErrValidation) {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group name and currency are required.", nil)
	}
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"group": toGroupResponse(group)})
}

func (s *Server) listGroups(c *fiber.Ctx) error {
	groups, err := s.groups.ListGroups(c.UserContext(), currentUser(c).ID)
	if err != nil {
		return err
	}

	response := make([]groupResponse, 0, len(groups))
	for i := range groups {
		response = append(response, toGroupResponse(&groups[i]))
	}
	return c.JSON(fiber.Map{"groups": response})
}

func toGroupResponse(group *domain.Group) groupResponse {
	return groupResponse{
		ID:              group.ID.String(),
		Name:            group.Name,
		DefaultCurrency: group.DefaultCurrency,
		Description:     group.Description,
		OwnerID:         group.OwnerID.String(),
	}
}

func (s *Server) addParticipant(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}

	var req addParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Request body is invalid.", nil)
	}

	participant, err := s.groups.AddParticipantByEmail(c.UserContext(), groupID, currentUser(c).ID, req.Email)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Only the Group Owner can add Participants.", nil)
	}
	if errors.Is(err, service.ErrUserNotFound) {
		return writeError(c, fiber.StatusBadRequest, "USER_NOT_FOUND", "No registered User has that email.", map[string]string{"email": "User not found."})
	}
	if errors.Is(err, service.ErrParticipantExists) {
		return writeError(c, fiber.StatusConflict, "PARTICIPANT_EXISTS", "That User is already a Participant.", map[string]string{"email": "Already in this Group."})
	}
	if errors.Is(err, service.ErrValidation) {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Email is invalid.", map[string]string{"email": "Enter a valid email."})
	}
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"participant": toParticipantResponse(participant)})
}

func (s *Server) listParticipants(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}

	participants, err := s.groups.ListParticipants(c.UserContext(), groupID, currentUser(c).ID)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "You cannot view Participants for this Group.", nil)
	}
	if err != nil {
		return err
	}

	response := make([]participantResponse, 0, len(participants))
	for i := range participants {
		response = append(response, toParticipantResponse(&participants[i]))
	}
	return c.JSON(fiber.Map{"participants": response})
}

func (s *Server) listBalances(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}

	balances, err := s.groups.ListBalances(c.UserContext(), groupID, currentUser(c).ID)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "You cannot view Balances for this Group.", nil)
	}
	if err != nil {
		return err
	}

	response := make([]balanceResponse, 0, len(balances))
	for _, balance := range balances {
		response = append(response, balanceResponse{
			Participant:     toParticipantResponse(&balance.Participant),
			PaidAmountMinor: balance.PaidAmountMinor,
			OwedAmountMinor: balance.OwedAmountMinor,
			AmountMinor:     balance.AmountMinor,
		})
	}
	return c.JSON(fiber.Map{"balances": response})
}

func (s *Server) createTag(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}
	var req createTagRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Request body is invalid.", nil)
	}
	participantIDs := make([]uuid.UUID, 0, len(req.ParticipantIDs))
	for _, participantID := range req.ParticipantIDs {
		id, err := uuid.Parse(participantID)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Tag details are invalid.", nil)
		}
		participantIDs = append(participantIDs, id)
	}
	tag, err := s.groups.CreateTag(c.UserContext(), groupID, currentUser(c).ID, req.Name, participantIDs)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "You cannot create Tags for this Group.", nil)
	}
	if errors.Is(err, service.ErrValidation) {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Tag details are invalid.", nil)
	}
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"tag": toTagResponse(tag)})
}

func (s *Server) listTags(c *fiber.Ctx) error {
	groupID, err := parseGroupID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Group ID is invalid.", nil)
	}
	tags, err := s.groups.ListTags(c.UserContext(), groupID, currentUser(c).ID)
	if errors.Is(err, service.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "You cannot view Tags for this Group.", nil)
	}
	if err != nil {
		return err
	}
	response := make([]tagResponse, 0, len(tags))
	for i := range tags {
		response = append(response, toTagResponse(&tags[i]))
	}
	return c.JSON(fiber.Map{"tags": response})
}

func toTagResponse(tag *domain.Tag) tagResponse {
	participants := make([]participantResponse, 0, len(tag.Participants))
	for i := range tag.Participants {
		participants = append(participants, toParticipantResponse(&tag.Participants[i]))
	}
	return tagResponse{ID: tag.ID.String(), Name: tag.Name, Participants: participants}
}

func parseGroupID(c *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(c.Params("groupID"))
}

func toParticipantResponse(participant *domain.Participant) participantResponse {
	return participantResponse{
		ID: participant.ID.String(),
		User: userResponse{
			ID:    participant.UserID.String(),
			Email: participant.User.Email,
		},
		Role:   participant.Role,
		Active: participant.Active,
	}
}
