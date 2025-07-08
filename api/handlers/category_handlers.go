package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/afreedicp/zolaris-backend-app/internal/services"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/response"
	"github.com/afreedicp/zolaris-backend-app/internal/utils"
)

// AddCategoryHandler handles requests to add a new category
type AddCategoryHandler struct {
	categoryService *services.CategoryService
}

// NewAddCategoryHandler creates a new AddCategoryHandler
func NewAddCategoryHandler(categoryService *services.CategoryService) *AddCategoryHandler {
	return &AddCategoryHandler{categoryService: categoryService}
}

// HandleGin handles requests using Gin framework
// @Summary Add a new category
// @Description Register a new category
// @Tags Category Management
// @Accept json
// @Produce json
// @Param category body dto.CategoryRequest true "Category information"
// @Success 201 {object} dto.Response "Category added successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 409 {object} dto.ErrorResponse "Category already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /category/add [post]
func (h *AddCategoryHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request dto.CategoryRequest
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

	// Call service to add category
	if err := h.categoryService.AddCategory(c.Request.Context(), request.Name, request.Type); err != nil {
		if err.Error() == "category with this name already exists" {
			response.Error(c, http.StatusConflict, "Category with this name already exists", "CONFLICT")
			return
		}
		log.Printf("Error adding category: %v", err)
		response.InternalError(c, "Failed to add category")
		return
	}

	response.Created(c, nil, "Category added successfully")
}

// GetCategoriesByTypeHandler handles requests to get categories by type
type GetCategoriesByTypeHandler struct {
	categoryService *services.CategoryService
}

// NewGetCategoriesByTypeHandler creates a new GetCategoriesByTypeHandler
func NewGetCategoriesByTypeHandler(categoryService *services.CategoryService) *GetCategoriesByTypeHandler {
	return &GetCategoriesByTypeHandler{categoryService: categoryService}
}

// HandleGin handles requests using Gin framework
// @Summary Get categories by type
// @Description Retrieve all categories of a specific type
// @Tags Category Management
// @Accept json
// @Produce json
// @Param type path string true "Category type"
// @Success 200 {array} dto.CategoryResponse "List of categories"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /category/type/{type} [get]
func (h *GetCategoriesByTypeHandler) HandleGin(c *gin.Context) {
	// Get type from URL parameter
	categoryType := c.Param("type")
	if categoryType == "" {
		response.BadRequest(c, "Category type is required")
		return
	}

	// Call service to get categories by type
	categories, err := h.categoryService.GetCategoriesByType(c.Request.Context(), categoryType)
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		response.InternalError(c, "Failed to get categories")
		return
	}

	// Return empty array if no categories found
	if categories == nil {
		categories = []*dto.CategoryResponse{}
	}

	response.OK(c, categories, "Categories retrieved successfully")
}

// ListAllCategoriesHandler handles requests to list all categories
type ListAllCategoriesHandler struct {
	categoryService *services.CategoryService
}

// NewListAllCategoriesHandler creates a new ListAllCategoriesHandler
func NewListAllCategoriesHandler(categoryService *services.CategoryService) *ListAllCategoriesHandler {
	return &ListAllCategoriesHandler{categoryService: categoryService}
}

// HandleGin handles requests using Gin framework
// @Summary Get all categories
// @Description Retrieve all categories
// @Tags Category Management
// @Produce json
// @Success 200 {array} dto.CategoryResponse "List of categories"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /category/all [get]
func (h *ListAllCategoriesHandler) HandleGin(c *gin.Context) {
	// Call service to get all categories
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		response.InternalError(c, "Failed to get categories")
		return
	}

	// Return empty array if no categories found
	if categories == nil {
		categories = []*dto.CategoryResponse{}
	}

	response.OK(c, categories, "Categories retrieved successfully")
}

