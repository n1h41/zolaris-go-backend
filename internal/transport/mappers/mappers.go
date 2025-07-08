package mappers

import (
	"encoding/json"
	"time"
"log"
	"github.com/afreedicp/zolaris-backend-app/internal/domain"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
)

// UserToResponse converts a domain User to a UserResponse DTO
func UserToResponse(user *domain.User) *dto.UserResponse {
	if user == nil {
		return nil
	}

	var response *dto.UserResponse

	if user.Address == nil {
		user.Address = &domain.Address{
			Street1: "",
			Street2: "",
			City:    "",
			Region:  "",
			Country: "",
			Zip:     "",
		}
	}

	response = &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Address: &dto.AddressOutput{
			Street1: user.Address.Street1,
			Street2: user.Address.Street2,
			City:    user.Address.City,
			Region:  user.Address.Region,
			Country: user.Address.Country,
			Zip:     user.Address.Zip,
		},
		CreatedAt: user.CreatedAt,
	}

	if user.ParentID != nil {
		response.ParentID = *user.ParentID
	}

	return response
}

// UserRequestToEntity converts a UserDetailsRequest to a domain User entity
func UserRequestToEntity(req *dto.UserDetailsRequest, existingUser *domain.User) *domain.User {
	var user *domain.User

	if existingUser != nil {
		user = existingUser
		user.UpdatedAt = time.Now()
	}else {
        // Option A: Update NewUser to take more parameters
        // This is generally cleaner if NewUser is meant to fully construct a user.
        // You'd need to modify domain.NewUser to accept cognitoID, referralMail, role.
        // Example: user = domain.NewUser(req.Email, req.FirstName, req.LastName, req.Phone, req.CognitoID, req.ReferralMail, req.Role)
        // And then in domain.NewUser:
        // func NewUser(email string, fn, ln, ph, cID, rm, r string) *User { ... return &User{..., CognitoID: &cID, Role: &r, ...} }

        // Option B: Create a new domain.User struct and populate fields here (more flexible for partial updates)
        // This is what I'll assume for now, as it's easier to integrate incrementally.
        // If domain.NewUser still creates the base user, assign the rest here.
        user = domain.NewUser(req.Email, req.FirstName, req.LastName, req.Phone) // Keep this call if it generates ID/timestamps

        // Set mandatory fields from DTO that NewUser might not handle
        if req.CognitoID != "" { // Assuming CognitoID is now in dto.UserDetailsRequest and validate:"required"
            user.CognitoID = &req.CognitoID
        } else {
            // This should ideally not happen if validate:"required" works correctly on CognitoID in the DTO
            log.Printf("Warning: CognitoID was empty in request DTO, but is required for user creation. This might lead to DB error if not handled upstream.")
        }
    }

	user.Email = req.Email
	user.FirstName = &req.FirstName
	user.LastName = &req.LastName
	user.Phone = &req.Phone
	user.Address = &domain.Address{
		Street1: req.Street1,
		Street2: req.Street2,
		City:    req.City,
		Region:  req.Region,
		Country: req.Country,
		Zip:     req.Zip,
	}

	if req.ParentID != "" {
		user.ParentID = &req.ParentID
	}

	  // Map Role (assuming it's now in dto.UserDetailsRequest and not required, so has a default logic)
    if req.Role != "" {
        user.Role = &req.Role
    } else {
        defaultRole := "user" // Default role if not provided in DTO
        user.Role = &defaultRole
    }

    // Map ReferralMail (assuming it's now in dto.UserDetailsRequest and optional)
    if req.ReferralMail != "" {
        user.ReferralMail = &req.ReferralMail
    }
    // No else needed for ReferralMail if it's genuinely optional and can be nil in DB

	return user
}

// DeviceToResponse converts a domain Device to a DeviceResponse DTO
func DeviceToResponse(device *domain.Device) *dto.DeviceResponse {
	if device == nil {
		return nil
	}

	response := &dto.DeviceResponse{
		DeviceID:   device.MacAddress,
		DeviceName: device.Name,
		CreatedAt:  device.CreatedAt,
	}

	if device.Category != nil {
		response.Category = *device.Category
	}

	if device.Description != nil {
		response.Description = *device.Description
	}

	return response
}

