package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/utils"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// TransactionHandler handles HTTP requests for transaction operations
// Maps HTTP layer to business logic
type TransactionHandler struct {
	service core.TransactionService
	log     logger.Logger
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(service core.TransactionService, log logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		log:     log,
	}
}

// GetTransactionByMessageId godoc
// @Summary Get Transaction by Message ID
// @Description Retrieves a transaction by its unique message ID.
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.TransactionDetailDto
// @Failure      404  {object}  map[string]string "Message not found"
// @Failure      500  {object}  map[string]string "Error fetching/processing transaction data"
// @Router       /audit/transactions/{messageId} [get]
// @Param        messageId   path   string  true  "The message ID"
func (h *TransactionHandler) GetTransactionByMessageId(c *gin.Context) {
	messageId := c.Param("messageId")

	transaction, err := h.service.GetTransactionByMessageId(c.Request.Context(), messageId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetTransactionByTransactionId godoc
// @Summary Get Transaction by ID
// @Description Retrieves a transaction by its unique id.
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.TransactionDetailDto
// @Failure      404  {object}  map[string]string "Key not found"
// @Failure      500  {object}  map[string]string "Error fetching/processing transaction data"
// @Router       /audit/transactions/dvp/{transactionId} [get]
// @Param        transactionId   path   string  true  "The transaction id"
func (h *TransactionHandler) GetTransactionByTransactionId(c *gin.Context) {
	transactionId := c.Param("transactionId")

	transaction, err := h.service.GetTransactionByTransactionId(c.Request.Context(), transactionId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetTransactions godoc
// @Summary List Transactions by Query Parameters
// @Description Retrieves a paginated list of transactions based on specified query parameters for auditing. Returns an empty array in data field if no transactions match the filters.
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Paginated list of transactions"
// @Failure      400  {object}  map[string]string "Invalid filter parameters"
// @Failure      500  {object}  map[string]string "Error fetching transactions"
// @Router       /audit/transactions [get]
// @Param 			 limit							query  string  false  "Filter by latest transactions (Default: 10)"
// @Param 			 page               query  string  false  "Filter by pagination page (Default: 1)"
// @Param        messageId          query  string  false  "Filter by messageId."
// @Param        sourceChainId      query  string  false  "Filter by source chain ID"
// @Param        destinationChainId query  string  false  "Filter by destination chain ID"
// @Param        fromAddress        query  string  false  "Filter by source address"
// @Param        toAddress          query  string  false  "Filter by destination address"
// @Param        resourceId         query  string  false  "Filter by resource ID"
// @Param        messageType        query  string  false  "Filter by message type" Enums(erc20, erc721, erc1155, enygma, dvp_erc721, dvp_erc1155)
// @Param        initiatedAfter     query  string  false  "Filter transactions initiated after timestamp. Accepts Unix ts, YYYY-MM-DD or ISO8601 formats"
// @Param        initiatedBefore    query  string  false  "Filter transactions initiated before timestamp. Accepts Unix ts, YYYY-MM-DD or ISO8601 formats"
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	var filters dto.MergedTransactionsFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters", "details": err.Error()})
		return
	}

	result, err := h.service.GetTransactionsList(c.Request.Context(), filters)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTransactionsByEnygmaBatchId godoc
// @Summary      Get Aggregated Transactions by Batch ID
// @Description  Retrieves paginated enygma transactions with the same batch ID.
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Paginated list of batch transactions"
// @Failure      404  {object}  map[string]string "Batch not found"
// @Failure      500  {object}  map[string]string "Error fetching/processing transaction data"
// @Router       /audit/transactions/enygma/batch/{batchId} [get]
// @Param        batchId   path   string  true  "The batch ID"
// @Param        page      query  int     false "Page number (default: 1)"
// @Param        limit     query  int     false "Items per page (default: 10, max: 100)"
func (h *TransactionHandler) GetTransactionsByEnygmaBatchId(c *gin.Context) {
	batchId := c.Param("batchId")
	pagination := utils.ParsePaginationParams(c, utils.DefaultPage, utils.DefaultBatchPageSize)

	result, err := h.service.GetEnygmaBatchTransactions(
		c.Request.Context(),
		batchId,
		pagination.Page,
		pagination.Limit,
	)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTransactionsByRegularBatchId godoc
// @Summary      Get Transactions by Regular Batch ID
// @Description  Retrieves paginated transactions with the same batch ID (non-enygma batches).
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Paginated list of batch transactions"
// @Failure      404  {object}  map[string]string "Batch not found"
// @Failure      500  {object}  map[string]string "Error fetching/processing transaction data"
// @Router       /audit/transactions/batch/{batchId} [get]
// @Param        batchId   path   string  true  "The batch ID"
// @Param        page      query  int     false "Page number (default: 1)"
// @Param        limit     query  int     false "Items per page (default: 10, max: 100)"
func (h *TransactionHandler) GetTransactionsByRegularBatchId(c *gin.Context) {
	batchId := c.Param("batchId")
	pagination := utils.ParsePaginationParams(c, utils.DefaultPage, utils.DefaultBatchPageSize)

	result, err := h.service.GetRegularBatchTransactions(
		c.Request.Context(),
		batchId,
		pagination.Page,
		pagination.Limit,
	)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTransactionsBySharedId godoc
// @Summary Get Transactions by Shared ID
// @Description Retrieves all transactions with the same shared_id (ZkDvp swaps).
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.BatchTransactionDto "List of transactions with the same shared_id"
// @Failure      404  {object}  map[string]string "Transactions not found"
// @Failure      500  {object}  map[string]string "Error fetching/processing transaction data"
// @Router       /audit/transactions/dvp/swap/{sharedId} [get]
// @Param        sharedId   path   string  true  "The shared ID"
func (h *TransactionHandler) GetTransactionsBySharedId(c *gin.Context) {
	sharedId := c.Param("sharedId")

	transactions, err := h.service.GetTransactionsBySharedId(c.Request.Context(), sharedId)
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetFlaggedTransactions godoc
// @Summary      Get All Flagged Transactions
// @Description  Retrieves all transactions that have been flagged by the flagger service
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{} "flagged: array of flagged transactions"
// @Failure      500  {object}  map[string]string "Error getting flagged transactions"
// @Router       /flagged [get]
func (h *TransactionHandler) GetFlaggedTransactions(c *gin.Context) {
	flagged, err := h.service.GetFlaggedTransactions(c.Request.Context())
	if err != nil {
		HandleError(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"flagged": flagged})
}
