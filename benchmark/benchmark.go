package benchmark

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/komuvill/api_benchmarker/httpclient"
	"github.com/komuvill/api_benchmarker/metrics"
)

type BenchmarkConfig struct {
	URL         string
	Method      string
	Requests    int
	Concurrency int
	Duration    int
	Body        string
}

func RunBenchmark(config *BenchmarkConfig) []metrics.RequestResult {
	fmt.Printf("Benchmarking %s with %s method, %d requests, %d concurrent requests, for %d seconds\n", config.URL, config.Method, config.Requests, config.Concurrency, config.Duration)

	results := make(chan metrics.RequestResult, config.Requests)

	go startWorkers(config, results)

	allResults := collectResults(results)
	return allResults
}

func startWorkers(config *BenchmarkConfig, results chan<- metrics.RequestResult) {
	var wg sync.WaitGroup
	concurrencySemaphore := make(chan struct{}, config.Concurrency)
	testDuration := time.Duration(config.Duration) * time.Second
	timer := time.NewTimer(testDuration)

	for i := 0; i < config.Requests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Create a new reader for each request inside the goroutine
			var requestBody io.Reader
			var err error
			if config.Body != "" {
				requestBody, _, err = httpclient.GetRequestBody(config.Body)
				if err != nil {
					results <- metrics.RequestResult{
						RequestID:    i,
						Response:     "Failed to construct request body",
						StatusCode:   0,
						ResponseTime: 0,
						Error:        err,
					}
					return
				}
			}

			select {
			case concurrencySemaphore <- struct{}{}:
				// This blocks if concurrency limit is reached
				startTime := time.Now()
				// Perform the HTTP request here and capture the result...
				responseBody, statusCode, err := httpclient.HttpRequest(config.Method, config.URL, requestBody)
				responseTime := time.Since(startTime)

				// Send the result to the results channel
				results <- metrics.RequestResult{
					RequestID:    i,
					Response:     responseBody,
					StatusCode:   statusCode,
					ResponseTime: responseTime,
					Error:        err,
				}

				// Release the concurrency semaphore
				<-concurrencySemaphore
			case <-timer.C:
				// If the timer has expired, stop making new requests
				return
			}
		}(i)
	}

	wg.Wait()
	close(results)
	timer.Stop()
}

func collectResults(results <-chan metrics.RequestResult) []metrics.RequestResult {
	var allResults []metrics.RequestResult

	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}
