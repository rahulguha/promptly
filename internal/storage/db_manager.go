package storage

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/rahulguha/promptly/internal/storage/sqlite"
	_ "modernc.org/sqlite"
)

// DBManager handles a pool of database connections, one for each user.
type DBManager struct {
	mu  sync.RWMutex
	dbs map[string]*sql.DB
}

// NewDBManager creates a new DBManager.
func NewDBManager() *DBManager {
	return &DBManager{dbs: make(map[string]*sql.DB)}
}

// GetDB returns a database connection for a given user.
// If a connection for the user does not exist, it creates a new one.
func (m *DBManager) GetDB(userID, email string) (*sql.DB, error) {
	key := fmt.Sprintf("%s-%s-promptly", userID, email)

	m.mu.RLock()
	db, ok := m.dbs[key]
	m.mu.RUnlock()

	if ok {
		return db, nil
	}

	// Open new connection if it doesn't exist in the map
	dbPath := fmt.Sprintf("./data/%s.db", key)
	newDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Initialize schema for the new DB
	if err := sqlite.InitializeSchema(newDB); err != nil {
		newDB.Close()
		return nil, fmt.Errorf("failed to initialize schema for new db: %w", err)
	}

	m.mu.Lock()
	m.dbs[key] = newDB
	m.mu.Unlock()

	return newDB, nil
}
