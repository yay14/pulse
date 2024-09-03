package main

import (
	"log"
	"net"
	"time"

	"github.com/gocql/gocql"
	"github.com/yay14/pulse/ingestion"
	"github.com/yay14/pulse/internal/cassandra"
	"github.com/yay14/pulse/internal/kafka"
	ingestionSvc "github.com/yay14/pulse/internal/service/ingestion"
	metricsSvc "github.com/yay14/pulse/internal/service/metrics"
	"github.com/yay14/pulse/metrics"
	"google.golang.org/grpc"
)

func main() {
	// Connect to Cassandra
	cluster := gocql.NewCluster("cassandra")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	repo, err := cassandra.NewRepository(cluster)
	if err != nil {
		log.Fatalf("failed to connect to Cassandra: %v", err)
	}
	defer repo.Close()

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":9400")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	ingestionService := ingestionSvc.NewIngestionService(repo)
	metricsService := metricsSvc.NewMetricsService()
	ingestion.RegisterIngestionServiceServer(grpcServer, ingestionService)
	metrics.RegisterMetricsServiceServer(grpcServer, metricsService)

	// Start Kafka consumer
	kafkaConfig := kafka.KafkaConfig{
		Brokers: []string{"kafka:9092"},
		GroupID: "pulse",
		Topic:   "metrics-topic",
	}
	go ingestionService.StartKafkaConsumer(kafkaConfig)

	log.Println("Starting gRPC server on :9400...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
