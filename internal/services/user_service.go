package services

import (
	"context"
	"fmt"
	"log"

	"github.com/afreedicp/zolaris-backend-app/internal/domain"
	"github.com/afreedicp/zolaris-backend-app/internal/repositories"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/mappers"
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserIdByCognitoId(ctx context.Context, cId string) (string, error) {
	// Corrected line 25:
	// Capture the values from the repository call first
	userID, err := s.userRepo.GetUserIdByCognitoId(ctx, cId)

	// Now you can use them in the log statement.
	// You need to decide what you want to log for the second %s.
	// It's usually the actual ID or an error message.
	// If userID is empty, it means not found, which is a success from repo's perspective.
	// If err is not nil, that's an actual error from the DB.
	if err != nil {
		log.Printf("Error getting user ID by Cognito ID %s: %v", cId, err)
		return "", fmt.Errorf("error retrieving user ID by Cognito ID: %w", err)
	}
	log.Printf("Getting user ID by Cognito ID: %s, Result: %s", cId, userID) // Changed the log message
	
	// Then return the results
	return userID, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	log.Printf("Getting user details for user %s", userID)
	return s.userRepo.GetUserByID(ctx, userID)
}

// CreateUser creates a new user account
func (s *UserService) CreateUser(ctx context.Context, req *dto.UserDetailsRequest) (*domain.User, error) {
	// Convert DTO to domain entity
	log.Printf("UserRequestToEntity")
	user := mappers.UserRequestToEntity(req, nil)

	// Save user to database
	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUserDetails updates a user's details
func (s *UserService) UpdateUserDetails(ctx context.Context, userID string, req *dto.UserDetailsRequest) (*domain.User, error) {
	// Get existing user
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if existingUser == nil {
		return nil, fmt.Errorf("user not found with ID: %s", userID)
	}

	// Update user with new details
	updatedUser := mappers.UserRequestToEntity(req, existingUser)

	// Save updated user to database
	err = s.userRepo.UpdateUser(ctx, updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// CheckHasParentID checks if a user has a parent ID
func (s *UserService) CheckHasParentID(ctx context.Context, userID string) (bool, error) {
	return s.userRepo.CheckHasParentID(ctx, userID)
}

func (s *UserService) ListReferredUsers(ctx context.Context, userID string) ([]*domain.User, error) {
	log.Printf("Listing referred users for user %s", userID)
	return s.userRepo.ListReferredUsers(ctx, userID)
}
