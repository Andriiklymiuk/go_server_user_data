package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, ErrorResponse{Message: message})
}

func ParseJsonValidationError(err error) string {
	errs := err.(validator.ValidationErrors)
	var errMessage string

	for _, e := range errs {
		fieldName := e.Field()
		lowercasedFieldName := strings.ToLower(fieldName[:1]) + fieldName[1:]
		// can translate each error one at a time.
		switch e.Tag() {
		case "required":
			errMessage += lowercasedFieldName + " is a required field. "
		default:
			errMessage += "Validation failed on " + e.Field() + ". "
		}
	}

	return errMessage
}
