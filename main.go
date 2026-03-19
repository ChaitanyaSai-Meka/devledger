package main

import (
    "log"
	"github.com/ChaitanyaSai-Meka/devledger/db"
)

func main() {
    conn, err := db.ConnectDB()
    if err != nil {
        log.SetFlags(0)
        log.Fatalf("Error: failed to connect to database: %v", err)
    }
    defer conn.Close()
}
