package services

import (
	"context"
	"errors"
	"log"

	"github.com/afreedicp/zolaris-backend-app/internal/repositories"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/mappers"
)

// CategoryService handles business logic for category operations
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryService creates a new category service instance
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// AddCategory handles the business logic for adding a new category
func (s *CategoryService) AddCategory(ctx context.Context, name, categoryType string) error {
	log.Printf("Adding category %s of type %s", name, categoryType)

	// Check if category already exists
	existingCategory, err := s.categoryRepo.GetCategoryByName(ctx, name)
	if err != nil {
		return err
	}

	if existingCategory != nil {
		return errors.New("category with this name already exists")
	}

	return s.categoryRepo.AddCategory(ctx, name, categoryType)
}

// GetCategoryByName retrieves a category by its name
func (s *CategoryService) GetCategoryByName(ctx context.Context, name string) (*dto.CategoryResponse, error) {
	log.Printf("Getting category with name %s", name)
	category, err := s.categoryRepo.GetCategoryByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Convert domain category to DTO
	if category != nil {
		return mappers.CategoryToResponse(category), nil
	}

	return nil, nil
} // GetCategoriesByType retrieves all categories of a specific type
func (s *CategoryService) GetCategoriesByType(ctx context.Context, categoryType string) ([]*dto.CategoryResponse, error) {
	log.Printf("Getting categories of type %s", categoryType)
	categories, err := s.categoryRepo.GetCategoriesByType(ctx, categoryType)
	if err != nil {
		return nil, err
	}

	return mappers.CategoriesToResponses(categories), nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*dto.CategoryResponse, error) {
	log.Println("List all categories")
	categories, err := s.categoryRepo.ListAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	return mappers.CategoriesToResponses(categories), nil
}
