package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// ParticipantHandler handles HTTP requests for participant operations
type ParticipantHandler struct {
	service core.ParticipantService
	log     logger.Logger
}

// NewParticipantHandler creates a new participant handler
func NewParticipantHandler(service core.ParticipantService, log logger.Logger) *ParticipantHandler {
	return &ParticipantHandler{
		service: service,
		log:     log,
	}
}

// GetParticipantByChainId godoc
// @Summary      Get a Participant by their ChainId
// @Description  Retrieves a participant's information by their chain ID if it exists in the database.
// @Tags         Audit
// @Param        chainId   path      string  true  "The chainId of the participant to retrieve"
// @Success      200      {object}  domain.Participant
// @Failure      400      {object}  map[string]string  "Invalid chainId specified"
// @Failure      404      {object}  map[string]string  "No participants found"
// @Failure      500      {object}  map[string]string  "Database error"
// @Router       /audit/participants/{chainId} [get]
// @Router       /participant_info/{chainId} [get]
func (h *ParticipantHandler) GetParticipantByChainId(c *gin.Context) {
	chainId := c.Param("chainId")

	participant, err := h.service.GetParticipantByChainId(c.Request.Context(), chainId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, participant)
}

// GetParticipantList godoc
// @Summary      List Participants by Query Parameters
// @Description  Retrieves a list of participants based on specified query parameters for auditing. Returns an empty array if no participants match the filters.
// @Tags         Audit
// @Param        name          query     string  false  "Participant's name (partial match)"
// @Param        chainId       query     string  false  "Chain ID"
// @Param        status        query     string  false  "Status of the participant."  Enums(new,active,inactive,frozen)
// @Param        role          query     string  false  "Role of the participant."  Enums(participant,issuer,auditor)
// @Param        createdAfter  query     string  false  "Creation date after which participants were added (inclusive). Accepts Unix ts, YYYY-MM-DD or ISO8601 formats"
// @Param        createdBefore query     string  false  "Creation date before which participants were added (exclusive). Accepts Unix ts, YYYY-MM-DD or ISO8601 formats"
// @Success      200  {array}   domain.Participant
// @Failure      400  {object}  map[string]string "Invalid parameters specified"
// @Failure      500  {object}  map[string]string "Database error"
// @Router       /audit/participants [get]
func (h *ParticipantHandler) GetParticipantList(c *gin.Context) {
	var filters dto.ParticipantListFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters", "details": err.Error()})
		return
	}

	participants, err := h.service.GetParticipantsList(c.Request.Context(), filters)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, participants)
}
