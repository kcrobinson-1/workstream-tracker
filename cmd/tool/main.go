// Command tool is the workstream-tracker local server. See
// design/v0.1-design.md for the full design and spec/ for the
// plan-doc structure the tool consumes.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kcrobinson-1/workstream-tracker/internal/api"
	"github.com/kcrobinson-1/workstream-tracker/internal/db"
	"github.com/kcrobinson-1/workstream-tracker/internal/site"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "workstream-tracker.db"
	}

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := db.Init(initCtx, database); err != nil {
		initCancel()
		log.Fatalf("init db: %v", err)
	}
	initCancel()

	apiServer := api.New(database)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/work-instances", apiServer.MountRoutes)
	r.Mount("/", site.Router())

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Println("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	log.Printf("workstream-tracker listening on %s (db: %s)", addr, dbPath)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
