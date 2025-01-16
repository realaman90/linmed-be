package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// paginated respose struct
type paginatedResponse struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}

// Helper function to write JSON response
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError) // Handle encoding error
	}
}

// Helper function to write error response in a JSON format, with error message

func errorResposne(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "error",
		"error":  message,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError) // Handle encoding error
	}
}

// validate page and limit query parameters
func (s *Server) validatePageLimit(page, limit string) (int, int) {
	// default values
	pageInt := 1
	limitInt := 10

	if page != "" {
		pageInt = convertToInt(page)
	}

	if limit != "" {
		limitInt = convertToInt(limit)
	}

	return pageInt, limitInt
}

// convert string to int
func convertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func (s *Server) stringToUint(str string) (uint, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}
