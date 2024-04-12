package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runBenchmark(cmd *cobra.Command, args []string) {
	// Retrieve the flag values
	url, _ := cmd.Flags().GetString("url")
	method, _ := cmd.Flags().GetString("method")
	requests, _ := cmd.Flags().GetInt("requests")
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	duration, _ := cmd.Flags().GetInt("duration")

	// TODO: Implement the benchmarking logic
	fmt.Printf("Benchmarking %s with %s method, %d requests, %d concurrent requests, for %d seconds\n", url, method, requests, concurrency, duration)
}

func validateFlags(cmd *cobra.Command, args []string) error {

	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}

	if url == "" {
		return fmt.Errorf("URL is required")
	}

	method, err := cmd.Flags().GetString("method")
	if err != nil {
		return err
	}

	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
	}

	if _, valid := validMethods[method]; !valid {
		return fmt.Errorf("'%s' is not a valid HTTP method. Following HTTP methods are supported: GET, POST, PUT, DELETE", method)
	}

	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "api_benchmarker",
		Short: "api_benchmarker is a CLI tool for benchmarking REST APIs",
		Long:  "A Fast and Flexible API benchmarking tool built with help from GPT-4",
		Run:   runBenchmark,
	}

	var (
		url         string
		method      string
		requests    int
		concurrency int
		duration    int
	)

	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The URL of the API endpoint to benchmark.")
	rootCmd.PersistentFlags().StringVarP(&method, "method", "m", "GET", "The HTTP method to use.")
	rootCmd.PersistentFlags().IntVarP(&requests, "requests", "r", 100, "The number of requests to perform.")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 10, "The level of concurrency for the requests.")
	rootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 10, "The duration of the test in seconds.")

	rootCmd.PreRunE = validateFlags

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
