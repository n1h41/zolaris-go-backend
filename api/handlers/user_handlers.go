package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/afreedicp/zolaris-backend-app/internal/middleware"
	"github.com/afreedicp/zolaris-backend-app/internal/services"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/mappers"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/response"
	"github.com/afreedicp/zolaris-backend-app/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// HandleGin handles GET /user/details requests
// @Summary Get user details
// @Description Retrieve authenticated user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Param X-Cognito-ID header string true "Cognito ID"
// @Success 200 {object} dto.Response{data=dto.UserResponse} "User details retrieved successfully"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/details [get]
func (h *UserHandler) HandleGetUserDetails(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get user from service
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving user details: %v", err)
		response.InternalError(c, "Failed to retrieve user details")
		return
	}

	if user == nil {
		response.NotFound(c, "User not found")
		return
	}

	// Convert domain model to response DTO
	userResponse := mappers.UserToResponse(user)
	response.OK(c, userResponse, "User details retrieved successfully")
}

// HandleGin handles POST /user/details requests
// @Summary Update user details
// @Description Update the authenticated user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param user body dto.UserDetailsRequest true "User details"
// @Success 200 {object} dto.Response{data=dto.UserResponse} "User details updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/details [post]
func (h *UserHandler) HandleUpdateUserDetails(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Parse request body
	var request dto.UserDetailsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))

		// Convert validation errors to DTO format
		var validationErrDTOs []dto.ValidationError
		for _, item := range validationErrs {
			validationErrDTOs = append(validationErrDTOs, dto.ValidationError{
				Field:   item.Field,
				Message: item.Error,
			})
		}

		response.ValidationErrors(c, validationErrDTOs)
		return
	}

	// Update user details
	updatedUser, err := h.userService.UpdateUserDetails(c.Request.Context(), userID, &request)
	if err != nil {
		log.Printf("Error updating user details: %v", err)
		response.InternalError(c, "Failed to update user details")
		return
	}

	// Convert domain model to response DTO
	userResponse := mappers.UserToResponse(updatedUser)
	response.OK(c, userResponse, "User details updated successfully")
}

// HandleGin handles GET /user/check-parent-id requests
// @Summary Check if user has parent ID
// @Description Checks if the authenticated user has a parent ID set in their profile
// @Tags User Management
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} map[string]bool "Returns has_parent_id flag"
// @Failure 400 {object} map[string]string "Error when user ID is not found in context"
// @Failure 500 {object} map[string]string "Error when checking parent ID fails"
// @Router /user/check-parent-id [get]
func (h *UserHandler) HandleCheckHasParentID(c *gin.Context) {
	// Extract user ID from the request context
	userID, exists := c.Get("userID")
	if !exists {
		response.BadRequest(c, "User ID not found in context")
		return
	}

	// Check if the user has a parent ID
	hasParentID, err := h.userService.CheckHasParentID(c.Request.Context(), userID.(string))
	if err != nil {
		log.Printf("Error checking parent ID: %v", err)
		response.InternalError(c, "Failed to check parent ID")
		return
	}

	response.OK(c, gin.H{"has_parent_id": hasParentID}, "Success")
}

// HandleListReferredUsers handles GET /user/referrals requests
// @Summary List referred users
// @Description Retrieve a list of users referred by the authenticated user
// @Tags User Management
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} dto.Response{data=[]dto.UserResponse} "Referred users retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "User ID not found in context"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/referrals [get]
func (h *UserHandler) HandleListReferredUsers(c *gin.Context) {
	log.Printf("Error listing referred usersasd:")
	userID, exists := c.Get("userID")
	if !exists {
		response.BadRequest(c, "User ID not found in context")
		return
	}

	referredUsers, err := h.userService.ListReferredUsers(c.Request.Context(), userID.(string))
	if err != nil {
		log.Printf("Error listing referred users: %v", err)
		response.InternalError(c, "Failed to list referred users")
		return
	}

	response.OK(c, mappers.UsersToResponses(referredUsers), "Referred users retrieved successfully")
}




// CreateUserDetails handles POST /user/createUser requests
// @Summary Create user details
// @Description Create a new user record in the system based on Cognito ID and request data
// @Tags User Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param user body dto.UserDetailsRequest true "User details"
// @Success 201 {object} dto.Response{data=dto.UserResponse} "User details created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/createUser [post]
func (h *UserHandler) CreateUserDetails(c *gin.Context) {

	var request dto.UserDetailsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request format: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}
	// Validate input
	if validationErrs := utils.Validate(request); validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		var validationErrDTOs []dto.ValidationError
		for _, ve := range validationErrs {
			validationErrDTOs = append(validationErrDTOs, dto.ValidationError{
				Field:   ve.Field,
				Message: ve.Error,
			})
		}
		response.ValidationErrors(c, validationErrDTOs)
		return
	}
	// Call service to create user
	createdUser, err := h.userService.CreateUser(c.Request.Context(), &request)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		response.InternalError(c, "Failed to create user")
		return
	}

	// Map domain model to response DTO
	userResp := mappers.UserToResponse(createdUser)
	response.Created(c, userResp, "User created successfully")
}

