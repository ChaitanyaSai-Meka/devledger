package db

import (
	"fmt"
	"database/sql"
	_ "embed"
	_ "modernc.org/sqlite"

)

//go:embed schema.sql
var schemaSQL string

func ConnectDB()(*sql.DB, error){
	conn, err := sql.Open("sqlite","./devledger.db")

	if err != nil{
		return nil,err
	}

	_,err=conn.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		conn.Close()
		return nil,err
	}
	_, err = conn.Exec(schemaSQL)
	if err != nil {
		conn.Close()
		return nil,err
	}

	fmt.Println("Database connected and schema created successfully.")
	return conn,nil
}