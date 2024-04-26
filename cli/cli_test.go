package cli

import (
	"testing"

	"github.com/komuvill/api_benchmarker/benchmark"
)

// TestValidateFlags tests the validation logic for command-line flags.
func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name    string
		config  benchmark.BenchmarkConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			config: benchmark.BenchmarkConfig{
				URL:         "http://example.com",
				Method:      "GET",
				Requests:    100,
				Concurrency: 10,
				Duration:    1,
				Body:        "",
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			config: benchmark.BenchmarkConfig{
				Method: "GET",
			},
			wantErr: true,
			errMsg:  "URL is required",
		},
		{
			name: "invalid HTTP method",
			config: benchmark.BenchmarkConfig{
				URL:    "http://example.com",
				Method: "INVALID",
			},
			wantErr: true,
			errMsg:  "'INVALID' is not a valid HTTP method. Supported methods are: GET, POST, PUT, DELETE",
		},
		{
			name: "missing body for POST",
			config: benchmark.BenchmarkConfig{
				URL:    "http://example.com",
				Method: "POST",
			},
			wantErr: true,
			errMsg:  "a request body is required for the POST method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFlags(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateFlags() gotErr = %v, wantErr %v", err, tt.errMsg)
			}
		})
	}
}
