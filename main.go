package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"net/http"

	metrics "github.com/yay14/pulse/metrics"
	"google.golang.org/grpc"
)

type server struct {
	metrics.UnimplementedMetricsServiceServer
}

// WriteMetrics writes metrics to VictoriaMetrics
func (s *server) WriteMetrics(ctx context.Context, req *metrics.WriteRequest) (*metrics.WriteResponse, error) {
	vmURL := os.Getenv("VICTORIA_METRICS_URL")

	// Prepare data for VictoriaMetrics
	var jsonData []map[string]interface{}
	for _, ts := range req.Timeseries {
		for _, sample := range ts.Samples {
			dataPoint := map[string]interface{}{
				"__name__":  "http_requests_total",
				"method":    ts.Labels["method"],
				"handler":   ts.Labels["handler"],
				"status":    ts.Labels["status"],
				"value":     sample.Value,
				"timestamp": sample.Timestamp,
			}
			jsonData = append(jsonData, dataPoint)
		}
	}

	jsonBody, err := json.Marshal(jsonData)
	if err != nil {
		return &metrics.WriteResponse{Status: "Error marshalling data"}, err
	}

	// Send data to VictoriaMetrics
	resp, err := http.Post(vmURL+"/api/v1/import", "application/json", bytes.NewReader(jsonBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		return &metrics.WriteResponse{Status: "Failed to send data to VictoriaMetrics"}, err
	}

	return &metrics.WriteResponse{Status: "Data sent to VictoriaMetrics successfully"}, nil
}

// QueryMetrics queries metrics from VictoriaMetrics
func (s *server) QueryMetrics(ctx context.Context, req *metrics.ReadRequest) (*metrics.ReadResponse, error) {
	vmURL := os.Getenv("VICTORIA_METRICS_URL")

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
	lis, err := net.Listen("tcp", ":50051")
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
