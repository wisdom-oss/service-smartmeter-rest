package data

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jszwec/csvutil"

	wisdomTypes "github.com/wisdom-oss/common-go/v3/types"

	"microservice/internal/db"
	"microservice/routes"
	"microservice/types"
)

type timeseriesBounds struct {
	Lower time.Time `binding:"omitempty"               form:"from"`
	Upper time.Time `binding:"omitempty,gtfield=Lower" form:"until"`
}

var ErrBadTimebounds = wisdomTypes.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Bad Time Bounds",
	Detail: "The provided timestamps are not usable as boundaries for the data",
}

const (
	MIME_CSV    = "text/csv"
	MIME_NDJSON = "application/x-ndjson"
)

func Timeseries(c *gin.Context) {
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

	var dataBounds timeseriesBounds
	err = c.ShouldBindQuery(&dataBounds)
	if err != nil {
		c.Abort()
		e := ErrBadTimebounds
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	params := []any{meterID}
	var queryName string
	switch {
	case dataBounds.Lower.IsZero() && dataBounds.Upper.IsZero():
		queryName = "get-data"
	case !dataBounds.Lower.IsZero() && dataBounds.Upper.IsZero():
		params = append(params, dataBounds.Lower)
		queryName = "get-data-from"
	case dataBounds.Lower.IsZero() && !dataBounds.Upper.IsZero():
		params = append(params, dataBounds.Upper)
		queryName = "get-data-until"
	case !dataBounds.Lower.IsZero() && !dataBounds.Upper.IsZero():
		params = append(params, dataBounds.Lower, dataBounds.Upper)
		queryName = "get-data-in-range"
	}

	query, err = db.Queries.Raw(queryName)
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	var datapoints []types.DataPoint
	err = pgxscan.Select(c, db.Pool, &datapoints, query, params...)
	if err != nil {
		if pgxscan.NotFound(err) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	format := c.NegotiateFormat(MIME_CSV, MIME_NDJSON)
	switch format {
	case MIME_CSV:
		c.Header("Content-Type", MIME_CSV)
		c.Status(http.StatusOK)

		csvWriter := csv.NewWriter(c.Writer)
		err = csvutil.NewEncoder(csvWriter).Encode(datapoints)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}
	case MIME_NDJSON:
		c.Header("Content-Type", MIME_NDJSON)
		c.Status(http.StatusOK)

		enc := json.NewEncoder(c.Writer)
		for _, datapoint := range datapoints {
			err = enc.Encode(datapoint)
			if err != nil {
				break
			}
		}
	default:
		c.JSON(http.StatusOK, datapoints)
		return
	}

	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
}
