package main

import (
	"fmt"
	"log"

	"github.com/Eevangelion/ewallet/config"
	"github.com/Eevangelion/ewallet/db"
	"github.com/Eevangelion/ewallet/server"
)

func main() {
	conf := config.GetConfig()
	port := conf.Server.Port
	r := server.GetRouter()
	conn, err := db.GetConn()
	if err != nil {
		log.Print("Error while connecting to DB:", err.Error())
	}
	defer conn.Close()
	r.Run(fmt.Sprintf(":%d", port))
}
