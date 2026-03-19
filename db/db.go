package db

import (
	"fmt"
	_ "embed"
	"database/sql"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

func ConnectDB()(*sql.DB, error){
	db, err := sql.Open("sqlite","../devledger.db")

	if err != nil{
		return nil,err
	}

	_,err=db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil,err
	}
	_, err = db.Exec(schemaSQL)
	if err != nil {
		return nil,err
	}

	fmt.Println("Database connected and schema created successfully.")
	return db,nil
}