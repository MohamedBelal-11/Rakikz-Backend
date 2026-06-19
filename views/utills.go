package views

import (
	"encoding/json"
	"io"
	"net/http"
	"rakkiz-backend/errors"
)

func isValidData(data io.Reader, v any) bool {
	return json.NewDecoder(data).Decode(v) == nil
}

func response(
	w http.ResponseWriter,
	data map[string]any,
	status int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(
	w http.ResponseWriter,
	err *errors.AppError,
	status int,
) {
	if err == nil {
		err = &errors.AppError{
			Code: 12001,
			Message: "Unknown Erorr",
		}
	}

	response(
		w,
		map[string]any{
			"success": false,
			"error_code": err.Code,
			"error_message": err.Message,
		},
		status,
	)
}