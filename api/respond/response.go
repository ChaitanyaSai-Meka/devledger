package respond

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteCreated(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"data":data,
	})
}

func WriteOK(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": data,
	})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]any{
		"success": false,
		"error": message,
	})
}