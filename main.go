package main

import (
	"go1f/pkg/server"
	"log"

	"github.com/Vorobey112/go-final/pkg/db"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetDB().Close()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
