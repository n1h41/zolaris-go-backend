package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/afreedicp/zolaris-backend-app/internal/config"
)

// Database holds the database clients and configuration
type Database struct {
	dynamoClient     *dynamodb.Client
	postgresPool     *pgxpool.Pool
	deviceTable      string
	machineDataTable string
}

// NewDatabase creates and initializes database clients
func NewDatabase(ctx context.Context, dynamoClient *dynamodb.Client, cfg *config.Config) (*Database, error) {
	// Initialize PostgreSQL
	pgDB, err := NewPostgresDB(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Database{
		dynamoClient:     dynamoClient,
		postgresPool:     pgDB.GetPool(),
		deviceTable:      cfg.Database.DeviceTableName,
		machineDataTable: cfg.Database.DataTableName,
	}, nil
}

// GetDynamoClient returns the DynamoDB client
func (db *Database) GetDynamoClient() *dynamodb.Client {
	return db.dynamoClient
}

// GetPostgresPool returns the PostgreSQL connection pool
func (db *Database) GetPostgresPool() *pgxpool.Pool {
	return db.postgresPool
}

// GetDeviceTableName returns the device table name
func (db *Database) GetDeviceTableName() string {
	return db.deviceTable
}

// GetMachineDataTableName returns the machine data table name
func (db *Database) GetMachineDataTableName() string {
	return db.machineDataTable
}

// Close closes all database connections
func (db *Database) Close() {
	// Close PostgreSQL connection if it's initialized
	if db.postgresPool != nil {
		db.postgresPool.Close()
	}
}
