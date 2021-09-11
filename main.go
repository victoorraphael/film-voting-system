package main

import (
	"github.com/victoorraphael/film-voting-system/db"
	"log"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	webserviceStart()
}
