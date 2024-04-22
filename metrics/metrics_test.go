package metrics

import (
	"testing"
	"time"
)

// helper function to create a successful RequestResult
func successfulRequest(responseTime time.Duration) RequestResult {
	return RequestResult{
		Response:     "OK",
		StatusCode:   200,
		ResponseTime: responseTime,
		Error:        nil,
	}
}

// helper function to create a failed RequestResult
func failedRequest() RequestResult {
	return RequestResult{
		Response:     "Internal Server Error",
		StatusCode:   500,
		ResponseTime: 100 * time.Millisecond, // Failed request still took time
		Error:        nil,
	}
}

func TestCalculateMetrics(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		requestResults []RequestResult
		want          AggregateMetrics
	}{
		{
			name: "All successful requests",
			requestResults: []RequestResult{
				successfulRequest(150 * time.Millisecond),
				successfulRequest(100 * time.Millisecond),
				successfulRequest(200 * time.Millisecond),
			},
			want: AggregateMetrics{
				TotalRequests:     3,
				FailedRequests:    0,
				SuccessRequests:   3,
				SuccessRate:       100.0,
				AverageResponse:   150 * time.Millisecond,
				MinResponse:       100 * time.Millisecond,
				MaxResponse:       200 * time.Millisecond,
				TotalResponseTime: 450 * time.Millisecond,
			},
		},
		{
			name: "Some failed requests",
			requestResults: []RequestResult{
				successfulRequest(120 * time.Millisecond),
				failedRequest(),
				successfulRequest(180 * time.Millisecond),
				failedRequest(),
			},
			want: AggregateMetrics{
				TotalRequests:     4,
				FailedRequests:    2,
				SuccessRequests:   2,
				SuccessRate:       50.0,
				AverageResponse:   150 * time.Millisecond,
				MinResponse:       120 * time.Millisecond,
				MaxResponse:       180 * time.Millisecond,
				TotalResponseTime: 300 * time.Millisecond,
			},
		},
		{
			name: "No requests",
			requestResults: []RequestResult{},
			want: AggregateMetrics{
				TotalRequests:     0,
				FailedRequests:    0,
				SuccessRequests:   0,
				SuccessRate:       0.0,
				AverageResponse:   0,
				MinResponse:       0,
				MaxResponse:       0,
				TotalResponseTime: 0,
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateMetrics(tt.requestResults)
			if got != tt.want {
				t.Errorf("CalculateMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
