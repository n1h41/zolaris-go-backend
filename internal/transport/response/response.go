package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
)

// Success sends a successful response
func Success(c *gin.Context, statusCode int, data any, message string) {
	c.JSON(statusCode, dto.Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Created sends a successful creation response
func Created(c *gin.Context, data any, message string) {
	Success(c, http.StatusCreated, data, message)
}

// OK sends a successful response with 200 status code
func OK(c *gin.Context, data any, message string) {
	Success(c, http.StatusOK, data, message)
}

// NoContent sends a successful response with no content
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, errorCode string) {
	c.JSON(statusCode, dto.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    errorCode,
	})
}

// BadRequest sends a 400 error response
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message, "BAD_REQUEST")
}

// NotFound sends a 404 error response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, "NOT_FOUND")
}

// Unauthorized sends a 401 error response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, "UNAUTHORIZED")
}

// Forbidden sends a 403 error response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, "FORBIDDEN")
}

// InternalError sends a 500 error response
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message, "INTERNAL_ERROR")
}

// ValidationErrors sends a response with validation errors
func ValidationErrors(c *gin.Context, errors []dto.ValidationError) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Validation failed",
		"code":    "VALIDATION_ERROR",
		"details": errors,
	})
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, items any, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: dto.PaginatedResponse{
			Items:      items,
			TotalItems: total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}
