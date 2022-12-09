package app

import (
	"context"
	"forum/internal/config"
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

func Run(cfg *config.Config) {
	db, err := storage.InitDB(cfg)
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

	storages := storage.NewStorage(db)
	services := service.NewService(storages)
	handlers := delivery.NewHandler(services)

	server := new(server.Server)
	go func() {
		if err := server.Start(cfg, handlers.InitRoutes()); err != nil {
			log.Println(err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Minute)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Println(err)
		return
	}
}
