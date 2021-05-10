package util

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON ....
func RespondWithJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		panic("when encoding json: " + err.Error()) // negroni catch it
	}

}

// RespondWithError ...
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, ResponseBase(false, message))
}

// ResponseBase ...
func ResponseBase(success bool, message string) map[string]interface{} {
	return map[string]interface{}{"success": success, "message": message}
}

// ResponseBaseOK ...
func ResponseBaseOK() map[string]interface{} {
	return ResponseBase(true, "OK")
}
