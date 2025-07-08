package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/afreedicp/zolaris-backend-app/internal/domain"
)

type CategoryType string

const (
	UserCategoryType     CategoryType = "user"
	OfficeCategoryType   CategoryType = "office"
	LocationCategoryType CategoryType = "location"
)

type EntityRepository struct {
	db *pgxpool.Pool
}

func NewEntityRepository(dbPool *pgxpool.Pool) EntityRepository {
	return EntityRepository{
		db: dbPool,
	}
}

func (r *EntityRepository) CheckEntityPresence(ctx context.Context, userId string) (bool, error) {
	var exists bool
	query := `select exists(select 1 from z_entity where user_id = $1)`

	if err := r.db.QueryRow(ctx, query, userId).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check entity presence: %w", err)
	}

	return exists, nil
}

func (r *EntityRepository) GetCategoryType(ctx context.Context, categoryId string) (CategoryType, error) {
	var categoryType string
	query := `SELECT type FROM z_category WHERE category_id = $1`
	err := r.db.QueryRow(ctx, query, categoryId).Scan(&categoryType) // Fixed: scanning into categoryType instead of categoryId
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("category with ID %s not found", categoryId)
		}
		return "", fmt.Errorf("failed to get category type: %w", err)
	}
	return CategoryType(categoryType), nil
}

func (r *EntityRepository) GetCategoryIDByEntityID(ctx context.Context, entityID string) (string, error) {
	var categoryID string
	query := `SELECT category_id FROM z_entity WHERE entity_id = $1`
	err := r.db.QueryRow(ctx, query, entityID).Scan(&categoryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("entity with ID %s not found", entityID)
		}
		return "", fmt.Errorf("failed to get category ID: %w", err)
	}
	return categoryID, nil
}

func (r *EntityRepository) CreateRootEntity(ctx context.Context, categoryId string, entityName string, userId string, details map[string]any) (string, error) {
	categoryType, err := r.GetCategoryType(ctx, categoryId)
	if err != nil {
		return "", err
	}

	var entityId string

	if categoryType == "" {
		return "", fmt.Errorf("category with ID %s not found", categoryId)
	}

	if categoryType == UserCategoryType {
		if userId == "" {
			return "", fmt.Errorf("user ID is required for user category entities")
		}

		query := `insert into z_entity (category_id, name, user_id) values ($1, $2, $3) returning entity_id`

		err = r.db.QueryRow(ctx, query, categoryId, entityName, userId).Scan(&entityId)
	} else {
		detailsJSON, jsonErr := json.Marshal(details)
		if jsonErr != nil {
			return "", fmt.Errorf("failed to marshal details: %w", jsonErr)
		}

		query := `insert into z_entity (category_id, name, details) values ($1, $2, $3) returning entity_id`
		err = r.db.QueryRow(ctx, query, categoryId, entityName, detailsJSON).Scan(&entityId)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create root entity: %w", err)
	}

	return entityId, nil
}

