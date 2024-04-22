package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/komuvill/api_benchmarker/metrics"
)

// SaveResults serializes the slice of RequestResults to JSON and saves it to a file.
// generateTimestampedFilename creates a filename with a timestamp.
func generateTimestampedFilename(prefix string) string {
	timestamp := time.Now().Format("01022006-150405") // DDMMYYYY-HHMMSS format
	return fmt.Sprintf("%s-%s.json", prefix, timestamp)
}

// RequestResultForStorage is a struct for storing RequestResult with a serializable Error field.
type RequestResultForStorage struct {
	RequestID    int           `json:"request_id"`
	Response     string        `json:"response"`
	StatusCode   int           `json:"status_code"`
	ResponseTime time.Duration `json:"response_time"`
	Error        string        `json:"error,omitempty"`
}

// ConvertRequestResults prepares a slice of RequestResult for storage by converting the Error field.
func ConvertRequestResults(results []metrics.RequestResult) []RequestResultForStorage {
	storageResults := make([]RequestResultForStorage, len(results))
	for i, result := range results {
		storageResults[i] = RequestResultForStorage{
			RequestID:    result.RequestID,
			Response:     result.Response,
			StatusCode:   result.StatusCode,
			ResponseTime: result.ResponseTime,
			Error:        "", // Default empty string if there's no error
		}
		if result.Error != nil {
			storageResults[i].Error = result.Error.Error() // Convert the error to a string
		}
	}
	return storageResults
}

// SaveResults serializes the slice of RequestResults to JSON and saves it to a file with a timestamp.
func SaveResults(results []metrics.RequestResult, outputDir string) error {
	// Convert the results to a storage-friendly format.
	storageResults := ConvertRequestResults(results)

	// Create the output file with a timestamped filename.
	filename := generateTimestampedFilename("results")
	filePath := filepath.Join(outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use the JSON encoder to write the results to the file with indentation for readability.
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(storageResults)
}

// SaveAggregatedMetrics serializes the AggregateMetrics to JSON and saves it to a file with a timestamp.
func SaveAggregatedMetrics(metrics metrics.AggregateMetrics, outputDir string) error {
	// Create the output file with a timestamped filename.
	filename := generateTimestampedFilename("aggregated")
	filePath := filepath.Join(outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use the JSON encoder to write the metrics to the file with indentation for readability.
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metrics)
}
