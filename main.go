package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/komuvill/api_benchmarker/metrics"
	"github.com/komuvill/api_benchmarker/benchmark"
	"github.com/komuvill/api_benchmarker/storage"
)

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

func main() {
	var config benchmark.BenchmarkConfig

	var rootCmd = &cobra.Command{
		Use:   "api_benchmarker",
		Short: "api_benchmarker is a CLI tool for benchmarking REST APIs",
		Long:  "A Fast and Flexible API benchmarking tool built with help from GPT-4",
		Run: func(cmd *cobra.Command, args []string) {
			results := benchmark.RunBenchmark(&config)
			aggregatedMetrics := metrics.CalculateMetrics(results)
			metrics.PrintMetrics(aggregatedMetrics)
			outputDir := "./output"
			os.MkdirAll(outputDir, os.ModePerm)
			storage.SaveResults(results, outputDir)
			storage.SaveAggregatedMetrics(aggregatedMetrics, outputDir)

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
