package main

import (
	"github.com/ChaitanyaSai-Meka/devledger/db"
)

func main() {
    db, err := db.ConnectDB()
    if err != nil {
        panic(err)
    }
    defer db.Close()
}
