package main

import (
	"ForumWeb/internal/config"
	"ForumWeb/internal/db"
	"ForumWeb/internal/server"
	"log"
)

func main() {
	// load config
	config.LoadEnv()

	// load db
	db.InitDB()
	defer db.CloseDB()

	// load server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	// gracefull shutdown
}
