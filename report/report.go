package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/komuvill/api_benchmarker/benchmark"
	"github.com/komuvill/api_benchmarker/metrics"
)

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

	// Define the path to the HTML template
	templatePath := filepath.Join("report", "report_template.html")

	// Define name and output path for the report
	filenameTimestamp := startTime.Format("020106-150405") // DDMMYY-HHMMSS format
	fileName := fmt.Sprintf("%s_benchmark_report.html", filenameTimestamp)
	filePath := filepath.Join(outputDir, fileName)

	// Parse the HTML template from file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	// Create an output file for the report

	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Execute the template and write the report to the file
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return err
	}

	return nil
}
