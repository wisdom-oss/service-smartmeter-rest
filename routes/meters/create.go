package meters

import (
	"errors"
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thanhpk/randstr"

	wisdomTypes "github.com/wisdom-oss/common-go/v3/types"

	"microservice/internal/db"
	"microservice/routes"
	"microservice/types"
)

type newMeterParameters struct {
	Name                 string         `binding:"required"          json:"name"`
	Key                  *string        `json:"key"`
	Latitude             float64        `binding:"required"          json:"latitude"`
	Longitude            float64        `binding:"required"          json:"longitude"`
	AdditionalParameters map[string]any `json:"additionalProperties"`
}

var ErrBadSmartmeterData = wisdomTypes.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "BadParameterData",
	Detail: "The provided data could not be used to create a new smartmeter",
}

const defaultKeyLength = 16

func Create(c *gin.Context) {
	var newMeter newMeterParameters
	if err := c.ShouldBindJSON(&newMeter); err != nil {
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

	if newMeter.Key == nil {
		newKey := randstr.Base62(defaultKeyLength)
		newMeter.Key = &newKey
	}

	query, err := db.Queries.Raw("insert-meter")
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	var smartmeter types.Smartmeter
	err = pgxscan.Get(c, db.Pool, &smartmeter, query,
		newMeter.Name, newMeter.Key,
		newMeter.Longitude, newMeter.Latitude,
		newMeter.AdditionalParameters,
	)
	if err != nil {
		c.Abort()
		e := routes.ErrDatabaseError
		e.Errors = append(e.Errors, err)
		e.Emit(c)
		return
	}

	c.JSON(http.StatusCreated, smartmeter)
}
