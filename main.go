package main

import (
	"fmt"
	"os"

	"github.com/komuvill/api_benchmarker/benchmark"
	"github.com/komuvill/api_benchmarker/cli"
)

func main() {
	var config benchmark.BenchmarkConfig

	rootCmd := cli.NewRootCmd(&config)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