func (r *EntityRepository) CreateSubEntity(ctx context.Context, categoryId string, entityName string, userId string, details map[string]any, parentEntityId string) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if parentEntityId == "" {
		getParentEntityIDQuery := `select entity_id from z_entity where user_id = $1 limit 1`
		err = tx.QueryRow(ctx, getParentEntityIDQuery, userId).Scan(&parentEntityId)
		if err != nil {
			return "", fmt.Errorf("failed to check parent entity: %w", err)
		}
	}

	categoryType, err := r.GetCategoryType(ctx, categoryId)
	if err != nil {
		return "", err
	}

	if categoryType == UserCategoryType && userId == "" {
		return "", fmt.Errorf("user ID is requried for user category entities")
	}

	var entityId string

	if userId != "" {
		var userHasEntity bool
		checkUserQuery := `select exists(select 1 from z_entity where user_id = $1)`
		if err := tx.QueryRow(ctx, checkUserQuery, userId).Scan(&userHasEntity); err != nil {
			return "", fmt.Errorf("failed to check user entities: %w", err)
		}

		if !userHasEntity {
			return "", fmt.Errorf("user with ID %s does not have any existing entities", userId)
		}

		if categoryType == UserCategoryType {
			query := `INSERT INTO z_entity (category_id, parent_id, name, user_id) VALUES ($1, $2, $3, $4) RETURNING entity_id`
			err = tx.QueryRow(ctx, query, categoryId, parentEntityId, entityName, userId).Scan(&entityId)
		} else {
			detailsJSON, jsonErr := json.Marshal(details)
			if jsonErr != nil {
				return "", fmt.Errorf("failed to marshal details: %w", jsonErr)
			}

			query := `INSERT INTO z_entity (category_id, parent_id, name, details) VALUES ($1, $2, $3, $4) RETURNING entity_id`

			err = tx.QueryRow(ctx, query, categoryId, parentEntityId, entityName, detailsJSON).Scan(&entityId)
		}
	}

	if err != nil {
		return "", fmt.Errorf("failed to create sub-entity: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return entityId, nil
}

// GetChildEntities retrieves all direct child entities of a given entity.
// If recursive is true, returns all descendants (children, grandchildren, etc.)
func (r *EntityRepository) GetChildEntities(ctx context.Context, entityId string, recursive bool) ([]*domain.Entity, error) {
	// First check if the parent entity exists
	var exists bool
	checkEntityQuery := `SELECT EXISTS(SELECT 1 FROM z_entity WHERE entity_id = $1)`

	if err := r.db.QueryRow(ctx, checkEntityQuery, entityId).Scan(&exists); err != nil {
		return nil, fmt.Errorf("failed to check entity existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("entity with ID %s not found", entityId)
	}

	// Get the path of the parent entity
	var parentPath string
	pathQuery := `SELECT path::text FROM z_entity WHERE entity_id = $1`

	if err := r.db.QueryRow(ctx, pathQuery, entityId).Scan(&parentPath); err != nil {
		return nil, fmt.Errorf("failed to get entity path: %w", err)
	}

	// Build query based on recursive flag
	var query string
	if recursive {
		// Query all descendants
		query = `
			SELECT 
				e.entity_id, e.user_id, e.name, e.details, e.category_id, 
				e.parent_id, e.path::text, e.depth, e.created_at, e.updated_at
			FROM 
				z_entity e
			WHERE 
				e.path <@ $1::ltree 
				AND e.entity_id != $2
			ORDER BY 
				e.path, e.depth
		`
	} else {
		// Query only immediate children
		query = `
			SELECT 
				e.entity_id, e.user_id, e.name, e.details, e.category_id, 
				e.parent_id, e.path::text, e.depth, e.created_at, e.updated_at
			FROM 
				z_entity e
			WHERE 
				e.parent_id = $1
			ORDER BY 
				e.name
		`
	}

	// Execute the query
	var rows pgx.Rows
	var err error

	if recursive {
		rows, err = r.db.Query(ctx, query, parentPath, entityId)
	} else {
		rows, err = r.db.Query(ctx, query, entityId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query child entities: %w", err)
	}
	defer rows.Close()

	// Process results
	entities := make([]*domain.Entity, 0)
	for rows.Next() {
		entity := new(domain.Entity)

		if err := rows.Scan(
			&entity.ID,
			&entity.UserID,
			&entity.Name,
			&entity.Details,
			&entity.CategoryID,
			&entity.ParentID,
			&entity.Path,
			&entity.Depth,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan entity row: %w", err)
		}

		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating entity rows: %w", err)
	}

	return entities, nil
}

// GetEntityHierarchy retrieves an entity and all its descendant entities as a hierarchy
// This method provides a structured view of the entity tree with proper parent-child relationships
func (r *EntityRepository) GetEntityHierarchy(ctx context.Context, rootEntityId string) (map[string]any, error) {
	// First check if the root entity exists
	var exists bool
	checkEntityQuery := `SELECT EXISTS(SELECT 1 FROM z_entity WHERE entity_id = $1)`

	if err := r.db.QueryRow(ctx, checkEntityQuery, rootEntityId).Scan(&exists); err != nil {
		return nil, fmt.Errorf("failed to check entity existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("entity with ID %s not found", rootEntityId)
	}

	// Get all entities in the hierarchy using recursive CTE
	query := `
		WITH RECURSIVE entity_hierarchy AS (
			-- Base case: the root entity
			SELECT 
				e.entity_id, 
				e.user_id, 
				e.name, 
				e.details, 
				e.category_id,
				c.name AS category_name,
				c.type AS category_type,
				e.parent_id, 
				e.path::text,
				e.depth, 
				e.created_at, 
				e.updated_at
			FROM 
				z_entity e
				JOIN z_category c ON e.category_id = c.category_id
			WHERE 
				e.entity_id = $1
				
			UNION ALL
			
			-- Recursive case: find all children
			SELECT 
				child.entity_id, 
				child.user_id, 
				child.name, 
				child.details, 
				child.category_id,
				c.name AS category_name,
				c.type AS category_type,
				child.parent_id, 
				child.path::text,
				child.depth, 
				child.created_at, 
				child.updated_at
			FROM 
				z_entity child
				JOIN entity_hierarchy parent ON child.parent_id = parent.entity_id
				JOIN z_category c ON child.category_id = c.category_id
		)
		SELECT 
			entity_id, 
			user_id, 
			name, 
			details, 
			category_id,
			category_name,
			category_type,
			parent_id, 
			path,
			depth, 
			created_at, 
			updated_at
		FROM 
			entity_hierarchy
		ORDER BY 
			depth, name
	`

	rows, err := r.db.Query(ctx, query, rootEntityId)
	if err != nil {
		return nil, fmt.Errorf("failed to query entity hierarchy: %w", err)
	}
	defer rows.Close()

	// Store all entities to build the hierarchy
	allEntities := make(map[string]map[string]any)
	var rootEntity map[string]any

	// Build a flat map of all entities
	for rows.Next() {
		var (
			entityId     string
			userId       *string
			name         string
			details      []byte
			categoryId   string
			categoryName string
			categoryType string
			parentId     *string
			path         string
			depth        int
			createdAt    time.Time
			updatedAt    time.Time
		)

		if err := rows.Scan(
			&entityId,
			&userId,
			&name,
			&details,
			&categoryId,
			&categoryName,
			&categoryType,
			&parentId,
			&path,
			&depth,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan entity row: %w", err)
		}

		// Create entity object
		entity := make(map[string]any)
		entity["id"] = entityId
		entity["name"] = name
		entity["category_id"] = categoryId
		entity["category_name"] = categoryName
		entity["category_type"] = categoryType
		entity["path"] = path
		entity["depth"] = depth
		entity["created_at"] = createdAt
		entity["updated_at"] = updatedAt
		entity["children"] = []map[string]any{} // Will be populated later

		// Add optional fields
		if userId != nil {
			entity["user_id"] = *userId
		}

		if parentId != nil {
			entity["parent_id"] = *parentId
		}

		// Parse details JSON if present
		if len(details) > 0 && string(details) != "null" {
			var detailsMap map[string]any
			if err := json.Unmarshal(details, &detailsMap); err != nil {
				return nil, fmt.Errorf("failed to unmarshal details JSON: %w", err)
			}
			entity["details"] = detailsMap
		} else {
			entity["details"] = map[string]any{}
		}

		// Store entity in the map
		allEntities[entityId] = entity

		// If this is the root entity, save a reference
		if entityId == rootEntityId {
			rootEntity = entity
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating entity rows: %w", err)
	}

	// If no root entity was found, something went wrong
	if rootEntity == nil {
		return nil, fmt.Errorf("root entity not found in query results")
	}

	// Build the hierarchy - for each entity, add it to its parent's children array
	for _, entity := range allEntities {
		parentId, hasParent := entity["parent_id"]
		if hasParent && parentId != nil {
			// Check if this entity is not the root entity itself
			if entity["id"].(string) != rootEntityId {
				if parent, exists := allEntities[parentId.(string)]; exists {
					children := parent["children"].([]map[string]any)
					parent["children"] = append(children, entity)
				}
			}
		}
	}

	return rootEntity, nil
}

// ListEntityChildren lists all children of a given entity with optional filtering
// level: 0 for direct children only, -1 for all descendants, or specific depth (1, 2, 3, etc.)
// categoryType: filter by category type (optional)
func (r *EntityRepository) ListEntityChildren(ctx context.Context, entityId string, level int, categoryType string) ([]*domain.Entity, error) {
	// First check if the entity exists
	var exists bool
	checkEntityQuery := `SELECT EXISTS(SELECT 1 FROM z_entity WHERE entity_id = $1)`

	if err := r.db.QueryRow(ctx, checkEntityQuery, entityId).Scan(&exists); err != nil {
		return nil, fmt.Errorf("failed to check entity existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("entity with ID %s not found", entityId)
	}

	// Build the query based on parameters
	var args []any
	args = append(args, entityId)

	// Base query
	query := `
		WITH RECURSIVE entity_children AS (
			-- Base case: direct children of the entity
			SELECT 
				e.entity_id, e.user_id, e.name, e.details, e.category_id, 
				e.parent_id, e.path::text, e.depth, e.created_at, e.updated_at,
				c.type as category_type, 
				1 as level
			FROM 
				z_entity e
				JOIN z_category c ON e.category_id = c.category_id
			WHERE 
				e.parent_id = $1
	`

	// Add category filter if provided
	if categoryType != "" {
		query += ` AND c.type = $2`
		args = append(args, categoryType)
	}

	// Continue recursive query
	query += `
			UNION ALL
			
			-- Recursive case: children of children
			SELECT 
				child.entity_id, child.user_id, child.name, child.details, child.category_id, 
				child.parent_id, child.path::text, child.depth, child.created_at, child.updated_at,
				c.type as category_type,
				parent.level + 1
			FROM 
				z_entity child
				JOIN entity_children parent ON child.parent_id = parent.entity_id
				JOIN z_category c ON child.category_id = c.category_id
	`

	// Apply level filter for recursive case if specified
	if level > 0 {
		query += fmt.Sprintf(" WHERE parent.level < %d", level)
	}

	// Close the CTE and add final selection
	query += `
		)
		SELECT 
			entity_id, user_id, name, details, category_id, 
			parent_id, path, depth, created_at, updated_at
		FROM 
			entity_children
	`

	// Add category filter for the final selection if provided
	if categoryType != "" {
		query += ` WHERE category_type = $2`
	}

	// Add ordering
	query += ` ORDER BY depth, name`

	// Execute query
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query entity children: %w", err)
	}
	defer rows.Close()

	// Process results
	entities := make([]*domain.Entity, 0)
	for rows.Next() {
		entity := new(domain.Entity)

		var rawDetails []byte
		if err := rows.Scan(
			&entity.ID,
			&entity.UserID,
			&entity.Name,
			&rawDetails,
			&entity.CategoryID,
			&entity.ParentID,
			&entity.Path,
			&entity.Depth,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan entity row: %w", err)
		}

		entity.Details = rawDetails
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating entity rows: %w", err)
	}

	return entities, nil
}

func (r *EntityRepository) GetEntityID(ctx context.Context, userId string) (string, error) {
    var entityID string
    query := `SELECT entity_id FROM z_entity WHERE user_id = $1 LIMIT 1`
    err := r.db.QueryRow(ctx, query, userId).Scan(&entityID)
    if err != nil {
        if err == pgx.ErrNoRows {
            return "", fmt.Errorf("no entity found for user ID %s", userId)
        }
        return "", fmt.Errorf("failed to get entity ID: %w", err)
    }
    return entityID, nil
}



