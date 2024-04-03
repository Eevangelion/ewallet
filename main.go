package main

import (
	"fmt"
	"log"

	"github.com/Eevangelion/ewallet/config"
	"github.com/Eevangelion/ewallet/db"
)

func main() {
	conf := config.GetConfig()
	port := conf.Server.Port
	r := GetRouter()
	conn, err := db.GetConn()
	if err != nil {
		log.Print("Error while connecting to DB:", err.Error())
	}
	defer conn.Close()
	r.Run(fmt.Sprintf(":%d", port))
}
