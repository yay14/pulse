package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/query", queryHandler)

	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Example of reading a metric from the request body
	var metric map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to VictoriaMetrics
	vmURL := os.Getenv("VICTORIA_METRICS_URL")
	resp, err := http.Post(vmURL+"/api/v1/import", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Failed to send data to VictoriaMetrics", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Respond to the client
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte("Data sent to VictoriaMetrics"))
	// This is where you'll handle incoming POST requests to send data to VictoriaMetrics
	w.Write([]byte("POST endpoint - data will be sent to VictoriaMetrics here"))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	// This is where you'll handle queries to retrieve data from VictoriaMetrics
	w.Write([]byte("QUERY endpoint - data will be retrieved from VictoriaMetrics here"))
}
