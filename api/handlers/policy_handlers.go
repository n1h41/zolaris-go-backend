package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/afreedicp/zolaris-backend-app/internal/services"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/response"
	"github.com/afreedicp/zolaris-backend-app/internal/utils"
)

// AttachIotPolicyHandler handles requests to attach an IoT policy to an identity
type AttachIotPolicyHandler struct {
	policyService *services.PolicyService
}

// NewAttachIotPolicyHandler creates a new AttachIotPolicyHandler
func NewAttachIotPolicyHandler(policyService *services.PolicyService) *AttachIotPolicyHandler {
	return &AttachIotPolicyHandler{policyService: policyService}
}

// HandleGin handles requests using Gin framework
// @Summary Attach IoT policy
// @Description Attach an AWS IoT policy to a Cognito identity
// @Tags Policy Management
// @Accept json
// @Produce json
// @Param request body dto.PolicyAttachRequest true "Identity information"
// @Success 200 {object} dto.Response "IoT policy attached successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} dto.ErrorResponse "Failed to attach IoT policy"
// @Router /device/attach-policy [post]
func (h *AttachIotPolicyHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request dto.PolicyAttachRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		response.ValidationErrors(c, utils.CreateDtoValidationErrors(validationErrs))
		return
	}

	// Call service to attach policy
	if err := h.policyService.AttachIoTPolicy(c.Request.Context(), request.IdentityID); err != nil {
		log.Printf("Error attaching IoT policy: %v", err)
		response.InternalError(c, "Failed to attach IoT policy")
		return
	}

	response.OK(c, nil, "IoT policy attached successfully")
}

