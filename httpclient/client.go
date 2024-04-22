package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Sends an HTTP request and returns the response body as a string, the status code, and an error if any.
func HttpRequest(method, url string, body io.Reader) (string, int, error) {
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
func GetRequestBody(bodyFlag string) (io.Reader, func(), error) {
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
