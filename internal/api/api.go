package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Start() {
	fmt.Println("Starting Go Blockchain...")
	http.HandleFunc("/", handleRoot)
	// http.HandleFunc("/download", handleDownload)
	http.ListenAndServe(":3003", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Message: "Golang Blockchain",
		Version: "v1.3.0",
	}
	json.NewEncoder(w).Encode(response)
}
