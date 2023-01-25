package app

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"forum/internal/delivery"
	"forum/internal/service"
	"forum/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	cfg *config.Config

	db *sql.DB

	httpServer  *http.Server
	httpHandler *delivery.Handler
}

func InitApp(cfg *config.Config) (*app, error) {
	log.Println("init app")

	db, err := storage.InitDB(cfg)
	if err != nil {
		return nil, err
	}

	storages := storage.NewStorage(db)
	services := service.NewService(storages, cfg)

	handlers := delivery.NewHandler(services)

	return &app{
		cfg:         cfg,
		db:          db,
		httpHandler: handlers,
	}, nil
}

func (a *app) RunApp() {
	log.Println("run app")

	go func() {
		if err := a.startHTTP(); err != nil {
			log.Println(err)
			return
		}
	}()
	log.Println("http server started on", a.cfg.Http.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Println()
	log.Println("Received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Println(err)
		return
	}

	if err := a.db.Close(); err != nil {
		log.Println(err)
	} else {
		log.Println("db closed")
	}
}

func (a *app) startHTTP() error {
	router := a.httpHandler.InitRoutes()

	a.httpServer = &http.Server{
		Handler:      router,
		Addr:         ":" + a.cfg.Http.Port,
		ReadTimeout:  time.Second * time.Duration(a.cfg.Http.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(a.cfg.Http.WriteTimeout),
	}

	return a.httpServer.ListenAndServe()
}
