package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// BalanceHandler handles HTTP requests for balance operations
type BalanceHandler struct {
	service core.BalanceService
	log     logger.Logger
}

// NewBalanceHandler creates a new balance handler
func NewBalanceHandler(service core.BalanceService, log logger.Logger) *BalanceHandler {
	return &BalanceHandler{
		service: service,
		log:     log,
	}
}

// GetBalancesInChain godoc
// @Summary      Get Balances in Chain
// @Description  Retrieves all balances for a chain or specific resource balance
// @Description  If resourceid is "/" returns all balances, otherwise returns specific resource
// @Tags         Balances
// @Accept       json
// @Produce      json
// @Param        chainid     path  string  true   "Chain ID"
// @Param        resourceid  path  string  false  "Resource ID (use / for all)"
// @Success      200  {object}  interface{} "Array of balances or single balance object"
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     CookieAuth
// @Router       /resources/{chainid}/{resourceid} [get]
func (h *BalanceHandler) GetBalancesInChain(c *gin.Context) {
	chainId := c.Param("chainid")
	resourceId := c.Param("resourceid")

	result, err := h.service.GetBalancesInChain(c.Request.Context(), chainId, resourceId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetBalanceAcrossAllChains godoc
// @Summary      Get Balance Across All Chains
// @Description  Retrieves balance information for a resource across all chains
// @Tags         Balances
// @Accept       json
// @Produce      json
// @Param        resourceid  path  string  true  "Resource ID"
// @Success      200  {array}   domain.Balance "Array of balances across chains"
// @Failure      400  {object}  map[string]string "Invalid request"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Resource not found in any chain"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     CookieAuth
// @Router       /resource_info_all_chains/{resourceid} [get]
func (h *BalanceHandler) GetBalanceAcrossAllChains(c *gin.Context) {
	resourceId := c.Param("resourceid")

	balances, err := h.service.GetBalanceAcrossAllChains(c.Request.Context(), resourceId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, balances)
}

// GetBalanceAcrossSpecificChains godoc
// @Summary      Get Balance Across Specific Chains
// @Description  Retrieves balance information for a resource across specific chains
// @Tags         Balances
// @Accept       json
// @Produce      json
// @Param        body  body  object{resource_id=string,chains=[]string}  true  "Request body"
// @Success      200  {array}   domain.Balance "Array of balances across specified chains"
// @Failure      400  {object}  map[string]string "Invalid request body"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Resource not found in specified chains"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     CookieAuth
// @Router       /resource_info_list_chains [post]
func (h *BalanceHandler) GetBalanceAcrossSpecificChains(c *gin.Context) {
	var request struct {
		ResourceId string   `json:"resource_id" binding:"required"`
		Chains     []string `json:"chains"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	balances, err := h.service.GetBalanceAcrossSpecificChains(c.Request.Context(), request.ResourceId, request.Chains)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, balances)
}
