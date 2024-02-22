package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	wisdomType "github.com/wisdom-oss/commonTypes"

	"github.com/joho/godotenv"

	"github.com/qustavo/dotsql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/wisdom-oss/service-smartmeter-rest/globals"
)

// DefaultAuth contains the default authentication configuration if no file
// is present. it only allows named users access to this service who use the
// same group as the service name
var DefaultAuth = wisdomType.AuthorizationConfiguration{
	Enabled:                   true,
	RequireUserIdentification: true,
	RequiredUserGroup:         globals.ServiceName,
}

var initLogger = log.With().Bool("startup", true).Logger()

// init is executed at every startup of the microservice and is always executed
// before main
func init() {
	// load the variables found in the .env file into the process environment
	err := godotenv.Load()
	if err != nil {
		initLogger.Debug().Msg("no .env files found")
	}

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "-health" {
		conn, err := net.Dial("tcp", tcpSocket)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to connect to tcp socket")
		}
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		conn.Write([]byte("ping\n"))

		inputBuffer := make([]byte, BUF_SIZE)
		n, err := conn.Read(inputBuffer)
		if err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
			log.Fatal().Err(err).Msg("heartbeat server responded too slow")
		}

		returnedMessage := strings.TrimSpace(string(inputBuffer[:n-1]))

		if returnedMessage != "success" {
			fmt.Fprint(os.Stderr, returnedMessage)
			os.Exit(1)
		}
		os.Exit(0)
	}

	configureLogger()
	loadServiceConfiguration()
	setupAuthorization()
	connectDatabase()
	loadPreparedQueries()

	initLogger.Info().Msg("initialization process finished")

}

// configureLogger handles the configuration of the logger used in the
// microservice. it reads the logging level from the `LOG_LEVEL` environment
// variable and sets it according to the parsed logging level. if an invalid
// value is supplied or no level is supplied, the service defaults to the
// `INFO` level
func configureLogger() {
	// set the time format to unix timestamps to allow easier machine handling
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// allow the logger to create an error stack for the logs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// now use the environment variable `LOG_LEVEL` to determine the logging
	// level for the microservice.
	rawLoggingLevel, isSet := os.LookupEnv("LOG_LEVEL")

	// if the value is not set, use the info level as default.
	var loggingLevel zerolog.Level
	if !isSet {
		loggingLevel = zerolog.InfoLevel
	} else {
		// now try to parse the value of the raw logging level to a logging
		// level for the zerolog package
		var err error
		loggingLevel, err = zerolog.ParseLevel(rawLoggingLevel)
		if err != nil {
			// since an error occurred while parsing the logging level, use info
			loggingLevel = zerolog.InfoLevel
			initLogger.Warn().Msg("unable to parse value from environment. using info")
		}
	}
	// since now a logging level is set, configure the logger
	zerolog.SetGlobalLevel(loggingLevel)
}

// loadServiceConfiguration handles loading the `environment.json` file which
// describes which environment variables are needed for the service to function
// and what variables are optional and their default values
func loadServiceConfiguration() {
	initLogger.Info().Msg("loading service configuration from environment")
	// now check if the default location for the environment configuration
	// was changed via the `ENV_CONFIG_LOCATION` variable
	location, locationChanged := os.LookupEnv("ENV_CONFIG_LOCATION")
	if !locationChanged {
		// since the location has not changed, set the default value
		location = "./environment.json"
		initLogger.Debug().Msg("location for environment config not changed")
	}
	initLogger.Debug().Str("path", location).Msg("loading environment requirements file")
	var c wisdomType.EnvironmentConfiguration
	err := c.PopulateFromFilePath(location)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to load environment requirements file")
	}
	initLogger.Info().Msg("validating environment variables")
	globals.Environment, err = c.ParseEnvironment()
	if err != nil {
		initLogger.Fatal().Err(err).Msg("environment validation failed")
	}
	initLogger.Info().Msg("loaded service configuration from environment")
}

// setupAuthorization configures the authorization requirements for the
// microservice
//
// It first checks if the authorization configuration file location is set in
// the environment variables.
// If it is not set, it logs a warning message and uses the default
// authorization configuration.
// If the path is set but empty, it logs a warning message and uses the default
// configuration.
// If the path is not empty, it tries to populate the authorization
// configuration from the file.
// If it encounters an error while populating the configuration, it logs a
// warning message and uses the default configuration.
// Finally, it logs a message indicating that the authorization requirements
// have been configured.
func setupAuthorization() {
	initLogger.Info().Msg("configuring authorization requirements")
	filePath, isSet := globals.Environment["AUTH_CONFIG_FILE_LOCATION"]
	if !isSet {
		initLogger.Warn().Interface("defaultConfig", DefaultAuth).Msg("no authorization configuration file found. using default")
		globals.AuthorizationConfiguration = DefaultAuth
		return
	}

	if strings.TrimSpace(filePath) == "" {
		initLogger.Warn().Interface("defaultConfig", DefaultAuth).Msg("empty path set for authorization configuration file. using default")
		globals.AuthorizationConfiguration = DefaultAuth
		return
	}

	err := globals.AuthorizationConfiguration.PopulateFromFilePath(filePath)
	if err != nil {
		initLogger.Warn().Err(err).Msg("unable to populate configuration from detected file. using default")
		globals.AuthorizationConfiguration = DefaultAuth
		return
	}

	initLogger.Info().Msg("configured authorization requirements")
}

// connectDatabase uses the previously read environment variables to connect the
// microservice to the PostgreSQL database used as the backend for all WISdoM
// services
func connectDatabase() {
	initLogger.Info().Msg("connecting to the database")

	address := fmt.Sprintf("postgres://%s:%s@%s:%s/wisdom",
		globals.Environment["PG_USER"], globals.Environment["PG_PASS"],
		globals.Environment["PG_HOST"], globals.Environment["PG_PORT"])

	var err error
	config, err := pgxpool.ParseConfig(address)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to create base configuration for connection pool")
	}
	globals.Db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to create database connection pool")
	}
	err = globals.Db.Ping(context.Background())
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to verify the connection to the database")
	}
	initLogger.Info().Msg("database connection established")
}

// loadPreparedQueries loads the prepared SQL queries from a file specified by
// the QUERY_FILE_LOCATION environment variable.
// It initializes the SqlQueries variable with the loaded queries.
// If there is an error loading the queries, it logs a fatal error and the
// program terminates.
// This function is typically called during the startup of the microservice.
func loadPreparedQueries() {
	initLogger.Info().Msg("loading prepared sql queries")
	var err error
	globals.SqlQueries, err = dotsql.LoadFromFile(globals.Environment["QUERY_FILE_LOCATION"])
	if err != nil {
		initLogger.Fatal().Err(err).Msg("failed to load prepared queries")
	}
}
