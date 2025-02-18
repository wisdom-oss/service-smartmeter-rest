package meters

import (
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"

	"microservice/internal/db"
	"microservice/routes"
	"microservice/types"
)

func All(c *gin.Context) {
	query, err := db.Queries.Raw("all-meters")
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	var meters []types.Smartmeter
	err = pgxscan.Select(c, db.Pool, &meters, query)
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	c.JSON(http.StatusOK, meters)
}
