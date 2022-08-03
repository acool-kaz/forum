package main

import (
	"fmt"
	"forum/internal/delivery"
	"forum/internal/service"
	"forum/internal/storage"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := storage.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	storages := storage.NewStorage(db)
	services := service.NewService(storages)
	handlers := delivery.NewHandler(services)
	handlers.InitRoutes()

	server := http.Server{
		Addr:    ":8080",
		Handler: handlers.Mux,
	}

	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
