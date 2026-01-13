package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikunj/rest-api/internal/config"
	"github.com/nikunj/rest-api/internal/config/http/handlers/student"
	"github.com/nikunj/rest-api/internal/storage/sqlite"
)

func main() {

	//config set up
	cfg := config.MustLoad()
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized ", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	//database setup
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	//w is response  r is request

	//set up server

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	fmt.Printf("server started %s", cfg.HTTPServer.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("error in server")
		}
	}()
	<-done
	slog.Info("shutting down")

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancle()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Info("failed to shutting down the server", slog.String("error", err.Error()))
	}
	slog.Info("shutting down successfully")
}
