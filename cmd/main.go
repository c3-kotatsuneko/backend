package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/c3-kotatsuneko/backend/internal/di"
	"github.com/c3-kotatsuneko/backend/pkg/config"
	"github.com/c3-kotatsuneko/backend/pkg/handler"
	"github.com/c3-kotatsuneko/backend/pkg/log"
	"github.com/c3-kotatsuneko/backend/pkg/middleware"
	"github.com/rs/cors"
)

func main() {
	config.LoadEnv()

	h := di.InitHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("GET     /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World")
	})

	mux.Handle("GET /ws/physics", handler.AppHandler(h.PhysicsHandler.Calculate()))
	mux.Handle("GET /ws/events", handler.AppHandler(h.EventHandler.ManageEvent()))

	c := cors.AllowAll()
	handler := middleware.Chain(mux, middleware.Context, c.Handler, middleware.Recover, middleware.Logger)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	go func() {
		log.Start()
		if err := server.ListenAndServe(); err != nil {
			slog.Error("server error", "error", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server error", "error", err.Error())
	}
}
