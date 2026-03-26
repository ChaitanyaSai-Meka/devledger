package main

import (
	"log"
	"net/http"

	"github.com/ChaitanyaSai-Meka/devledger/api"
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
	log.Println("Server is running on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error: failed to start server: %v", err)
	}
}
