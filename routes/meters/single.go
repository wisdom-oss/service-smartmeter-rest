package meters

import (
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"

	"microservice/internal/db"
	"microservice/routes"
	"microservice/types"
)

func GetSingle(c *gin.Context) {
	meterID := c.Param("meterID")

	query, err := db.Queries.Raw("get-meter")
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	var smartmeter types.Smartmeter
	err = pgxscan.Get(c, db.Pool, &smartmeter, query, meterID)
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

	c.JSON(http.StatusOK, smartmeter)
}
