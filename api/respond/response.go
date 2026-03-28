package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal JSON response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
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
