package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"microservice/internal"
	"microservice/internal/db"
	"microservice/internal/router"
	"microservice/routes/data"
	"microservice/routes/meters"
)

// constants for controlling server behaviour.
const (
	headerReadTimeout     = 10 * time.Second
	serverShutdownTimeout = 20 * time.Second
)

func main() {
	_ = internal.ParseConfiguration()       // error ignored as function always returns nil
	configuration := internal.Configuration // alias for the configuration to keep the code short

	err := db.Connect()
	if err != nil {
		slog.Error("unable to connect to the database", "error", err)
		os.Exit(1)
	}

	err = db.MigrateDatabase()
	if err != nil {
		slog.Error("failed to execute database migrations", "error", err)
		os.Exit(1)
	}

	err = db.LoadQueries()
	if err != nil {
		slog.Error("unable to load database queries", "error", err)
		os.Exit(1)
	}

	router, err := router.GenerateRouter()
	if err != nil {
		slog.Error("unable to generate router instance", "error", err)
		os.Exit(1)
	}

	meterAPI := router.Group("/meters")
	{
		meterAPI.GET("/", meters.All)
		meterAPI.PUT("/", meters.Create)
		meterAPI.GET("/:meterID", meters.GetSingle)
	}

	dataAPI := router.Group("/data")
	{
		dataAPI.POST("/:meterID", data.Import)
	}

	// create a http server to handle the requests
	server := http.Server{
		Addr:              net.JoinHostPort(configuration.GetString(internal.ConfigKey_Http_Host), configuration.GetString(internal.ConfigKey_Http_Port)), //nolint:lll
		Handler:           router.Handler(),
		ReadHeaderTimeout: headerReadTimeout,
	}

	// Start the server and log errors that happen while running it
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("unable to start http server", "error", err)
		}
	}()

	// Set up some the signal handling to allow the server to shut down gracefully
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	// Block further code execution until the shutdown signal was received
	<-shutdownSignal

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("unable to shutdown api gracefully", "error", err)
		slog.Error("forcing shutdown...")
		return
	}

}
