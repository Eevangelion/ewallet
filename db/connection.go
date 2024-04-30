package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Eevangelion/ewallet/config"
	_ "github.com/lib/pq"
)

var Conn *sql.DB = nil

func GetConn() (db *sql.DB, err error) {
	if Conn == nil {
		conf := config.GetConfig()

		dbHost := conf.DB.Host
		dbPort := conf.DB.Port
		dbName := conf.DB.Name
		dbUser := conf.DB.User
		dbPassword := conf.DB.Password
		conn_info := fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			dbHost,
			dbPort,
			dbName,
			dbUser,
			dbPassword,
		)

		Conn, err = sql.Open("postgres", conn_info)
		if err != nil {
			log.Fatal("Connection Error:", err)
			return nil, err
		}
		err = CreateTables(Conn)
		if err != nil {
			return
		}
	}
	return Conn, err
}

func CreateTables(db *sql.DB) (err error) {
	qry := `begin;
	
	CREATE TABLE IF NOT EXISTS public."wallet"
	(
		id varchar(50),
		balance float
	);
	
	COMMIT;`

	_, err = db.Exec(qry)

	if err != nil {
		log.Fatal("Create DB tables error:", err)
	}

	return
}
