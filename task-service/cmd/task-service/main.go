package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"task-service/internal/config"
	"task-service/internal/http-server/handlers/task/delete"
	"task-service/internal/http-server/handlers/task/get"
	"task-service/internal/http-server/handlers/task/get_by_author"
	"task-service/internal/http-server/handlers/task/save"
	"task-service/internal/http-server/handlers/task/update"
	"task-service/internal/http-server/middleware/logger"
	"task-service/internal/lib/logger/sl"
	"task-service/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/task/", save.New(log, storage))
	router.Put("/task/", update.New(log, storage))
	router.Delete("/task/{id}", delete.New(log, storage))
	router.Get("/task/{id}", get.New(log, storage))
	router.Get("/task/author/{author}", get_by_author.New(log, storage))
	log.Info("Server starting...", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}
	log.Info("Server stopped", slog.String("address", cfg.Address))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
