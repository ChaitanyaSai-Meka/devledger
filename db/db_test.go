package db

import "testing"

func TestConnectDB(t *testing.T) {
    conn, err := ConnectDB()
    if err != nil {
        t.Fatalf("ConnectDB failed: %v", err)
    }
    defer conn.Close()
}