package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/yay14/pulse/metrics"
)

type MetricsService struct {
	metrics.UnimplementedMetricsServiceServer
}

// NewIngestionService creates a new IngestionService
func NewMetricsService() *MetricsService {
	return &MetricsService{}
}


// WriteMetrics writes metrics to VictoriaMetrics
func (s *MetricsService) WriteMetrics(ctx context.Context, req *metrics.WriteRequest) (*metrics.WriteResponse, error) {
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

func (s *MetricsService) InstantQueryMetrics(ctx context.Context, req *metrics.InstantQueryReadRequest) (*metrics.ReadResponse, error) {
	vmURL := os.Getenv("VICTORIA_METRICS_URL")

	// Send the request to VictoriaMetrics
	encodedQuery := url.QueryEscape(req.Query)

	queryEndpoint := fmt.Sprintf("%s/api/v1/query?query=%s", vmURL, encodedQuery)
	log.Printf("Querying VictoriaMetrics at %s with query: %s", vmURL, queryEndpoint)

	resp, err := http.Get(queryEndpoint)
	if err != nil {
		log.Printf("Error sending request to VictoriaMetrics: %v", err)
		return &metrics.ReadResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %s", resp.Status)
		return &metrics.ReadResponse{}, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Decode the response
	var vmResponse struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"`
			} `json:"result"`
		} `json:"data"`
		Stats struct {
			SeriesFetched     string `json: "seriesFetched"`
			ExecutionTimeMsec int64  `json: "executionTimeMsec"`
		} `json:"stats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&vmResponse); err != nil {
		log.Printf("Error decoding response body: %v", err)
		return &metrics.ReadResponse{}, err
	}

	log.Printf("Decoded response: %+v", vmResponse)

	readResponse := &metrics.ReadResponse{}

	// Process the results
	for _, result := range vmResponse.Data.Result {
		timeseries := &metrics.TimeseriesData{
			Labels: result.Metric,
		}

		value := result.Value
		timestamp := int64(value[0].(float64))
		metricValue := value[1].(string)

		sample := &metrics.Sample{
			Value:     mustParseFloat64(metricValue),
			Timestamp: timestamp,
		}
		timeseries.Samples = append(timeseries.Samples, sample)

		readResponse.Timeseries = append(readResponse.Timeseries, timeseries)
	}

	log.Printf("Final ReadResponse: %+v", readResponse)

	return readResponse, nil
}

func (s *MetricsService) RangeQueryMetrics(ctx context.Context, req *metrics.RangeQueryReadRequest) (*metrics.ReadResponse, error) {
	vmURL := os.Getenv("VICTORIA_METRICS_URL")

	// Encode the query
	encodedQuery := url.QueryEscape(req.Query)

	// Check if a time range is provided
	var queryEndpoint string
	if req.Start != "" && req.End != "" && req.Step != "" {
		// Use the query_range API for range queries
		queryEndpoint = fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d",
			vmURL, encodedQuery, req.Start, req.End, req.Step)
	} else {
		// Default to the instant query API
		queryEndpoint = fmt.Sprintf("%s/api/v1/query?query=%s", vmURL, encodedQuery)
	}

	log.Printf("Querying VictoriaMetrics with query: %s", queryEndpoint)

	resp, err := http.Get(queryEndpoint)
	if err != nil {
		log.Printf("Error sending request to VictoriaMetrics: %v", err)
		return &metrics.ReadResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %s", resp.Status)
		return &metrics.ReadResponse{}, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var vmResponse struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Values [][]interface{}   `json:"values"` // Array of [timestamp, value]
			} `json:"result"`
		} `json:"data"`
		Stats struct {
			SeriesFetched     string `json:"seriesFetched"`
			ExecutionTimeMsec int64  `json:"executionTimeMsec"`
		} `json:"stats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&vmResponse); err != nil {
		log.Printf("Error decoding response body: %v", err)
		return &metrics.ReadResponse{}, err
	}

	log.Printf("Decoded response: %+v", vmResponse)

	readResponse := &metrics.ReadResponse{}

	for _, result := range vmResponse.Data.Result {
		timeseries := &metrics.TimeseriesData{
			Labels: result.Metric,
		}

		for _, valuePair := range result.Values {
			timestamp := int64(valuePair[0].(float64))
			metricValue := valuePair[1].(string)

			sample := &metrics.Sample{
				Value:     mustParseFloat64(metricValue),
				Timestamp: timestamp,
			}
			timeseries.Samples = append(timeseries.Samples, sample)
		}

		readResponse.Timeseries = append(readResponse.Timeseries, timeseries)
	}

	log.Printf("Final ReadResponse: %+v", readResponse)

	return readResponse, nil
}

func mustParseFloat64(value string) float64 {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing float64 value: %s, error: %v", value, err)
		panic(err)
	}
	return v
}
