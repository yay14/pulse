package ingestion

import (
	"context"
	"log"

	"github.com/yay14/pulse/ingestion"
	"github.com/yay14/pulse/internal/cassandra"
	"github.com/yay14/pulse/internal/kafka"
)

// IngestionService implements the IngestionServiceServer
type IngestionService struct {
	ingestion.UnimplementedIngestionServiceServer
	repo *cassandra.Repository
}

// NewIngestionService creates a new IngestionService
func NewIngestionService(repo *cassandra.Repository) *IngestionService {
	return &IngestionService{repo: repo}
}

// IngestData ingests metrics data from Kafka and writes to Cassandra
func (s *IngestionService) IngestData(ctx context.Context, req *ingestion.IngestDataRequest) (*ingestion.IngestDataResponse, error) {
	log.Println("Ingesting data for source:", req.SourceId)

	for _, metric := range req.Metrics {
		// Write each metric to Cassandra
		err := s.repo.WriteMetric(ctx, metric, req)
		if err != nil {
			log.Printf("Error writing metric to Cassandra: %v", err)
			return &ingestion.IngestDataResponse{Status: "Failed to write to Cassandra"}, err
		}
	}

	return &ingestion.IngestDataResponse{Status: "Data ingested successfully"}, nil
}

// AddMetricValidation adds a validation rule to the database
func (s *IngestionService) AddMetricValidation(ctx context.Context, req *ingestion.NewValidationRequest) (*ingestion.NewValidationResponse, error) {
	err := s.repo.AddMetricValidation(ctx, req)
	if err != nil {
		return &ingestion.NewValidationResponse{
			Success: false,
		}, err
	}

	return &ingestion.NewValidationResponse{Success: true}, nil
}

// ValidateData validates metrics data based on predefined rules
func (s *IngestionService) ValidateData(ctx context.Context, req *ingestion.ValidateDataRequest) (*ingestion.ValidateDataResponse, error) {
	for _, metric := range req.Metrics {
		isValid, validationMessage, err := s.repo.ValidateMetric(ctx, metric.Name, req.SourceId, metric.Value)
		if err != nil {
			return &ingestion.ValidateDataResponse{
				Success: false,
				Message: "Error during validation: " + err.Error(),
			}, err
		}

		if !isValid {
			return &ingestion.ValidateDataResponse{
				Success: false,
				Message: validationMessage,
			}, nil
		}
	}

	return &ingestion.ValidateDataResponse{
		Success: true,
		Message: "All metrics are valid",
	}, nil
}

// StartKafkaConsumer starts consuming messages from Kafka and processes them
func (s *IngestionService) StartKafkaConsumer(cfg kafka.KafkaConfig) {
	kafka.Consumer(cfg, func(message kafka.Message) {
		// Process Kafka message and ingest metrics
		req := &ingestion.IngestDataRequest{
			SourceId:   message.SourceID,
			SourceType: message.SourceType,
			Metrics:    message.Metrics,
		}
		_, err := s.IngestData(context.Background(), req)
		if err != nil {
			log.Printf("Failed to ingest data from Kafka message: %v", err)
		}
	})
}
