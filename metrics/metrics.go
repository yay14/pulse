package main

import (
	"fmt"
)

// Metric represents a telemetry metric with its associated properties
type Metric struct {
	Name      string            // Name of the metric
	Tags      map[string]string // Tags associated with the metric (key-value pairs)
	Value     float64           // Value of the metric
	Timestamp int64             // Timestamp when the metric was recorded (Unix time)
}

// NewMetric creates a new Metric instance
func NewMetric(name string, tags map[string]string, value float64, timestamp int64) *Metric {
	return &Metric{
		Name:      name,
		Tags:      tags,
		Value:     value,
		Timestamp: timestamp,
	}
}

// FormatForVictoriaMetrics formats the metric for sending to VictoriaMetrics
func (m *Metric) FormatForVictoriaMetrics() string {
	var tagsBuffer string
	for key, val := range m.Tags {
		if tagsBuffer != "" {
			tagsBuffer += "," // Add comma between tags
		}
		tagsBuffer += fmt.Sprintf("%s=%s", key, val)
	}

	return fmt.Sprintf("%s,%s %f %d\n", m.Name, tagsBuffer, m.Value, m.Timestamp)
}
