#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# The name of the output binary
OUTPUT_BINARY="api_benchmarker"

# The directory where the binary will be placed
OUTPUT_DIR="${GOPATH}/bin"

# Create the output directory if it doesn't exist
mkdir -p ${OUTPUT_DIR}

# Build the application
echo "Building the application..."
go build -o ${OUTPUT_DIR}/${OUTPUT_BINARY}

# Give execution permissions to the binary
chmod +x ${OUTPUT_DIR}/${OUTPUT_BINARY}

echo "Build complete. Binary located at ${OUTPUT_DIR}/${OUTPUT_BINARY}"
