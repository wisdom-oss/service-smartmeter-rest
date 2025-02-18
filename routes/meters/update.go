package meters

import (
	"errors"
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"microservice/internal/db"
	"microservice/routes"
	"microservice/types"
)

type updateParameters struct {
	Name                 *string        `json:"name"`
	Latitude             *float64       `binding:"required_with=Longitude" json:"latitude"`
	Longitude            *float64       `binding:"required_with=Latitude"  json:"longitude"`
	AdditionalParameters map[string]any `json:"additionalProperties"`
}

func Update(c *gin.Context) {
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

	var updateValues updateParameters
	err = c.ShouldBindBodyWithJSON(&updateValues)
	if err != nil {
		c.Abort()
		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			e := ErrBadSmartmeterData
			e.Errors = routes.GenerateErrors(validationError)
			e.Emit(c)
			return
		}
		_ = c.Error(err)
		return
	}

	if updateValues.Name != nil {
		query, err := db.Queries.Raw("update-meter-name")
		if err != nil {
			c.Abort()
			e := routes.ErrDatabaseError
			e.Errors = append(e.Errors, err)
			e.Emit(c)
			return
		}

		_, err = db.Pool.Exec(c, query, updateValues.Name, meterID)
		if err != nil {
			c.Abort()
			e := routes.ErrDatabaseError
			e.Errors = append(e.Errors, err)
			e.Emit(c)
			return
		}
	}

	if updateValues.Latitude != nil && updateValues.Longitude != nil {
		query, err := db.Queries.Raw("update-meter-location")
		if err != nil {
			c.Abort()
			e := routes.ErrDatabaseError
			e.Errors = append(e.Errors, err)
			e.Emit(c)
			return
		}

		_, err = db.Pool.Exec(c, query, updateValues.Longitude, updateValues.Latitude, meterID)
		if err != nil {
			c.Abort()
			e := routes.ErrDatabaseError
			e.Errors = append(e.Errors, err)
			e.Emit(c)
			return
		}
	}

	if updateValues.AdditionalParameters != nil {
		query, err := db.Queries.Raw("update-meter-properties")
		if err != nil {
			c.Abort()
			e := routes.ErrDatabaseError
			e.Errors = append(e.Errors, err)
			e.Emit(c)
			return
		}

		_, err = db.Pool.Exec(c, query, updateValues.AdditionalParameters, meterID)
		if err != nil {
			c.Abort()
			e := routes.ErrDatabaseError
			e.Errors = append(e.Errors, err)
			e.Emit(c)
			return
		}
	}

	query, err = db.Queries.Raw("get-meter")
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
