package database

import (
	"database/sql"
	"log"
	"time"
)

// DbConn : Database Connection
var DbConn *sql.DB

// SetupDatabase sets up a new database connection
func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "root:admin@tcp(192.168.99.100:3306)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}

	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}
