package ingestion

import (
	"context"
	"reflect"
	"testing"

	"github.com/yay14/pulse/ingestion"
	"github.com/yay14/pulse/internal/cassandra"
)

func TestIngestionService_IngestData(t *testing.T) {
	type fields struct {
		UnimplementedIngestionServiceServer ingestion.UnimplementedIngestionServiceServer
		repo                                *cassandra.Repository
	}
	type args struct {
		ctx context.Context
		req *ingestion.IngestDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ingestion.IngestDataResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IngestionService{
				UnimplementedIngestionServiceServer: tt.fields.UnimplementedIngestionServiceServer,
				repo:                                tt.fields.repo,
			}
			got, err := s.IngestData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IngestionService.IngestData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngestionService.IngestData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngestionService_AddMetricValidation(t *testing.T) {
	type fields struct {
		UnimplementedIngestionServiceServer ingestion.UnimplementedIngestionServiceServer
		repo                                *cassandra.Repository
	}
	type args struct {
		ctx context.Context
		req *ingestion.NewValidationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ingestion.NewValidationResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IngestionService{
				UnimplementedIngestionServiceServer: tt.fields.UnimplementedIngestionServiceServer,
				repo:                                tt.fields.repo,
			}
			got, err := s.AddMetricValidation(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IngestionService.AddMetricValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngestionService.AddMetricValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngestionService_ValidateData(t *testing.T) {
	type fields struct {
		UnimplementedIngestionServiceServer ingestion.UnimplementedIngestionServiceServer
		repo                                *cassandra.Repository
	}
	type args struct {
		ctx context.Context
		req *ingestion.ValidateDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ingestion.ValidateDataResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IngestionService{
				UnimplementedIngestionServiceServer: tt.fields.UnimplementedIngestionServiceServer,
				repo:                                tt.fields.repo,
			}
			got, err := s.ValidateData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IngestionService.ValidateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngestionService.ValidateData() = %v, want %v", got, tt.want)
			}
		})
	}
}
