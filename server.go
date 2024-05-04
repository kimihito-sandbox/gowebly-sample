package main

import (
	"fmt"
	"embed"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	gowebly "github.com/gowebly/helpers"
)

//go:embed all:static
var static embed.FS

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() error {
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "7777"))
	if err != nil {
		return err
	}

	// Create a new chi router.
	router := chi.NewRouter()

	// Use chi middlewares.
	router.Use(middleware.Logger)

	// Handle static files from the embed FS (with a custom handler).
	router.Handle("/static/*", gowebly.StaticFileServerHandler(http.FS(static)))

	// Handle index page view.
	router.Get("/", indexViewHandler)

	// Handle API endpoints.
	router.Get("/api/hello-world", showContentAPIHandler)

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router, // handle all chi routes
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Send log message.
	slog.Info("Starting server...", "port", port)

	return server.ListenAndServe()
}
