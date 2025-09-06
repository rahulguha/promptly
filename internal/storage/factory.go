package storage

import (
	"fmt"

	"github.com/rahulguha/promptly/internal/storage/jsonstore"
	"github.com/rahulguha/promptly/internal/storage/sqlite"
)

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeJSON   StorageType = "json"
	StorageTypeSQLite StorageType = "sqlite"
)

// StorageConfig holds configuration for storage backends
type StorageConfig struct {
	Type     StorageType
	JSONPath string  // Path to JSON file (for JSON storage)
	DBPath   string  // Path to SQLite database file (for SQLite storage)
}

// NewStorage creates a new storage instance based on the provided configuration
func NewStorage(config StorageConfig) (Storage, error) {
	switch config.Type {
	case StorageTypeJSON:
		if config.JSONPath == "" {
			return nil, fmt.Errorf("JSON path is required for JSON storage")
		}
		return jsonstore.NewFileStorage(config.JSONPath)
	
	case StorageTypeSQLite:
		if config.DBPath == "" {
			return nil, fmt.Errorf("database path is required for SQLite storage")
		}
		return sqlite.NewSQLiteStorage(config.DBPath)
	
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}

// NewProfileStorage creates a new profile storage instance based on the provided configuration
func NewProfileStorage(config StorageConfig) (ProfileStorage, error) {
	switch config.Type {
	case StorageTypeSQLite:
		if config.DBPath == "" {
			return nil, fmt.Errorf("database path is required for SQLite storage")
		}
		return sqlite.NewSQLiteStorage(config.DBPath)
	default:
		return nil, fmt.Errorf("unsupported storage type for profiles: %s", config.Type)
	}
}

// ValidateStorageType checks if the provided storage type is valid
func ValidateStorageType(storageType string) (StorageType, error) {
	switch StorageType(storageType) {
	case StorageTypeJSON, StorageTypeSQLite:
		return StorageType(storageType), nil
	default:
		return "", fmt.Errorf("invalid storage type '%s', must be 'json' or 'sqlite'", storageType)
	}
}