// DeviceRequestToEntity converts a DeviceRequest to a domain Device entity
func DeviceRequestToEntity(req *dto.DeviceRequest, userID string) *domain.Device {
	device := domain.NewDevice(req.DeviceID, userID, req.DeviceName)

	if req.Category != "" {
		device.Category = &req.Category
	}

	if req.Description != "" {
		device.Description = &req.Description
	}

	return device
}

// SensorReadingToResponse converts a domain SensorReading to a SensorDataResponse DTO
func SensorReadingToResponse(reading *domain.SensorReading) *dto.SensorDataResponse {
	if reading == nil {
		return nil
	}

	return &dto.SensorDataResponse{
		Timestamp:   reading.Timestamp.UnixMilli(),
		Amperage:    reading.Amperage,
		Temperature: reading.Temperature,
		Humidity:    reading.Humidity,
	}
}

// CategoryToResponse converts a domain Category to a CategoryResponse DTO
func CategoryToResponse(category *domain.Category) *dto.CategoryResponse {
	if category == nil {
		return nil
	}

	return &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Type: category.Type,
	}
}

// Batch conversion helpers
func UsersToResponses(users []*domain.User) []*dto.UserResponse {
	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserToResponse(user)
	}
	return responses
}

func DevicesToResponses(devices []*domain.Device) []*dto.DeviceResponse {
	responses := make([]*dto.DeviceResponse, len(devices))
	for i, device := range devices {
		responses[i] = DeviceToResponse(device)
	}
	return responses
}

func SensorReadingsToResponses(readings []*domain.SensorReading) []*dto.SensorDataResponse {
	responses := make([]*dto.SensorDataResponse, len(readings))
	for i, reading := range readings {
		responses[i] = SensorReadingToResponse(reading)
	}
	return responses
}

func CategoriesToResponses(categories []*domain.Category) []*dto.CategoryResponse {
	responses := make([]*dto.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryToResponse(category)
	}
	return responses
}

// EntityToResponse converts a domain Entity to an EntityResponse DTO
func EntityToResponse(entity *domain.Entity) *dto.EntityResponse {
	if entity == nil {
		return nil
	}

	response := &dto.EntityResponse{
		ID:         entity.ID,
		Name:       entity.Name,
		CategoryID: entity.CategoryID,
		Depth:      entity.Depth,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}

	if entity.UserID != nil {
		response.UserID = *entity.UserID
	}

	if entity.ParentID != nil {
		response.ParentID = *entity.ParentID
	}

	// Parse details if available
	if len(entity.Details) > 0 {
		var details map[string]any
		if err := json.Unmarshal(entity.Details, &details); err == nil {
			response.Details = details
		}
	}

	return response
}

// EntitiesToResponses converts a slice of domain Entity to EntityResponse DTOs
func EntitiesToResponses(entities []*domain.Entity) []*dto.EntityResponse {
	responses := make([]*dto.EntityResponse, len(entities))
	for i, entity := range entities {
		responses[i] = EntityToResponse(entity)
	}
	return responses
}

// HierarchyMapToResponse converts a map-based entity hierarchy to EntityHierarchyResponse
func HierarchyMapToResponse(entityMap map[string]any) *dto.EntityHierarchyResponse {
	if entityMap == nil {
		return nil
	}

	response := &dto.EntityHierarchyResponse{
		EntityResponse: dto.EntityResponse{
			ID:        entityMap["id"].(string),
			Name:      entityMap["name"].(string),
			Depth:     entityMap["depth"].(int),
			CreatedAt: entityMap["created_at"].(time.Time),
		},
	}

	// Handle optional fields
	if userID, ok := entityMap["user_id"].(string); ok {
		response.UserID = userID
	}

	if parentID, ok := entityMap["parent_id"].(string); ok {
		response.ParentID = parentID
	}

	if categoryID, ok := entityMap["category_id"].(string); ok {
		response.CategoryID = categoryID
	}

	if categoryName, ok := entityMap["category_name"].(string); ok {
		response.CategoryName = categoryName
	}

	if categoryType, ok := entityMap["category_type"].(string); ok {
		response.CategoryType = categoryType
	}

	if details, ok := entityMap["details"].(map[string]any); ok {
		response.Details = details
	}

	// Process children recursively
	if children, ok := entityMap["children"].([]map[string]any); ok && len(children) > 0 {
		response.Children = make([]dto.EntityHierarchyResponse, len(children))
		for i, child := range children {
			childResponse := HierarchyMapToResponse(child)
			if childResponse != nil {
				response.Children[i] = *childResponse
			}
		}
	}

	return response
}
