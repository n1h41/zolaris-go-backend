package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the system
type User struct {
	ID        string    `json:"id" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	FirstName *string   `json:"firstName" db:"first_name"`
	LastName  *string   `json:"lastName" db:"last_name"`
	Phone     *string   `json:"phone" db:"phone"`
	CognitoID *string   `json:"cognitoId,omitempty" db:"cognito_id"`
	ReferralMail *string   `json:"referralMail,omitempty" db:"referral_mail"`
	Role 		*string   `json:"role,omitempty" db:"role"`
	Address   *Address  `json:"address"`
	ParentID  *string   `json:"parentId,omitempty" db:"parent_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Address represents a physical address
type Address struct {
	Street1 string `json:"street1" db:"street1"`
	Street2 string `json:"street2,omitempty" db:"street2"`
	City    string `json:"city" db:"city"`
	Region  string `json:"region" db:"region"`
	Country string `json:"country" db:"country"`
	Zip     string `json:"zip" db:"zip"`
}

// Device represents an IoT device entity
type Device struct {
	MacAddress  string    `json:"macAddress" db:"mac_address"`
	UserID      string    `json:"userId" db:"user_id"`
	Name        string    `json:"name" db:"device_name"`
	Category    *string   `json:"category,omitempty" db:"category"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// SensorReading represents data from a device sensor
type SensorReading struct {
	DeviceID    string    `json:"deviceId" db:"mac_id"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Amperage    string    `json:"amperage" db:"amperage"`
	Temperature string    `json:"temperature" db:"temperature"`
	Humidity    string    `json:"humidity" db:"humidity"`
	RawData     string    `json:"-" db:"raw_data"`
}

// Category represents a device category
type Category struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// NewUser creates a new User with default values
func NewUser(email, firstName, lastName, phone string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		FirstName: &firstName,
		LastName:  &lastName,
		Phone:     &phone,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewDevice creates a new Device with default values
func NewDevice(macAddress, userID, name string) *Device {
	now := time.Now()
	return &Device{
		MacAddress: macAddress,
		UserID:     userID,
		Name:       name,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Entity represents a hierarchical entity in the system
type Entity struct {
	ID         string          `json:"id" db:"entity_id"`
	UserID     *string         `json:"userId,omitempty" db:"user_id"`
	Name       string          `json:"name" db:"name"`
	Details    json.RawMessage `json:"details" db:"details"`
	CategoryID string          `json:"categoryId" db:"category_id"`
	ParentID   *string         `json:"parentId,omitempty" db:"parent_id"`
	Path       string          `json:"path,omitempty" db:"path"`
	Depth      int             `json:"depth" db:"depth"`
	CreatedAt  time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time       `json:"updatedAt" db:"updated_at"`
}

// NewCategory creates a new Category with default values
func NewCategory(name, categoryType string) *Category {
	return &Category{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      categoryType,
		CreatedAt: time.Now(),
	}
}
