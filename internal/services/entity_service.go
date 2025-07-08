package services

import (
	"context"
	"fmt"
	"github.com/afreedicp/zolaris-backend-app/internal/domain"
	"github.com/afreedicp/zolaris-backend-app/internal/repositories"
)

// EntityService provides entity-related business operations
type EntityService struct {
	repo repositories.EntityRepository
	userRepo  repositories.UserRepositoryInterface 
}

// NewEntityService creates a new entity service with the provided repository
func NewEntityService(repo repositories.EntityRepository, userRepo repositories.UserRepositoryInterface)*EntityService {
	return &EntityService{
		repo: repo,
		userRepo: userRepo,
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
func (s *EntityService) CreateSubEntity(ctx context.Context, categoryId string, entityName string, userId string, details map[string]any, parentEntityID string) (string, error) {
	if categoryId == "" {
		return "", fmt.Errorf("category ID cannot be empty")
	}

	if entityName == "" {
		return "", fmt.Errorf("entity name cannot be empty")
	}
	if details == nil {
		details = make(map[string]any)
	}

	parentCategoryID, err := s.repo.GetCategoryIDByEntityID(ctx, parentEntityID)
	if err != nil {
		return "", fmt.Errorf("failed to get parent's category ID: %w", err)
	}

	parentCategoryType, err := s.repo.GetCategoryType(ctx, parentCategoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get parent's category type: %w", err)
	}

	currentCategoryType, err := s.repo.GetCategoryType(ctx, categoryId)
	if err != nil {
		return "", fmt.Errorf("failed to get parent's category type: %w", err)
	}

	subentityID, err := s.repo.CreateSubEntity(ctx, categoryId, entityName, userId, details, parentEntityID)
	if err != nil {
		return "", fmt.Errorf("failed to create sub-entity: %w", err)
	}
	if parentCategoryType == "user" &&  currentCategoryType == "user" {
		subuserRaw, ok := details["subuser_id"]
		if !ok {
			return "", fmt.Errorf("subuser_id not found in details : %w", parentCategoryType)
		}
		subuserID, ok := subuserRaw.(string)
		if !ok {
			return "", fmt.Errorf("subuser_id must be a string")
		}


		if err := s.userRepo.UpdateUserParentID(ctx, subuserID, &parentEntityID); err != nil {
			return "", fmt.Errorf("failed to update user parent ID: %w", err)
		}
	}
	return subentityID, nil
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
func (s *EntityService) GetEntityHierarchy(ctx context.Context, rootEntityId string) (map[string]any, error) {
	if rootEntityId == "" {
		return nil, fmt.Errorf("root entity ID cannot be empty")
	}

	return s.repo.GetEntityHierarchy(ctx, rootEntityId)
}

// ListEntityChildren lists all children of a given entity with optional filtering
// level: 0 for direct children only, -1 for all descendants, or specific depth (1, 2, 3, etc.)
// categoryType: filter by category type (optional)
func (s *EntityService) ListEntityChildren(ctx context.Context, entityId string, level int, categoryType string) ([]*domain.Entity, error) {
	if entityId == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}

	// Validate level parameter
	if level < -1 {
		return nil, fmt.Errorf("invalid level: must be -1 (all levels), 0 (direct children only), or a positive integer")
	}

	return s.repo.ListEntityChildren(ctx, entityId, level, categoryType)
}




func (s *EntityService) GetEntityID(ctx context.Context, userId string) (string, error) {
	if userId == "" {
		return "", fmt.Errorf("user ID is empty")
	}

	return s.repo.GetEntityID(ctx, userId) // You must have this repo method implemented
}