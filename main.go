package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ChaitanyaSai-Meka/devledger/api"
	"github.com/ChaitanyaSai-Meka/devledger/cli"
	"github.com/ChaitanyaSai-Meka/devledger/db"
)

func main() {
	log.SetFlags(0)

	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error: failed to connect to database: %v", err)
	}
	defer conn.Close()

	router := api.SetupRouter(conn)
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
