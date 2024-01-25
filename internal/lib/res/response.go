package res

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Could not marshal JSON: %v\nError: %s\n", payload, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)

}

type errorResponse struct {
	ErrMessage string `json:"error"`
}

func JSONError(w http.ResponseWriter, status int, errMessage string) {
	if status >= 500 {
		log.Printf("Responding with HTTP %v Status.\nError: %s\n", status, errMessage)
	}

	JSONResponse(w, status, errorResponse{ErrMessage: errMessage})
}
