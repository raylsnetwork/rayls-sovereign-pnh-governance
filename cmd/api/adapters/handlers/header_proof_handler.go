package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

// HeaderProofHandler handles HTTP requests for header proof operations
type HeaderProofHandler struct {
	service core.HeaderProofService
	log     logger.Logger
}

// NewHeaderProofHandler creates a new header proof handler
func NewHeaderProofHandler(service core.HeaderProofService, log logger.Logger) *HeaderProofHandler {
	return &HeaderProofHandler{
		service: service,
		log:     log,
	}
}

// GetHeaderProofsList godoc
// @Summary      List Header Proofs by Block Range
// @Description  Retrieves blockchain header proofs for a specific Privacy Node (PN) within a given block range. Returns paginated results ordered by block number. Returns an empty array if no headers match the filters.
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Param        chainId      query     string  true   "Chain ID of the Privacy Node"  example("11155111")
// @Param        startBlock   query     string  true   "Start block number (inclusive)"       example("1")
// @Param        endBlock     query     string  true   "End block number (inclusive)"         example("100")
// @Param        page         query     int     false  "Page number (default: 1)"             example(1)
// @Param        pageSize     query     int     false  "Items per page (default: 50, max: 1000)" example(50)
// @Success      200  {object}  dto.HeaderProofListResponse "List of header proofs, empty array if none found"
// @Failure      400  {object}  map[string]string "Invalid parameters specified"
// @Failure      500  {object}  map[string]string "Database error"
// @Router       /audit/header-proofs [get]
func (h *HeaderProofHandler) GetHeaderProofsList(c *gin.Context) {
	var filters dto.HeaderProofFilters

	// Bind query parameters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters specified", "details": err.Error()})
		return
	}

	// Call service
	response, err := h.service.GetHeaderProofsList(c.Request.Context(), filters)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
