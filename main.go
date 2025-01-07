package main

import (
	"log/slog"
	"os"

	"microservice/internal"
	"microservice/internal/db"
)

func main() {
	_ = internal.ParseConfiguration() // error ignored as function always returns nil

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

}
