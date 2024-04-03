package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Conn *sql.DB = nil

func GetConn() (db *sql.DB, err error) {
	if Conn == nil {
		dbHost := "localhost"
		dbName := "wallet"
		dbUser := "wallet"
		dbPassword := "password"
		conn_info := fmt.Sprintf(
			"host=%s dbname=%s user=%s password=%s",
			dbHost,
			dbName,
			dbUser,
			dbPassword,
		)

		Conn, err = sql.Open("postgres", conn_info)
		if err != nil {
			log.Fatal("Connection Error:", err)
			return nil, err
		}
	}
	return Conn, err
}
