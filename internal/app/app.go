package app

import (
	"context"
	"forum/internal/delivery"
	"forum/internal/server"
	"forum/internal/service"
	"forum/internal/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
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

	server := new(server.Server)
	go func() {
		if err := server.Start(":8080", handlers.InitRoutes()); err != nil {
			log.Println(err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Println(err)
		return
	}
}
