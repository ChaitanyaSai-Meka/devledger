package db

import (
	"os"
	"testing"
)

func TestConnectDB(t *testing.T) {
	conn, err := ConnectDB()
	if err != nil {
		t.Fatalf("ConnectDB failed: %v", err)
	}
	defer os.Remove("./devledger.db")
	defer conn.Close()
}

func TestForeignKeysEnabled(t *testing.T) {
	conn, err := ConnectDB()
	if err != nil {
		t.Fatalf("ConnectDB failed: %v", err)
	}
	defer os.Remove("./devledger.db")
	defer conn.Close()

	var enabled int
	err = conn.QueryRow("PRAGMA foreign_keys;").Scan(&enabled)
	if err != nil {
		t.Fatalf("failed to check foreign_keys pragma: %v", err)
	}

	if enabled != 1 {
		t.Fatalf("expected foreign_keys = 1, got %d", enabled)
	}
}

func TestUsersTableCreated(t *testing.T) {
	conn, err := ConnectDB()
	if err != nil {
		t.Fatalf("ConnectDB failed: %v", err)
	}
	defer os.Remove("./devledger.db")
	defer conn.Close()

	var count int
	query := `
		SELECT COUNT(*)
		FROM sqlite_master
		WHERE type='table' AND name='Users'
	`

	err = conn.QueryRow(query).Scan(&count)
	if err != nil {
		t.Fatalf("failed to check Users table: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected Users table to exist, got count = %d", count)
	}
}
