package metrics

import (
	"context"
	"reflect"
	"testing"

	"github.com/yay14/pulse/metrics"
)

func Test_mustParseFloat64(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustParseFloat64(tt.args.value); got != tt.want {
				t.Errorf("mustParseFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsService_RangeQueryMetrics(t *testing.T) {
	type fields struct {
		UnimplementedMetricsServiceServer metrics.UnimplementedMetricsServiceServer
	}
	type args struct {
		ctx context.Context
		req *metrics.RangeQueryReadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *metrics.ReadResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{
				UnimplementedMetricsServiceServer: tt.fields.UnimplementedMetricsServiceServer,
			}
			got, err := s.RangeQueryMetrics(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsService.RangeQueryMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricsService.RangeQueryMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsService_InstantQueryMetrics(t *testing.T) {
	type fields struct {
		UnimplementedMetricsServiceServer metrics.UnimplementedMetricsServiceServer
	}
	type args struct {
		ctx context.Context
		req *metrics.InstantQueryReadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *metrics.ReadResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{
				UnimplementedMetricsServiceServer: tt.fields.UnimplementedMetricsServiceServer,
			}
			got, err := s.InstantQueryMetrics(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsService.InstantQueryMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricsService.InstantQueryMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsService_WriteMetrics(t *testing.T) {
	type fields struct {
		UnimplementedMetricsServiceServer metrics.UnimplementedMetricsServiceServer
	}
	type args struct {
		ctx context.Context
		req *metrics.WriteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *metrics.WriteResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{
				UnimplementedMetricsServiceServer: tt.fields.UnimplementedMetricsServiceServer,
			}
			got, err := s.WriteMetrics(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsService.WriteMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricsService.WriteMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
