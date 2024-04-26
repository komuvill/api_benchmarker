package metrics

import (
	"fmt"
	"math"
	"time"
)

// RequestResult stores results from each individual request
type RequestResult struct {
	RequestID	int
	Response     string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

// AggregateMetrics is used for calculating metrics across the whole test
type AggregateMetrics struct {
	TotalRequests     int
	FailedRequests    int
	SuccessRequests   int
	SuccessRate       float64
	AverageResponse   time.Duration
	MinResponse       time.Duration
	MaxResponse       time.Duration
	TotalResponseTime time.Duration // for calculating average response time
}

func NewAggregateMetrics() *AggregateMetrics {
	return &AggregateMetrics{
		MinResponse: time.Duration(math.MaxInt64), // Initialize with the maximum possible value
		MaxResponse: time.Duration(0),             // Initialize with zero
	}
}

func CalculateMetrics(results []RequestResult) AggregateMetrics {
	metrics := NewAggregateMetrics()

	for _, result := range results {
		metrics.TotalRequests++

		if result.Error != nil || result.StatusCode < 200 || result.StatusCode >= 300 {
			metrics.FailedRequests++
		} else {
			// Only successful requests are considered for these metrics
			metrics.SuccessRequests++
			metrics.TotalResponseTime += result.ResponseTime

			if result.ResponseTime < metrics.MinResponse {
				metrics.MinResponse = result.ResponseTime
			}
			if result.ResponseTime > metrics.MaxResponse {
				metrics.MaxResponse = result.ResponseTime
			}
		}
	}

	// Calculate the average response time for successful requests
	if metrics.SuccessRequests > 0 {
		metrics.AverageResponse = metrics.TotalResponseTime / time.Duration(metrics.SuccessRequests)
	}

	// Calculate the success rate
	if metrics.TotalRequests > 0 {
		metrics.SuccessRate = (float64(metrics.SuccessRequests) / float64(metrics.TotalRequests)) * 100
	}

	// Reset MinResponse if no successful requests were recorded
	if metrics.MinResponse == time.Duration(math.MaxInt64) {
		metrics.MinResponse = 0
	}

	return *metrics
}

func PrintMetrics(metrics AggregateMetrics) {
	fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("Successful Requests: %d\n", metrics.SuccessRequests)
	fmt.Printf("Failed Requests: %d\n", metrics.FailedRequests)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.SuccessRate)
	fmt.Printf("Average Response Time: %s\n", metrics.AverageResponse)
	fmt.Printf("Minimum Response Time: %s\n", metrics.MinResponse)
	fmt.Printf("Maximum Response Time: %s\n", metrics.MaxResponse)
}
