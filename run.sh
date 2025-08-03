#!/bin/bash

# Set the secret file path - use the environment variable or default
json_file_path="${PATH_TO_SECRET_FILE:-secret_file_dev.json}"

echo "Using secret file: $json_file_path"

# Check if the JSON file exists
if [ ! -f "$json_file_path" ]; then
    echo "Error: JSON file not found at '$json_file_path'"
    exit 1
fi

eval $(go run configs/env/secret_base64.go -p $json_file_path)
export NOT_DEBUG=true
echo "================================================"
echo "	Finished export all env"
echo " SECRET_MANAGER_BASE64: $SECRET_MANAGER_BASE64"
echo "================================================"
echo "Done export all env"


# Run the compiled binary
go run cmd/main.go
