package meters

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-gonic/gin"

	testdataloader "github.com/peteole/testdata-loader"

	"microservice/internal"
	"microservice/internal/db"
)

const routePrefix = "/meters"

var router *gin.Engine
var testMeters []newMeterParameters

// TestMain sets up the general testing for the routes in this package.
// This includes creating some smart meters using Queries before testing the
// insertion routes and other routes.
func TestMain(m *testing.M) {
	_ = internal.ParseConfiguration() // error ignored as function always returns nil

	err := db.Connect()
	if err != nil {
		slog.Error("unable to connect to the database", "error", err)
		os.Exit(1)
	}

	slog.Info("checking and running applyable database migrations")
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

	err = json.Unmarshal(testdataloader.GetTestFile("testdata/meters.json"), &testMeters)
	if err != nil {
		slog.Error("unable to load meter testdata")
		os.Exit(1)
	}

	router = gin.New()
	meterAPI := router.Group("/meters")
	{
		meterAPI.GET("/", All)
		meterAPI.PUT("/", Create)
		meterAPI.GET("/:meterID", GetSingle)
	}

	os.Exit(m.Run())
}

func TestMeterAPI(t *testing.T) {
	t.Run("Create", testCreation)
	t.Run("List", listMeters)
	t.Run("Single_Meter", testSingleMeter)
	t.Cleanup(func() {
		for _, testMeter := range testMeters {
			query, err := db.Queries.Raw("cleanup-smartmeters")
			if err != nil {
				t.Fatal("unable to load cleanup query. database left in inconsistent state")
				return
			}

			_, err = db.Pool.Exec(context.Background(), query, testMeter.Key)
			if err != nil {
				t.Fatal("unable to execute cleanup query. database left in inconsistent state")
				return
			}
		}

	})
}
