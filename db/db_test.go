package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "devledger.db")
	conn, err := connectDB(dbPath)
	if err != nil {
		t.Fatalf("connectDB failed: %v", err)
	}

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("failed to close test database: %v", err)
		}
	})

	return conn
}

func TestConnectDB(t *testing.T) {
	openTestDB(t)
}

func TestForeignKeysEnabled(t *testing.T) {
	conn := openTestDB(t)

	var enabled int
	err := conn.QueryRow("PRAGMA foreign_keys;").Scan(&enabled)
	if err != nil {
		t.Fatalf("failed to check foreign_keys pragma: %v", err)
	}

	if enabled != 1 {
		t.Fatalf("expected foreign_keys = 1, got %d", enabled)
	}
}

func TestSchemaTablesCreated(t *testing.T) {
	conn := openTestDB(t)

	expectedTables := []string{
		"Users",
		"Groups",
		"GroupMembers",
		"Expenses",
		"Splits",
	}

	query := `
		SELECT COUNT(*)
		FROM sqlite_master
		WHERE type='table' AND name=?
	`

	for _, tableName := range expectedTables {
		var count int
		err := conn.QueryRow(query, tableName).Scan(&count)
		if err != nil {
			t.Fatalf("failed to check %s table: %v", tableName, err)
		}

		if count != 1 {
			t.Fatalf("expected %s table to exist, got count = %d", tableName, count)
		}
	}
}
