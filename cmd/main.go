package main

import (
	"forum/internal/app"
	"forum/internal/config"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.InitConfig("./config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.InitApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.RunApp()
}
