package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/yay14/pulse/metrics"
	"google.golang.org/grpc"
)

type server struct {
	metrics.UnimplementedMetricsServiceServer
}

// WriteMetrics writes metrics to VictoriaMetrics
func (s *server) WriteMetrics(ctx context.Context, req *metrics.WriteRequest) (*metrics.WriteResponse, error) {
	vmURL := os.Getenv("VICTORIA_METRICS_URL")
	log.Printf("VictoriaMetrics URL: %s", vmURL)

	// Prepare data for VictoriaMetrics in the correct format
	var jsonData []string
	for _, ts := range req.Timeseries {
		// Prepare the metric object
		metric := map[string]interface{}{
			"__name__": ts.Labels["__name__"], // Ensure the metric name is included
		}

		// Add additional labels to the metric object
		for key, value := range ts.Labels {
			if key != "__name__" { // Skip the metric name
				metric[key] = value
			}
		}

		// Prepare values and timestamps slices
		values := make([]float64, len(ts.Samples))
		timestamps := make([]int64, len(ts.Samples))

		for i, sample := range ts.Samples {
			values[i] = sample.Value
			timestamps[i] = sample.Timestamp // Ensure timestamps are in milliseconds
		}

		// Create the JSON object for the current time series
		jsonLine := map[string]interface{}{
			"metric":     metric,
			"values":     values,
			"timestamps": timestamps,
		}

		// Marshal the JSON object to a string
		jsonDataLine, err := json.Marshal(jsonLine)
		if err != nil {
			return &metrics.WriteResponse{Status: "Error marshalling data"}, err
		}

		// Append the JSON line to the data slice
		jsonData = append(jsonData, string(jsonDataLine))
	}

	// Send data to VictoriaMetrics
	for _, line := range jsonData {
		_, err := http.Post(vmURL+"/api/v1/import", "application/json", bytes.NewReader([]byte(line)))
		if err != nil {
			return &metrics.WriteResponse{Status: "Failed to send data to VictoriaMetrics"}, err
		}
	}
	return &metrics.WriteResponse{Status: "Data sent to VictoriaMetrics successfully"}, nil
}

// QueryMetrics queries metrics from VictoriaMetrics
func (s *server) QueryMetrics(ctx context.Context, req *metrics.ReadRequest) (*metrics.ReadResponse, error) {
	vmURL := os.Getenv("VICTORIA_METRICS_URL")
	log.Printf("VictoriaMetrics URL: %s", vmURL)

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=%s", vmURL, req.Query))
	if err != nil {
		return &metrics.ReadResponse{}, err
	}
	defer resp.Body.Close()

	var readResponse metrics.ReadResponse
	if err := json.NewDecoder(resp.Body).Decode(&readResponse); err != nil {
		return &metrics.ReadResponse{}, err
	}

	return &readResponse, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9400")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	metrics.RegisterMetricsServiceServer(grpcServer, &server{})

	log.Println("Starting gRPC server on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
