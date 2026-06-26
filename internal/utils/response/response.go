package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// constructors

func NewSuccess(msg string, data any) Response {
	return Response{
		Success: true,
		Message: msg,
		Data:    data,
	}
}

func NewError(err error) Response {
	return Response{
		Success: false,
		Message: err.Error(),
	}
}
