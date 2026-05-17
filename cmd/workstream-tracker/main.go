// Command tool is the workstream-tracker local server. See
// design/v0.1-design.md for the full design and spec/ for the
// plan-doc structure the tool consumes.
package main

import (
	"context"
	"errors"
	"log/slog"
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

// main dispatches subcommands. With no subcommand the binary runs
// the server (unchanged v0.1 behavior); the `register` subcommand
// performs one best-effort work-instance registration.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "register" {
		os.Exit(runRegister(os.Args[2:], os.Getenv, os.Stdout, os.Stderr))
	}
	runServer()
}

func runServer() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "workstream-tracker.db"
	}

	plansPath := os.Getenv("PLANS_PATH")
	if plansPath == "" {
		plansPath = "docs/plans"
	}

	database, err := db.Open(dbPath)
	if err != nil {
		slog.Error("open db", "path", dbPath, "err", err)
		os.Exit(1)
	}
	defer database.Close()

	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := db.Init(initCtx, database); err != nil {
		initCancel()
		slog.Error("init db", "path", dbPath, "err", err)
		os.Exit(1)
	}
	initCancel()

	apiServer := api.New(database)
	siteServer := site.New(database, plansPath)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", health)
	r.Route("/work-instances", apiServer.MountRoutes)
	r.Mount("/", siteServer.Router())

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		slog.Info("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	slog.Info("listening", "addr", addr, "db", dbPath, "plans", plansPath)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("listen and serve", "err", err)
		os.Exit(1)
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
