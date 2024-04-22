package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type BenchmarkConfig struct {
	URL         string
	Method      string
	Requests    int
	Concurrency int
	Duration    int
	Body        string
}

type RequestResult struct {
	Response     string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

func runBenchmark(config *BenchmarkConfig) {
	fmt.Printf("Benchmarking %s with %s method, %d requests, %d concurrent requests, for %d seconds\n", config.URL, config.Method, config.Requests, config.Concurrency, config.Duration)

	results := make(chan RequestResult, config.Requests)
	done := make(chan struct{})

	go startWorkers(config, results)
	go collectResults(results, done)

	<-done // Wait for the results processing to complete
}

func startWorkers(config *BenchmarkConfig, results chan<- RequestResult) {
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
				requestBody, _, err = getRequestBody(config.Body)
				if err != nil {
					results <- RequestResult{
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
				responseBody, statusCode, err := httpRequest(config.Method, config.URL, requestBody)
				responseTime := time.Since(startTime)

				// Send the result to the results channel
				results <- RequestResult{
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

func collectResults(results <-chan RequestResult, done chan<- struct{}) {
	var allResults []RequestResult
	for result := range results {
		allResults = append(allResults, result)
	}
	calculateAndPrintMetrics(allResults)
	close(done)
}

func calculateAndPrintMetrics(results []RequestResult) {
	// Iterate through the results and print each one
	for _, result := range results {
		fmt.Printf("Result: %+v\n", result)
	}
}

// httpRequest sends an HTTP request and returns the response body as a string, the status code, and an error if any.
func httpRequest(method, url string, body io.Reader) (string, int, error) {
	// Create a context with a timeout to allow for request cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new HTTP request with the context
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %v", err)
	}

	// Set the Content-Type header if there is a body.
	// For the purpose of the project, assume the API to be tested expects JSON requests.
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Create an HTTP client with timeout settings
	client := &http.Client{
		Timeout: time.Second * 30, // Optionally set a timeout for the client
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("error reading response body: %v", err)
	}

	return string(respBody), resp.StatusCode, nil
}

// getRequestBody handles the retrieval of the request body.
// It checks if the body is provided as a raw string or a file path.
func getRequestBody(bodyFlag string) (io.Reader, func(), error) {
	if strings.HasPrefix(bodyFlag, "@") {
		filePath := strings.TrimPrefix(bodyFlag, "@")
		file, err := os.Open(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("error opening file: %v", err)
		}

		// Return the file and a cleanup function that closes the file
		return file, func() {
			file.Close()
		}, nil
	}

	// For raw string bodies, no cleanup is needed, so we return a no-op function
	return strings.NewReader(bodyFlag), func() {}, nil
}

func validateFlags(config *BenchmarkConfig) error {
	// Validate URL
	if config.URL == "" {
		return fmt.Errorf("URL is required")
	}

	// Validate Method
	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
	}

	if _, valid := validMethods[config.Method]; !valid {
		return fmt.Errorf("'%s' is not a valid HTTP method. Supported methods are: GET, POST, PUT, DELETE", config.Method)
	}

	// Validate Body
	if config.Method == "POST" || config.Method == "PUT" || config.Method == "PATCH" {
		if config.Body == "" {
			return fmt.Errorf("a request body is required for the %s method", config.Method)
		}
		if strings.HasPrefix(config.Body, "@") {
			filePath := strings.TrimPrefix(config.Body, "@")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("the file specified for the request body does not exist: %s", filePath)
			}
		}
	}

	return nil
}

func main() {
	var config BenchmarkConfig

	var rootCmd = &cobra.Command{
		Use:   "api_benchmarker",
		Short: "api_benchmarker is a CLI tool for benchmarking REST APIs",
		Long:  "A Fast and Flexible API benchmarking tool built with help from GPT-4",
		Run: func(cmd *cobra.Command, args []string) {
			runBenchmark(&config) // Pass the config pointer to runBenchmark
		},
	}

	rootCmd.PersistentFlags().StringVarP(&config.URL, "url", "u", "", "The URL of the API endpoint to benchmark.")
	rootCmd.PersistentFlags().StringVarP(&config.Method, "method", "m", "GET", "The HTTP method to use.")
	rootCmd.PersistentFlags().IntVarP(&config.Requests, "requests", "r", 10000, "The number of requests to perform.")
	rootCmd.PersistentFlags().IntVarP(&config.Concurrency, "concurrency", "c", 1000, "The level of concurrency for the requests.")
	rootCmd.PersistentFlags().IntVarP(&config.Duration, "duration", "d", 10, "The duration of the test in seconds.")
	rootCmd.PersistentFlags().StringVarP(&config.Body, "body", "b", "", "The request body for POST/PUT requests. Prefix with @ to point to a file")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return validateFlags(&config)
	}

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
