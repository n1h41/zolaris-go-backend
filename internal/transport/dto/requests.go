package dto

import "time"

// UserDetailsRequest represents a request to create or update user details
type UserDetailsRequest struct {
	Email        string `json:"email" validate:"required,email"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Street1      string `json:"street1" validate:"required"`
	Street2      string `json:"street2"`
	City         string `json:"city" validate:"required"`
	Region       string `json:"region" validate:"required"`
	Country      string `json:"country" validate:"required"`
	Zip          string `json:"zip" validate:"required"`
	ParentID     string `json:"parentId,omitempty"`
	Role         string  `json:"role"`
	ReferralMail string `json:"referralMail,omitempty"`
	ID string `json:"-"`
	 CognitoID    string `json:"cognitoId" validate:"required"`
}

// DeviceRequest represents a request to add a new device
type DeviceRequest struct {
	DeviceID    string `json:"deviceId" validate:"required,min=3,max=50"`
	DeviceName  string `json:"deviceName" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
}

// CategoryRequest represents a request to add a new category
type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
	Type string `json:"type" validate:"required,min=2,max=50"`
}

// PolicyAttachRequest represents a request to attach an IoT policy
type PolicyAttachRequest struct {
	IdentityID string `json:"identityId" validate:"required"`
}

// SensorDataRequest represents a request to get device sensor data
type SensorDataRequest struct {
	DeviceMacID string `json:"deviceMacId" validate:"required"`
	Timestamp   string `json:"timestamp" validate:"required"`
	DateMode    string `json:"dateMode" validate:"required,oneof=hourly daily weekly monthly yearly"`
}

// TimeRange defines start and end times for data filtering
type TimeRange struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// PaginationParams contains pagination parameters
type PaginationParams struct {
	Page     int `json:"page" form:"page" default:"0"`
	PageSize int `json:"pageSize" form:"pageSize" default:"20"`
}

// CreateRootEntityRequest represents a request to create a root entity
type CreateRootEntityRequest struct {
	CategoryID string         `json:"categoryId" validate:"required,uuid"`
	Name       string         `json:"name" validate:"required,min=2,max=100"`
	Details    map[string]any `json:"details,omitempty"`
}

// CreateSubEntityRequest represents a request to create a child entity
type CreateSubEntityRequest struct {
	CategoryID     string         `json:"categoryId" validate:"required,uuid"`
	Name           string         `json:"name" validate:"required,min=2,max=100"`
	Details        map[string]any `json:"details,omitempty"`
	ParentEntityID string         `json:"parentEntityId,omitempty" validate:"omitempty,uuid"`
}

// GetEntityChildrenRequest represents a request to get children of an entity
type GetEntityChildrenRequest struct {
	Recursive    bool   `json:"recursive" form:"recursive" default:"false"`
	Level        int    `json:"level" form:"level" default:"0"`
	CategoryType string `json:"categoryType" form:"categoryType"`
}

// GetEntityHierarchyRequest represents a request to get an entity hierarchy
type GetEntityHierarchyRequest struct {
	MaxDepth int `json:"maxDepth" form:"maxDepth" default:"10"`
}
