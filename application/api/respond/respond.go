package respond

import (
	"encoding/json"
	"net/http"
)

func With(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err, ok := data.(error); ok {
		data = map[string]interface{}{"error": err.Error()}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic("Failed to encode response: %v" + err.Error())
	}
}
