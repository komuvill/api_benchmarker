package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func runBenchmark(cmd *cobra.Command, args []string) {
	// Retrieve the flag values
	url, _ := cmd.Flags().GetString("url")
	method, _ := cmd.Flags().GetString("method")
	requests, _ := cmd.Flags().GetInt("requests")
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	duration, _ := cmd.Flags().GetInt("duration")
	body, _ := cmd.Flags().GetString("body")

	fmt.Printf("Benchmarking %s with %s method, %d requests, %d concurrent requests, for %d seconds\n", url, method, requests, concurrency, duration)

	responseBody, statusCode, err := httpRequest(method, url, body)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Response Status Code: %d\n", statusCode)
	fmt.Printf("Response Body: %s\n", responseBody)
}

// httpRequest sends an HTTP request and returns the response body as a string, the status code, and an error if any.
func httpRequest(method, url string, bodyFlag string) (string, int, error) {
    var body io.Reader
    var err error
    var cleanup func()

    if method == "POST" || method == "PUT" || method == "PATCH" {
        body, cleanup, err = getRequestBody(bodyFlag)
        if err != nil {
            return "", 0, err
        }
        // Defer the cleanup function to close the file when done
        defer cleanup()
    }

    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return "", 0, fmt.Errorf("error creating request: %v", err)
    }

    // Set the Content-Type header if there is a body.
	// For purpose of the project, assume the API to be tested expects json requests.
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", 0, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
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

	body, err := cmd.Flags().GetString("body")
	if err != nil {
		return err
	}

	if method == "POST" || method == "PUT" {
		if body == "" {
			return fmt.Errorf("a request body is required for the %s method", method)
		}
	}

	if strings.HasPrefix(body, "@") {
		filePath := strings.TrimPrefix(body, "@")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("the file specified for the request body does not exist: %s", filePath)
		}
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
		body        string
	)

	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The URL of the API endpoint to benchmark.")
	rootCmd.PersistentFlags().StringVarP(&method, "method", "m", "GET", "The HTTP method to use.")
	rootCmd.PersistentFlags().IntVarP(&requests, "requests", "r", 100, "The number of requests to perform.")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 10, "The level of concurrency for the requests.")
	rootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 10, "The duration of the test in seconds.")
	rootCmd.PersistentFlags().StringVarP(&body, "body", "b", "", "The request body for POST/PUT requests. Prefix with @ to point to a file")

	rootCmd.PreRunE = validateFlags

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
