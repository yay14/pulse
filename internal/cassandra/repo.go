package cassandra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocql/gocql"
	"github.com/yay14/pulse/ingestion"
)

type Repository struct {
	session *gocql.Session
}

// NewRepository creates a new Cassandra repository
func NewRepository(cluster *gocql.ClusterConfig) (*Repository, error) {
	// Create a new session
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cassandra session: %w", err)
	}

	// Create keyspace
	err = session.Query(`
		CREATE KEYSPACE IF NOT EXISTS metrics_keyspace 
		WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};
	`).Exec()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create keyspace: %w", err)
	}

	// Create table metrics
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS metrics_keyspace.metrics (
			id UUID,
			source_id UUID,
			source_type TEXT,
			metric_name TEXT,
			metric_value DOUBLE,
			labels TEXT,
			timestamp TIMESTAMP,
			PRIMARY KEY (id)
		);
	`).Exec()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create metrics table: %w", err)
	}

	// Create table metrics_validation
	err = session.Query(`
	CREATE TABLE IF NOT EXISTS metrics_keyspace.metric_validation (
		id UUID,
		metric_name TEXT PRIMARY KEY,
		source_id UUID,
		min_value DOUBLE,
		max_value DOUBLE,
		PRIMARY KEY (id)
	);
	`).Exec()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create metrics validation table: %w", err)
	}

	return &Repository{session: session}, nil
}

// WriteMetric writes a metric to Cassandra
func (r *Repository) WriteMetric(ctx context.Context, metric *ingestion.MetricData, req *ingestion.IngestDataRequest) error {
	// Generate a new UUID for the primary key 'id'
	id := gocql.TimeUUID()

	// Prepare the CQL query
	query := `INSERT INTO metrics_keyspace.metrics (
        id, 
        source_id, 
        source_type,
        metric_name, 
        metric_value,
        labels,
        timestamp
    ) VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Convert the labels to JSON
	labelsJSON, err := json.Marshal(metric.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	// Execute the CQL query
	if err := r.session.Query(query, id, req.SourceId, req.SourceType, metric.Name, metric.Value, labelsJSON, metric.Timestamp).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	log.Println("Successfully wrote metric to Cassandra")
	return nil
}

// AddMetricValidation adds a validation rule to the metric_validation table
func (r *Repository) AddMetricValidation(ctx context.Context, validation *ingestion.NewValidationRequest) error {
	id := gocql.TimeUUID()

	query := `INSERT INTO metrics_keyspace.metric_validation (
		id,
		metric_name, 
		source_id, 
		min_value, 
		max_value
		) VALUES (?, ?, ?, ?, ?)`

	if err := r.session.Query(query, id, validation.MetricName, validation.SourceId, validation.MinValue, validation.MaxValue).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to add validation: %w", err)
	}

	log.Println("Successfully added validation to Cassandra")
	return nil
}

// ValidateMetric checks if the given metric value is within the predefined range
func (r *Repository) ValidateMetric(ctx context.Context, metricName, sourceId string, metricValue float64) (bool, string, error) {
	var minVal, maxVal float64

	query := `SELECT min_value, max_value FROM metrics_keyspace.metric_validation WHERE metric_name = ? AND source_id = ? LIMIT 1`
	if err := r.session.Query(query, metricName, sourceId).Scan(&minVal, &maxVal); err != nil {
		if err == gocql.ErrNotFound {
			return false, "Validation rule not found for metric", nil
		}
		return false, "", fmt.Errorf("failed to query validation rules: %w", err)
	}

	if metricValue < minVal || metricValue > maxVal {
		return false, fmt.Sprintf("Metric value %f is out of the allowed range [%f, %f]", metricValue, minVal, maxVal), nil
	}

	return true, "Metric is valid", nil
}

// Close closes the Cassandra session
func (r *Repository) Close() {
	r.session.Close()
}
