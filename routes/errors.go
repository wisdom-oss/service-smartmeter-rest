package routes

import (
	"net/http"

	"github.com/wisdom-oss/common-go/v3/types"
)

// ErrDatabaseError is a common error which is used to wrap errors that occurred
// during the database queries and preperation thereof.
var ErrDatabaseError = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.6.1",
	Status: http.StatusInternalServerError,
	Title:  "Database Error",
	Detail: "Unable to load requested data from the database",
}

// ErrUnkownSmartmeter is a common error which is sent back if the provided
// smartmeter ID is not available in the database.
var ErrUnknownSmartmeter = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: http.StatusNotFound,
	Title:  "Unknown Smartmeter",
	Detail: "The requested smartmeter does not exist",
}
