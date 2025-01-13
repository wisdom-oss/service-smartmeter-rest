package data

import (
	"encoding/csv"
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jszwec/csvutil"

	wisdomTypes "github.com/wisdom-oss/common-go/v3/types"

	"microservice/internal/db"
	"microservice/routes"
	"microservice/types"
)

type InsertResponse struct {
	InsertedDatapoints int64 `json:"insertedDatapoints"`
	InvalidDatapoints  []int `json:"invalidDatapoints,omitempty"`
}

var ErrInvalidContentType = wisdomTypes.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.6.1",
	Status: http.StatusBadRequest,
	Title:  "Unsupported Content Type",
	Detail: "This endpoint only supports the 'text/csv' content type",
}

func Import(c *gin.Context) {
	if c.ContentType() != "text/csv" {
		c.Abort()
		ErrInvalidContentType.Emit(c)
		return
	}

	meterID := c.Param("meterID")

	query, err := db.Queries.Raw("meter-exists")
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	var exists bool
	err = pgxscan.Get(c, db.Pool, &exists, query, meterID)
	if err != nil {
		if pgxscan.NotFound(err) {
			c.Abort()
			routes.ErrUnknownSmartmeter.Emit(c)
			return
		}
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	var datapoints []types.DataPoint
	csvreader := csv.NewReader(c.Request.Body)
	decoder, err := csvutil.NewDecoder(csvreader)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	err = decoder.Decode(&datapoints)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var validDatapoints []types.DataPoint //nolint: prealloc
	var invalidDatapointLines []int
	for idx, datapoint := range datapoints {
		if !datapoint.Timestamp.Valid {
			invalidDatapointLines = append(invalidDatapointLines, idx)
			continue
		}
		validDatapoints = append(validDatapoints, datapoint)
	}

	insertedRows, err := db.Pool.CopyFrom(
		c,
		pgx.Identifier{"timeseries", "smartmeter_data"},
		[]string{"time", "smart_meter", "flow_rate"},
		pgx.CopyFromSlice(len(validDatapoints), func(i int) ([]any, error) {
			return []any{validDatapoints[i].Timestamp, meterID, validDatapoints[i].FlowRate}, nil
		}),
	)
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	c.JSON(http.StatusCreated, InsertResponse{
		InsertedDatapoints: insertedRows,
		InvalidDatapoints:  invalidDatapointLines,
	})
}
