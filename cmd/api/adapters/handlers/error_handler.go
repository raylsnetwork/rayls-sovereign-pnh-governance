package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// HandleError maps domain errors to HTTP status codes
func HandleError(c *gin.Context, log logger.Logger, err error) {
	var notFoundErr *core.NotFoundError
	var validationErr *core.ValidationError
	var internalErr *core.InternalError
	var conflictErr *core.ConflictError

	switch {
	case errors.As(err, &conflictErr):
		log.Debug("Conflict error", "resource", conflictErr.Resource)
		c.JSON(http.StatusConflict, gin.H{"error": conflictErr.Error()})

	case errors.As(err, &notFoundErr):
		log.Debug("Resource not found", "resource", notFoundErr.Resource, "id", notFoundErr.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})

	case errors.As(err, &validationErr):
		log.Debug("Validation error", "field", validationErr.Field, "message", validationErr.Message)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})

	case errors.As(err, &internalErr):
		log.Error("Internal error", "error", internalErr.Err, "operation", internalErr.Operation)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

	default:
		log.Error("Unknown error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
