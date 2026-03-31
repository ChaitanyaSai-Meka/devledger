package main

import (
	"context"
	"log"
	"net"
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
	server := &http.Server{Handler: router}
	serverErrCh := make(chan error, 1)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error: failed to bind port 8080: %v", err)
	}
	defer listener.Close()

	go func() {
		serverErrCh <- server.Serve(listener)
	}()

	cliErr := cli.Execute()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error: failed to shut down server cleanly: %v", err)
	}

	if err := <-serverErrCh; err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
	if cliErr != nil {
		log.Fatal(cliErr)
	}
}
