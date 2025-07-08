package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/afreedicp/zolaris-backend-app/internal/domain"
)

// SensorDataDBModel represents how sensor readings are stored in the database
type SensorDataDBModel struct {
	Timestamp   int64  `dynamodbav:"timestamp"`
	Amperage    string `dynamodbav:"amperage"`
	Temperature string `dynamodbav:"temperature"`
	Humidity    string `dynamodbav:"humidity"`
}

// DeviceRepository handles all device-related database operations
type DeviceRepository struct {
	pgPool       *pgxpool.Pool    // PostgreSQL connection pool for device data
	dynamoClient *dynamodb.Client // DynamoDB client for sensor data
	machineTable string           // DynamoDB table for sensor readings
}

// NewDeviceRepository creates a new device repository instance
func NewDeviceRepository(pgPool *pgxpool.Pool, dynamoClient *dynamodb.Client) *DeviceRepository {
	return &DeviceRepository{
		pgPool:       pgPool,
		dynamoClient: dynamoClient,
		machineTable: "machine_data_table",
	}
}

// WithMachineTable sets the machine data table name for the repository
func (r *DeviceRepository) WithMachineTable(machineTable string) *DeviceRepository {
	r.machineTable = machineTable
	return r
}

// AddDevice adds a new device to the PostgreSQL database
func (r *DeviceRepository) AddDevice(ctx context.Context, deviceID, deviceName, userID string) error {
	query := `
		INSERT INTO z_device (
			mac_address, user_id, device_name, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (mac_address) DO UPDATE SET
			device_name = $3,
			updated_at = $5
	`

	now := time.Now()
	_, err := r.pgPool.Exec(
		ctx,
		query,
		deviceID,
		userID,
		deviceName,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to add device: %w", err)
	}

	return nil
}

// GetDevicesByUserID retrieves all devices for a specific user from PostgreSQL
func (r *DeviceRepository) GetDevicesByUserID(ctx context.Context, userID string) ([]*domain.Device, error) {
	query := `
		SELECT mac_address, user_id, device_name, category, description, created_at, updated_at
		FROM z_device
		WHERE user_id = $1
		ORDER BY device_name
	`

	rows, err := r.pgPool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var devices []*domain.Device
	for rows.Next() {
		device := &domain.Device{}
		err := rows.Scan(
			&device.MacAddress,
			&device.UserID,
			&device.Name,
			&device.Category,
			&device.Description,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device row: %w", err)
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device rows: %w", err)
	}

	return devices, nil
}

// GetSensorData retrieves sensor data from DynamoDB for a specific device within a time range
func (r *DeviceRepository) GetSensorData(ctx context.Context, macID string, startTime, endTime int64) ([]*domain.SensorReading, error) {
	log.Printf("Table name: %s", r.machineTable)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.machineTable),
		KeyConditionExpression: aws.String("mac_id = :macId AND #ts BETWEEN :startTime AND :endTime"),
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":macId":     &types.AttributeValueMemberS{Value: macID},
			":startTime": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", startTime)},
			":endTime":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", endTime)},
		},
	}

	log.Printf("Querying sensor data for device %s from %d to %d", macID, startTime, endTime)

	result, err := r.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var dbSensorData []SensorDataDBModel
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dbSensorData)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	domainReadings := make([]*domain.SensorReading, len(dbSensorData))
	for i, reading := range dbSensorData {
		// Convert timestamp from milliseconds to time.Time
		timestamp := time.UnixMilli(reading.Timestamp)

		domainReadings[i] = &domain.SensorReading{
			DeviceID:    macID,
			Timestamp:   timestamp,
			Amperage:    reading.Amperage,
			Temperature: reading.Temperature,
			Humidity:    reading.Humidity,
			RawData:     "", // We don't have this in DB currently
		}
	}

	return domainReadings, nil
}
