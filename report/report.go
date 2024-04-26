package report

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/komuvill/api_benchmarker/benchmark"
	"github.com/komuvill/api_benchmarker/metrics"
)

//go:embed report_template.html
var reportTemplate embed.FS

// ReportData holds all the data necessary for the report
type ReportData struct {
	Config           benchmark.BenchmarkConfig
	StartTime        string
	AggregateMetrics metrics.AggregateMetrics
	RequestResults   []metrics.RequestResult
}

func GenerateHTMLReport(config benchmark.BenchmarkConfig, aggregateMetrics metrics.AggregateMetrics, requestResults []metrics.RequestResult, startTime time.Time, outputDir string) error {
	// Create a ReportData struct with all necessary data
	data := ReportData{
		Config:           config,
		StartTime:        startTime.Format(time.RFC1123), // Format the start time as a string
		AggregateMetrics: aggregateMetrics,
		RequestResults:   requestResults,
	}

	// Define name and output path for the report
	filenameTimestamp := startTime.Format("020106-150405") // DDMMYY-HHMMSS format
	fileName := fmt.Sprintf("%s_benchmark_report.html", filenameTimestamp)
	filePath := filepath.Join(outputDir, fileName)

	// Parse the HTML template from the embedded file system
	tmpl, err := template.ParseFS(reportTemplate, "report_template.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create an output file for the report
	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating report file: %w", err)
	}
	defer outputFile.Close()

	// Execute the template and write the report to the file
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
