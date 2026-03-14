package utils

import (
	encoding/json
	net/http
)

// Response is the structure for API responses
type Response struct {
	Success bool        `json:"success"` 
	Message string      `json:"message"`
	Data   interface{} `json:"data,omitempty"`
}

// SendResponse sends a JSON response to the client
func SendResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Success: success,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}