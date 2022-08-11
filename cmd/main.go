package main

import (
	"context"
	"fmt"
	"forum/internal/delivery"
	"forum/internal/service"
	"forum/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()
	Run(ctx)
}

func Run(ctx context.Context) {
	db, err := storage.InitDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("can't close db err: %v\n", err)
		} else {
			log.Printf("db closed")
		}
	}()

	if err := storage.CreateTables(db); err != nil {
		log.Println(err)
		return
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
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			cancel()
			return
		}
	}()
	<-ctx.Done()
	ctx, cancel = context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Println(err)
		return
	}
}
