package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/afreedicp/zolaris-backend-app/internal/middleware"
	"github.com/afreedicp/zolaris-backend-app/internal/services"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
	"github.com/afreedicp/zolaris-backend-app/internal/transport/response"
	"github.com/afreedicp/zolaris-backend-app/internal/utils"
)

// AddDeviceHandler handles requests to add a new device
type AddDeviceHandler struct {
	deviceService *services.DeviceService
}

// NewAddDeviceHandler creates a new AddDeviceHandler
func NewAddDeviceHandler(deviceService *services.DeviceService) *AddDeviceHandler {
	return &AddDeviceHandler{deviceService: deviceService}
}

// HandleGin handles requests using Gin framework
// @Summary Add a new device
// @Description Register a new IoT device for the authenticated user
// @Tags Device Management
// @Accept json
// @Produce json
// @Param X-Cognito-ID header string true "Cognito ID"
// @Param device body dto.DeviceRequest true "Device information"
// @Success 201 {object} dto.Response "Device added successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /device/add [post]
func (h *AddDeviceHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request dto.DeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		response.ValidationErrors(c, utils.CreateDtoValidationErrors(validationErrs))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Call service to add device
	if err := h.deviceService.AddDevice(c.Request.Context(), request.DeviceID, request.DeviceName, userID); err != nil {
		log.Printf("Error adding device: %v", err)
		response.InternalError(c, "Failed to add device")
		return
	}

	response.Created(c, nil, "Device added successfully")
}

// ListUserDevicesHandler handles requests to list devices for a user
type ListUserDevicesHandler struct {
	deviceService *services.DeviceService
}

// NewListUserDevicesHandler creates a new ListUserDevicesHandler
func NewListUserDevicesHandler(deviceService *services.DeviceService) *ListUserDevicesHandler {
	return &ListUserDevicesHandler{deviceService: deviceService}
}

// HandleGin handles requests using Gin framework
// @Summary List user devices
// @Description Get all devices registered to the authenticated user
// @Tags Device Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} dto.DeviceResponse "List of user devices"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/devices [get]
func (h *ListUserDevicesHandler) HandleGin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Call service to get user devices
	devices, err := h.deviceService.GetUserDevices(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user devices: %v", err)
		response.InternalError(c, "Failed to retrieve user devices")
		return
	}

	response.OK(c, devices, "Devices retrieved successfully")
}

// GetDeviceSensorDataHandler handles requests to get sensor data for a device
type GetDeviceSensorDataHandler struct {
	deviceService *services.DeviceService
}

// NewGetDeviceSensorDataHandler creates a new GetDeviceSensorDataHandler
func NewGetDeviceSensorDataHandler(deviceService *services.DeviceService) *GetDeviceSensorDataHandler {
	return &GetDeviceSensorDataHandler{deviceService: deviceService}
}

// HandleGin handles requests using Gin framework
// @Summary Get device sensor data
// @Description Retrieve sensor data for a specific device with time filtering
// @Tags Device Data
// @Accept json
// @Produce json
// @Param request body dto.SensorDataRequest true "Request parameters"
// @Success 200 {object} dto.Response{data=[]dto.SensorDataResponse} "Sensor data for the device"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /device/sensor-data [post]
func (h *GetDeviceSensorDataHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request dto.SensorDataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		response.ValidationErrors(c, utils.CreateDtoValidationErrors(validationErrs))
		return
	}

	// Use timestamp directly from request
	baseTimestamp := request.Timestamp

	// Call service to get sensor data
	data, err := h.deviceService.GetDeviceSensorData(c.Request.Context(), request.DeviceMacID, request.DateMode, baseTimestamp)
	if err != nil {
		log.Printf("Error getting sensor data: %v", err)
		response.InternalError(c, "Failed to retrieve sensor data")
		return
	}

	// Use the data directly in the response
	response.OK(c, data, "Data retrieved successfully")
}
