package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/komuvill/api_benchmarker/benchmark"
	"github.com/komuvill/api_benchmarker/metrics"
	"github.com/komuvill/api_benchmarker/report"
	"github.com/komuvill/api_benchmarker/storage"
	"github.com/spf13/cobra"
)

func NewRootCmd(config *benchmark.BenchmarkConfig) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "api_benchmarker",
		Short: "api_benchmarker is a CLI tool for benchmarking REST APIs",
		Long:  "A Fast and Flexible API benchmarking tool built with help from GPT-4",
		Run: func(cmd *cobra.Command, args []string) {
			executeBenchmark(config)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&config.URL, "url", "u", "", "The URL of the API endpoint to benchmark.")
	rootCmd.PersistentFlags().StringVarP(&config.Method, "method", "m", "GET", "The HTTP method to use. Accepted methods: GET, POST, PUT, DELETE")
	rootCmd.PersistentFlags().IntVarP(&config.Requests, "requests", "r", 10000, "The number of requests to perform.")
	rootCmd.PersistentFlags().IntVarP(&config.Concurrency, "concurrency", "c", 1000, "The level of concurrency for the requests.")
	rootCmd.PersistentFlags().IntVarP(&config.Duration, "duration", "d", 10, "The duration of the test in seconds.")
	rootCmd.PersistentFlags().StringVarP(&config.Body, "body", "b", "", "The request body for POST/PUT requests. Prefix with @ to point to a file. Currently only json formatted bodies are accepted")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return validateFlags(config)
	}

	return rootCmd
}

func validateFlags(config *benchmark.BenchmarkConfig) error {
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

func executeBenchmark(config *benchmark.BenchmarkConfig) {
	startTime := time.Now()
	results := benchmark.RunBenchmark(config)
	aggregatedMetrics := metrics.CalculateMetrics(results)
	metrics.PrintMetrics(aggregatedMetrics)

	outputDir := "./output"
	os.MkdirAll(outputDir, os.ModePerm)
	storage.SaveResults(results, outputDir)
	storage.SaveAggregatedMetrics(aggregatedMetrics, outputDir)

	err := report.GenerateHTMLReport(*config, aggregatedMetrics, results, startTime, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating HTML report: %v\n", err)
		os.Exit(1)
	}
}
