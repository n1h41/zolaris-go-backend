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

// EntityHandler handles all entity-related HTTP requests
type EntityHandler struct {
	entityService *services.EntityService
}

// NewEntityHandler creates a new EntityHandler
func NewEntityHandler(entityService *services.EntityService) *EntityHandler {
	return &EntityHandler{entityService: entityService}
}

// HandleCreateRootEntity handles requests to create a root entity
// @Summary Create a root entity
// @Description Create a new top-level entity with optional user association
// @Tags Entity Management
// @Accept json
// @Produce json
// @Param X-Cognito-ID header string true "Cognito ID"
// @Param entity body dto.CreateRootEntityRequest true "Entity information"
// @Success 201 {object} dto.Response "Entity created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /entity/root [post]
func (h *EntityHandler) HandleCreateRootEntity(c *gin.Context) {
	// Parse request body
	var request dto.CreateRootEntityRequest
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

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Call service to create root entity
	entityID, err := h.entityService.CreateRootEntity(
		c.Request.Context(),
		request.CategoryID,
		request.Name,
		userID,
		request.Details,
	)
	if err != nil {
		log.Printf("Error creating root entity: %v", err)
		response.InternalError(c, "Failed to create root entity")
		return
	}

	response.Created(c, map[string]string{"entityId": entityID}, "Root entity created successfully")
}

// HandleCreateSubEntity handles requests to create a sub-entity
// @Summary Create a sub-entity
// @Description Create a new entity as a child of an existing entity
// @Tags Entity Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param entity body dto.CreateSubEntityRequest true "Entity information"
// @Success 201 {object} dto.Response "Sub-entity created successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 404 {object} dto.ErrorResponse "Parent entity not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /entity/sub [post]
func (h *EntityHandler) HandleCreateSubEntity(c *gin.Context) {
	// Parse request body
	var request dto.CreateSubEntityRequest
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

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Call service to create sub-entity
	entityID, err := h.entityService.CreateSubEntity(
		c.Request.Context(),
		request.CategoryID,
		request.Name,
		userID,
		request.Details,
		request.ParentEntityID,
	)
	if err != nil {
		if err.Error() == "user with ID "+userID+" does not have any existing entities" {
			response.BadRequest(c, "User does not have any existing entities")
			return
		}
		log.Printf("Error creating sub-entity: %v", err)
		response.InternalError(c, "Failed to create sub-entity")
		return
	}

	response.Created(c, map[string]string{"entityId": entityID}, "Sub-entity created successfully")
}

// HandleGetEntityChildren handles requests to get children of an entity
// @Summary Get entity children
// @Description Get all children of a specific entity, with optional recursion and filtering
// @Tags Entity Management
// @Accept json
// @Produce json
// @Param entity_id path string true "Entity ID"
// @Param recursive query bool false "Whether to include all descendants"
// @Param level query int false "Maximum depth level for descendants (0 for direct children only, -1 for all)"
// @Param category_type query string false "Filter by category type"
// @Success 200 {object} dto.Response{data=dto.EntityChildrenResponse} "Entity children retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Entity not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /entity/{entity_id}/children [get]
func (h *EntityHandler) HandleGetEntityChildren(c *gin.Context) {
	// Get entity ID from URL path
	entityID := c.Param("entity_id")
	if entityID == "" {
		response.BadRequest(c, "Entity ID is required")
		return
	}

	// Parse query parameters
	var request dto.GetEntityChildrenRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		response.BadRequest(c, "Invalid query parameters")
		return
	}

	// Set level based on recursive flag and level parameter
	level := 0
	if request.Recursive {
		if request.Level != 0 {
			level = request.Level
		} else {
			level = -1 // All levels
		}
	}

	// Call service to get entity children
	entities, err := h.entityService.ListEntityChildren(
		c.Request.Context(),
		entityID,
		level,
		request.CategoryType,
	)
	if err != nil {
		if err.Error() == "entity with ID "+entityID+" not found" {
			response.NotFound(c, "Entity not found")
			return
		}
		log.Printf("Error getting entity children: %v", err)
		response.InternalError(c, "Failed to retrieve entity children")
		return
	}

	// Convert entities to responses
	childResponses := mappers.EntitiesToResponses(entities)

	// Create the response structure
	result := dto.EntityChildrenResponse{
		ParentID: entityID,
		Children: childResponses,
		Count:    len(childResponses),
	}

	response.OK(c, result, "Entity children retrieved successfully")
}

// HandleGetEntityHierarchy handles requests to get an entity hierarchy
// @Summary Get entity hierarchy
// @Description Get an entity and all its descendants as a hierarchical structure
// @Tags Entity Management
// @Accept json
// @Produce json
// @Param entity_id path string true "Entity ID"
// @Param max_depth query int false "Maximum depth to include (default: 10)"
// @Success 200 {object} dto.Response{data=dto.EntityHierarchyResponse} "Entity hierarchy retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Entity not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /entity/{entity_id}/hierarchy [get]
func (h *EntityHandler) HandleGetEntityHierarchy(c *gin.Context) {
	// Get entity ID from URL path
	entityID := c.Param("entity_id")
	if entityID == "" {
		response.BadRequest(c, "Entity ID is required")
		return
	}

	// Call service to get entity hierarchy
	hierarchy, err := h.entityService.GetEntityHierarchy(
		c.Request.Context(),
		entityID,
	)
	if err != nil {
		if err.Error() == "entity with ID "+entityID+" not found" {
			response.NotFound(c, "Entity not found")
			return
		}
		log.Printf("Error getting entity hierarchy: %v", err)
		response.InternalError(c, "Failed to retrieve entity hierarchy")
		return
	}

	// Convert hierarchy map to response structure
	hierarchyResponse := mappers.HierarchyMapToResponse(hierarchy)
	if hierarchyResponse == nil {
		response.InternalError(c, "Failed to process entity hierarchy")
		return
	}

	response.OK(c, hierarchyResponse, "Entity hierarchy retrieved successfully")
}

// HandleCheckEntityPresence handles requests to check if a user has any entities
// @Summary Check entity presence
// @Description Check if the authenticated user has any entities
// @Tags Entity Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} dto.Response{data=map[string]bool} "Entity presence check successful"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/has-entity [get]
func (h *EntityHandler) HandleCheckEntityPresence(c *gin.Context) {
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Check if the user has an entity
	hasEntity, err := h.entityService.CheckEntityExists(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error checking entity presence: %v", err)
		response.InternalError(c, "Failed to check entity presence")
		return
	}

	result := map[string]interface{}{"hasEntity": hasEntity}

	if hasEntity {
		entityID, err := h.entityService.GetEntityID(c.Request.Context(), userID)
		if err != nil {
			log.Printf("Error fetching entity ID: %v", err)
			response.InternalError(c, "Failed to fetch entity ID")
			return
		}
		result["entityId"] = entityID
	}

	response.OK(c, result, "Entity presence check successful")
}
