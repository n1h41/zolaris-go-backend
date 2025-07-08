package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
"log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/afreedicp/zolaris-backend-app/internal/domain"
)

// UserRepository handles all user-related database operations with PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(dbPool *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{
		db: dbPool,
	}
}

func (r *UserRepository) GetUserIdByCognitoId(ctx context.Context, cId string) (string, error) {
	var userId string

	query := `select user_id from z_users where cognito_id = $1`

	if err := r.db.QueryRow(ctx, query, cId).Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil // User not found
		}
		return "", fmt.Errorf("failed to get user ID by Cognito ID: %w", err)
	}

	return userId, nil
}

// GetUserByID retrieves a user by ID from PostgreSQL
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
    query := `
        SELECT user_id, email, first_name, last_name, phone,
               cognito_id, referral_mail, role, -- <--- ADDED THESE FIELDS HERE
               address, parent_id, created_at, updated_at
        FROM z_users
        WHERE user_id = $1
    `

    row := r.db.QueryRow(ctx, query, userID)

    user := &domain.User{}
    var addressJSON []byte
    var parentID *string

    err := row.Scan(
        &user.ID,
        &user.Email,
        &user.FirstName,
        &user.LastName,
        &user.Phone,
        &user.CognitoID,    // <--- ADDED SCAN DESTINATION
        &user.ReferralMail, // <--- ADDED SCAN DESTINATION
        &user.Role,         // <--- ADDED SCAN DESTINATION
        &addressJSON,
        &parentID,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil // User not found, return nil without error
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    // Parse address from JSON
    if len(addressJSON) > 0 && string(addressJSON) != "null" {
        if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
            return nil, fmt.Errorf("failed to parse address JSON: %w", err)
        }
    }

	user.ParentID = parentID // May be nil
	return user, nil
}

// CreateUser creates a new user in PostgreSQL
func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Convert address struct to JSON
	addressJSON, err := json.Marshal(user.Address)
	if err != nil {
		return fmt.Errorf("failed to convert address to JSON: %w", err)
	}

	query := `
		INSERT INTO z_users (
			user_id, email, first_name, last_name, phone,
			address, parent_id, cognito_id, referral_mail, role,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5,
				  $6, $7, $8, $9, $10,
				  $11, $12)
	`

	_, err = r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		addressJSON,
		user.ParentID,
		user.CognitoID,
		user.ReferralMail,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUser updates user in PostgreSQL
func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	// Convert address struct to JSON
	addressJSON, err := json.Marshal(user.Address)
	if err != nil {
		return fmt.Errorf("failed to convert address to JSON: %w", err)
	}

	query := `
		UPDATE z_users SET
			first_name = $1,
			last_name = $2,
			phone = $3,
			address = $4,
			updated_at = $5
		WHERE user_id = $6
	`

	// Update the timestamp
	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		addressJSON,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", user.ID)
	}

	return nil
}

