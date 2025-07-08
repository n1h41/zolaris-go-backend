package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/afreedicp/zolaris-backend-app/internal/domain"
)

// CategoryRepository handles all category-related database operations
type CategoryRepository struct {
	db *pgxpool.Pool
}

// NewCategoryRepository creates a new category repository instance
func NewCategoryRepository(dbPool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		db: dbPool,
	}
}

// AddCategory adds a new category to the database
// This implementation matches the current signature while internally
// creating a proper Category object
func (r *CategoryRepository) AddCategory(ctx context.Context, name, categoryType string) error {
	query := `
		INSERT INTO z_category (
			category_id, name, type, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5)
	`

	// Generate a new UUID for the category
	categoryID := uuid.New().String()
	now := time.Now()

	_, err := r.db.Exec(
		ctx,
		query,
		categoryID,
		name,
		categoryType,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to add category: %w", err)
	}

	return nil
}

// GetCategoryByName retrieves a category by its name
func (r *CategoryRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	query := `
		SELECT category_id, name, type, created_at
		FROM z_category 
		WHERE name = $1
	`

	row := r.db.QueryRow(ctx, query, name)

	category := &domain.Category{}
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Type,
		&category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Category not found, return nil without error
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return category, nil
}

// GetCategoriesByType retrieves all categories of a specific type
func (r *CategoryRepository) GetCategoriesByType(ctx context.Context, categoryType string) ([]*domain.Category, error) {
	query := `
		SELECT category_id, name, type, created_at
		FROM z_category 
		WHERE type = $1
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query, categoryType)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		category := &domain.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Type,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning category row: %w", err)
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating category rows: %w", err)
	}

	return categories, nil
}

// ListAllCategories retrieves all categories from the database
func (r *CategoryRepository) ListAllCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `
		SELECT category_id, name, type, created_at
		FROM z_category
		ORDER BY type, name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		category := &domain.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Type,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning category row: %w", err)
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating category rows: %w", err)
	}

	return categories, nil
}
