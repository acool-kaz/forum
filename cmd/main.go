package main

import (
	"forum/internal/app"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app.Run()
}
