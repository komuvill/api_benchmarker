package main

import (
	"fmt"
    "io"
    "io/ioutil"
    "net/http"
	"github.com/spf13/cobra"
)

func runBenchmark(cmd *cobra.Command, args []string) {
	// Retrieve the flag values
	url, _ := cmd.Flags().GetString("url")
	method, _ := cmd.Flags().GetString("method")
	requests, _ := cmd.Flags().GetInt("requests")
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	duration, _ := cmd.Flags().GetInt("duration")

	fmt.Printf("Benchmarking %s with %s method, %d requests, %d concurrent requests, for %d seconds\n", url, method, requests, concurrency, duration)

    //TODO: construct responseBody if possible

    responseBody, statusCode, err := httpRequest(method, url, nil)

    if err != nil {
        fmt.Printf("%v\n", err)
        return
    }
    
    fmt.Printf("Response Status Code: %d\n", statusCode)
    fmt.Printf("Response Body: %s\n", responseBody)
}

func httpRequest(method, url string, body io.Reader) (string, int, error) {
    // Create a new request with the method, URL, and body if provided
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return "", 0, fmt.Errorf("error creating request: %v", err)
    }

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", 0, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    // Read the response body
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", resp.StatusCode, fmt.Errorf("error reading response body: %v", err)
    }

    return string(respBody), resp.StatusCode, nil
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

    //TODO: If method is POST/PUT/DELETE, make sure body is also supplied?

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
    //TODO: Flag for defining body for requests like POST

	rootCmd.PreRunE = validateFlags

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
