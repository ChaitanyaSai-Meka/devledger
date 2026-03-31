package main

import (
	"log"
	"net"
	"net/http"

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

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error: failed to bind port 8080: %v", err)
	}

	go func() {
		if err := http.Serve(listener, router); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