// CheckHasParentID checks if a user has a parent ID in PostgreSQL
func (r *UserRepository) CheckHasParentID(ctx context.Context, userID string) (bool, error) {
	query := `
		SELECT 
			CASE WHEN parent_id IS NULL THEN false ELSE true END 
		FROM z_users 
		WHERE user_id = $1
	`

	var hasParent bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&hasParent)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("user not found: %w", err)
		}
		return false, fmt.Errorf("database error: %w", err)
	}

	return hasParent, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
        SELECT user_id, email, first_name, last_name, phone,
               cognito_id, referral_mail, role, -- <--- ADDED THESE FIELDS HERE
               address, parent_id, created_at, updated_at
        FROM z_users
        WHERE user_id = $1
    `

	row := r.db.QueryRow(ctx, query, email)

	user := &domain.User{}
	var addressJSON []byte
	var parentID *string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&addressJSON,
		&user.CognitoID,    // <--- ADDED SCAN DESTINATION
        &user.ReferralMail, // <--- ADDED SCAN DESTINATION
        &user.Role,         // <--- ADDED SCAN DESTINATION
		&parentID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User not found, return nil without error
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Parse address from JSON
	if len(addressJSON) > 0 && string(addressJSON) != "null" {
		if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
			return nil, fmt.Errorf("failed to parse address JSON: %w", err)
		}
	}

	user.ParentID = parentID
	return user, nil
}

// GetChildUsers gets all child users for a parent user
func (r *UserRepository) GetChildUsers(ctx context.Context, parentID string) ([]*domain.User, error) {
	query := `
        SELECT user_id, email, first_name, last_name, phone,
               cognito_id, referral_mail, role, -- <--- ADDED THESE FIELDS HERE
               address, parent_id, created_at, updated_at
        FROM z_users
        WHERE user_id = $1
    `

	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var addressJSON []byte
		var parentID *string

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&addressJSON,
			&user.CognitoID,    // <--- ADDED SCAN DESTINATION
        &user.ReferralMail, // <--- ADDED SCAN DESTINATION
        &user.Role,         // <--- ADDED SCAN DESTINATION
			&parentID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}

		// Parse address from JSON
		if len(addressJSON) > 0 && string(addressJSON) != "null" {
			if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
				return nil, fmt.Errorf("failed to parse address JSON: %w", err)
			}
		}

		user.ParentID = parentID
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

func (r *UserRepository) ListReferredUsers(ctx context.Context, userID string) ([]*domain.User, error) {
	// Step 1: Get the email of the referring user
	var referrerEmail string
	err := r.db.QueryRow(ctx, `
		SELECT email
		FROM z_users
		WHERE user_id = $1
	`, userID).Scan(&referrerEmail)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Referrer not found — treat as no referrals
			return []*domain.User{}, nil
		}
		return nil, fmt.Errorf("failed to get email for user %s: %w", userID, err)
	}

	// Step 2: Get all users who were referred using that email
	query := `
		SELECT user_id, email, first_name, last_name, phone,
			   cognito_id, referral_mail, role,
			   address, parent_id, created_at, updated_at
		FROM z_users
		WHERE referral_mail = $1;
	`

	rows, err := r.db.Query(ctx, query, referrerEmail)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		// Removed local *string variables like firstName, lastName etc.
		// Scanning directly into &user.Field, assuming domain.User fields are *string if nullable,
		// or string if not nullable (and DB ensures NOT NULL).
		var addressJSON []byte // For the JSONB 'address' column
		var parentID *string   // For nullable parent_id

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.CognitoID,
			&user.ReferralMail,
			&user.Role,
			&addressJSON,
			&parentID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}

		user.ParentID = parentID // Assign the scanned nullable parentID

		// --- CORRECTED ADDRESS PARSING LOGIC ---
		// This block is now structurally correct.
		if len(addressJSON) > 0 && string(addressJSON) != "null" {
			user.Address = &domain.Address{} // Initialize Address struct only if there's data to unmarshal
			if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
				// If unmarshaling fails, return an error for this row.
				return nil, fmt.Errorf("failed to parse address JSON in ListReferredUsers: %w", err)
			}
		} else {
			// If addressJSON is empty or "null", explicitly set user.Address to nil
			user.Address = nil
		}
		// --- END CORRECTED ADDRESS PARSING LOGIC ---

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}


// UpdateUserParentID updates the parent_id field for a specific user in the z_users table.
// It accepts a *string for parentID to correctly handle NULL values in the database.
func (r *UserRepository) UpdateUserParentID(ctx context.Context, userID string, parentID *string) error {
	// The query to update the parent_id and updated_at timestamp for a given user.
	 log.Println("ERROR: r.db is nil in UpdateUserParentID!")
	// Assuming 'id' is the primary key column for the user in z_users.
	  if r.db == nil {
        log.Println("ERROR: r.db is nil in UpdateUserParentID!")
        return fmt.Errorf("database connection is not initialized")
    }
	query := `UPDATE z_users SET parent_id = $1, updated_at = NOW() WHERE user_id = $2`

	result, err := r.db.Exec(ctx, query, parentID, userID)
	if err != nil {
		return fmt.Errorf("error updating user parent_id for user %s: %w", userID, err)
	}

	// Check how many rows were affected to verify the update occurred.
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Warning: No user found with ID %s to update parent_id", userID)
		// Depending on your business logic, you might want to return an error here
		// if it's critical that the user exists for the update to succeed.
	}
	return nil
}