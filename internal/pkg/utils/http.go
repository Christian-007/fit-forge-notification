package utils

import (
	"encoding/json"
	"net/http"
)

type CollectionRes[Entity any] struct {
	Results []Entity `json:"results"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func SendResponse(w http.ResponseWriter, statusCode int, response any) {
	jsonRes, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}
