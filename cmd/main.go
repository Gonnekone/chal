package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gonnekone/challenge/internal/config"
	"github.com/Gonnekone/challenge/internal/http-server/handlers/del"
	"github.com/Gonnekone/challenge/internal/http-server/handlers/get"
	"github.com/Gonnekone/challenge/internal/http-server/handlers/save"
	"github.com/Gonnekone/challenge/internal/http-server/handlers/update"
	mwLogger "github.com/Gonnekone/challenge/internal/http-server/middleware/logger"
	"github.com/Gonnekone/challenge/internal/lib/logger/handlers/slogpretty"
	"github.com/Gonnekone/challenge/internal/lib/logger/sl"
	"github.com/Gonnekone/challenge/internal/storage/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/swaggo/http-swagger/v2"
	_ "github.com/Gonnekone/challenge/docs"
)

// @title			Cars catalog API
// @version		1.0
// @description	This is a testovoe zadanie.

// @contact.url    https://t.me/Gonnekone
// @contact.email	opacha2018@yandex.ru

// @host 			localhost:8082
// @schemes 		http
func main() {
	cfg := config.MustLoad()

	log := setupPrettySlog()
	log.Info("starting up the application")

	storage, err := postgres.New(log, cfg.DbURL)
	if err != nil {
		log.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	router.Post("/", save.New(log, storage))
	router.Get("/", get.New(log, storage))
	router.Patch("/", update.New(log, storage))
	router.Delete("/", del.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server {
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
