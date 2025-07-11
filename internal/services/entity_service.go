package services

import (
	"context"
	"fmt"

	"n1h41/zolaris-backend-app/internal/domain"
	"n1h41/zolaris-backend-app/internal/repositories"
)

// EntityService provides entity-related business operations
type EntityService struct {
	repo repositories.EntityRepository
}

// NewEntityService creates a new entity service with the provided repository
func NewEntityService(repo repositories.EntityRepository) *EntityService {
	return &EntityService{
		repo: repo,
	}
}

// CheckEntityExists determines if an entity exists for a given user
func (s *EntityService) CheckEntityExists(ctx context.Context, userId string) (bool, error) {
	if userId == "" {
		return false, fmt.Errorf("user ID cannot be empty")
	}

	return s.repo.CheckEntityPresence(ctx, userId)
}

// CreateRootEntity creates a new top-level entity without a parent
func (s *EntityService) CreateRootEntity(ctx context.Context, categoryId string, entityName string, userId string, details map[string]any) (string, error) {
	if categoryId == "" {
		return "", fmt.Errorf("category ID cannot be empty")
	}

	if entityName == "" {
		return "", fmt.Errorf("entity name cannot be empty")
	}

	// If details is nil, initialize it as an empty map
	if details == nil {
		details = make(map[string]any)
	}

	return s.repo.CreateRootEntity(ctx, categoryId, entityName, userId, details)
}

// CreateSubEntity creates a new entity as a child of an existing entity
func (s *EntityService) CreateSubEntity(ctx context.Context, categoryId string, entityName string, userId string, details map[string]any, parentEntityID string, subUserID string) (string, error) {
	if categoryId == "" {
		return "", fmt.Errorf("category ID cannot be empty")
	}

	if entityName == "" {
		return "", fmt.Errorf("entity name cannot be empty")
	}

	// If details is nil, initialize it as an empty map
	if details == nil {
		details = make(map[string]any)
	}

	return s.repo.CreateSubEntity(ctx, categoryId, entityName, userId, details, parentEntityID, subUserID)
}

// GetChildEntities retrieves all direct child entities of a given entity
// If recursive is true, returns all descendants (children, grandchildren, etc.)
func (s *EntityService) GetChildEntities(ctx context.Context, entityId string, recursive bool) ([]*domain.Entity, error) {
	if entityId == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}

	return s.repo.GetChildEntities(ctx, entityId, recursive)
}

// GetEntityHierarchy retrieves an entity and all its descendant entities as a hierarchical structure
func (s *EntityService) GetEntityHierarchy(ctx context.Context, userId string) (map[string]any, error) {
	if userId == "" {
		return nil, fmt.Errorf("root entity ID cannot be empty")
	}

	return s.repo.GetEntityHierarchy(ctx, userId)
}

// ListEntityChildren lists all children of a given entity with optional filtering
// level: 0 for direct children only, -1 for all descendants, or specific depth (1, 2, 3, etc.)
// categoryType: filter by category type (optional)
func (s *EntityService) ListEntityChildren(ctx context.Context, userId string, level int, categoryType string) ([]*domain.Entity, error) {
	if userId == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}

	// Validate level parameter
	if level < -1 {
		return nil, fmt.Errorf("invalid level: must be -1 (all levels), 0 (direct children only), or a positive integer")
	}

	return s.repo.ListEntityChildren(ctx, userId, level, categoryType)
}

// GetCategoryType retrieves the type of a category by its ID
func (s *EntityService) GetCategoryType(ctx context.Context, categoryId string) (repositories.CategoryType, error) {
	if categoryId == "" {
		return "", fmt.Errorf("category ID cannot be empty")
	}

	return s.repo.GetCategoryType(ctx, categoryId)
}
