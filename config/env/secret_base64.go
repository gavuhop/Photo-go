package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
)

func main() {
	filePath := flag.String("p", "", "Path to the file")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path using -p flag")
		return
	}

	// Read the file contents
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(data)

	// Print the export command
	fmt.Printf("export SECRET_MANAGER_BASE64=%s\n", encoded)
}
