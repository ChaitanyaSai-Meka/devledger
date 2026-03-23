package main

import (
	"github.com/ChaitanyaSai-Meka/devledger/db"
	"log"
)

func main() {
	log.SetFlags(0)

	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error: failed to connect to database: %v", err)
	}
	defer conn.Close()
}
