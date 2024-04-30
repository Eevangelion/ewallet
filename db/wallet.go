package db

import (
	"log"

	"github.com/google/uuid"
)

func Create(balance float32) (id string, err error) {
	db, err := GetConn()
	if err != nil {
		log.Println("Connection error:", err)
		return
	}
	qry := `INSERT INTO public."wallet" VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(qry, uuid.New(), balance).Scan(&id)
	if err != nil {
		log.Println("Error while trying to create wallet:", err)
		return
	}
	return
}
