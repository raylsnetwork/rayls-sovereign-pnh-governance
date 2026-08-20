package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// TokenResponse is the HTTP response format for tokens
type TokenResponse struct {
	ResourceId        string                          `json:"resourceId"`
	Name              string                          `json:"name"`
	Symbol            string                          `json:"symbol"`
	MetadataUrl       string                          `json:"metadataURL"`
	Decimals          uint8                           `json:"decimals"`
	IssuerId          string                          `json:"issuerId"`
	Status            string                          `json:"status"`
	ErcStandard       string                          `json:"ercStandard"`
	TotalSupply       decimal.Decimal                 `json:"totalSupply"`
	CirculatingSupply []domain.CirculatingSupplyEntry `json:"circulatingSupply"`
	FrozenChainIds    []string                        `json:"frozenChainIds"`
	CreatedAt         time.Time                       `json:"createdAt"`
	UpdatedAt         time.Time                       `json:"updatedAt"`
}

// TokenHandler handles HTTP requests for token operations
type TokenHandler struct {
	service core.TokenService
	log     logger.Logger
}

// NewTokenHandler creates a new token handler
func NewTokenHandler(service core.TokenService, log logger.Logger) *TokenHandler {
	return &TokenHandler{
		service: service,
		log:     log,
	}
}

// GetTokenByResourceId godoc
// @Summary      Get Token by Resource ID
// @Description  Retrieves a token by its resource ID
// @Tags         Audit
// @Param        resourceId   path      string  true  "The resource ID of the token"
// @Success      200      {object}  handlers.TokenResponse
// @Failure      400      {object}  map[string]string  "Invalid resourceId specified"
// @Failure      404      {object}  map[string]string  "No tokens found"
// @Failure      500      {object}  map[string]string  "Database error"
// @Router       /audit/tokens/{resourceId} [get]
func (h *TokenHandler) GetTokenByResourceId(c *gin.Context) {
	resourceId := c.Param("resourceId")

	token, err := h.service.GetTokenByResourceId(c.Request.Context(), resourceId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	response := h.tokenToResponse(*token)
	c.JSON(http.StatusOK, response)
}

// GetTokenList godoc
// @Summary      List Tokens by Query Parameters
// @Description  Retrieves a paginated list of tokens based on specified query parameters for auditing. Returns an empty array if no tokens match the filters.
// @Tags         Audit
// @Param        limit          query     int     false  "Items per page (Default: 10, Max: 100)"
// @Param        page           query     int     false  "Page number (Default: 1)"
// @Param        name           query     string  false  "Token's name (partial match)"
// @Param        issuerId       query     string  false  "Issuer ID"
// @Param        status         query     string  false  "Status of the token."  Enums(new,active,inactive)
// @Param        ercStandard    query     string  false  "Token type." Enums(erc20, erc721, erc1155, enygma, custom message)
// @Param        symbol         query     string  false  "Symbol of the token"
// @Param        decimals       query     int     false  "Decimals of the token"
// @Param        createdAfter   query     string  false  "Creation date after which tokens were added (inclusive). Accepts Unix ts, YYYY-MM-DD or ISO8601 formats"
// @Param        createdBefore  query     string  false  "Creation date before which tokens were added (exclusive). Accepts Unix ts, YYYY-MM-DD or ISO8601 formats"
// @Success      200  {object}  types.Paginated[handlers.TokenResponse]
// @Failure      400            {object}  map[string]string "Invalid parameters specified"
// @Failure      500            {object}  map[string]string "Database error"
// @Router       /audit/tokens [get]
func (h *TokenHandler) GetTokenList(c *gin.Context) {
	var filters dto.TokenListFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters", "details": err.Error()})
		return
	}

	result, err := h.service.GetTokensList(c.Request.Context(), filters)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	response := make([]TokenResponse, len(result.Data))
	for i, token := range result.Data {
		response[i] = h.tokenToResponse(token)
	}

	c.JSON(http.StatusOK, types.Paginated[TokenResponse]{
		Data:  response,
		Total: result.Total,
		Limit: result.Limit,
		Page:  result.Page,
	})
}

// GetTokenRegistryStatus godoc
// @Summary      Get Token Registry Status
// @Description  Retrieves token status from the database
// @Tags         Blockchain
// @Accept       json
// @Produce      json
// @Param        resourceid  path  string  true  "Resource ID"
// @Success      200  {object}  dto.TokenRegistryStatusDto
// @Failure      400  {object}  map[string]string "Invalid resource ID"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Token not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Security     CookieAuth
// @Router       /token_status/{resourceid} [get]
func (h *TokenHandler) GetTokenRegistryStatus(c *gin.Context) {
	resourceId := c.Param("resourceid")

	tokenStatus, err := h.service.GetTokenRegistryStatus(c.Request.Context(), resourceId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, tokenStatus)
}

// tokenToResponse converts domain TokenWithBalancesAndFreezeState to HTTP TokenResponse
func (h *TokenHandler) tokenToResponse(token domain.TokenWithBalancesAndFreezeState) TokenResponse {
	return TokenResponse{
		ResourceId:        token.ResourceId,
		Name:              token.Name,
		Symbol:            token.Symbol,
		MetadataUrl:       token.MetadataUrl,
		Decimals:          token.Decimals,
		IssuerId:          token.IssuerId,
		Status:            domain.TokenStatusToString[int(token.Status)],
		ErcStandard:       types.AssetTypeToString[uint8(token.ErcStandard)],
		TotalSupply:       token.TotalSupply,
		CirculatingSupply: token.CirculatingSupply,
		FrozenChainIds:    token.FrozenChainIds,
		CreatedAt:         token.CreatedAt,
		UpdatedAt:         token.UpdatedAt,
	}
}
