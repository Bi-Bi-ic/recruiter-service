package utils

import (
	"encoding/json"
	"net/http"

	"github.com/rgrs-x/service/api/models"
)

// Message ...
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

// Respond ...
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//ResultRepository ...
type ResultRepository struct {
	Result  interface{}
	Message string
	Error   error
	Code    models.Code
}

// Response ...
type Response struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
	Code    *models.Code `json:"code,omitempty"`
}

// Response ...
type BTResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
}
