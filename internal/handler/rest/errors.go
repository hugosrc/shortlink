package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hugosrc/shortlink/internal/util"
)

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func handleError(w http.ResponseWriter, err error, message string) {
	w.Header().Set("Content-Type", "application/json")

	response := ErrorResponse{
		Code:  http.StatusInternalServerError,
		Error: message,
	}

	var appError *util.Error
	if errors.As(err, &appError) {
		switch appError.Code() {
		case util.ErrCodeInvalidArgument:
			response.Code = http.StatusBadRequest
		case util.ErrCodeNotFound:
			response.Code = http.StatusNotFound
		case util.ErrCodeUnauthorized:
			response.Code = http.StatusUnauthorized
		case util.ErrCodeUnknown:
			response.Code = http.StatusBadRequest
		}
	}

	w.WriteHeader(response.Code)
	_ = json.NewEncoder(w).Encode(response)
}
