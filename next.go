package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	// Check if the file path is provided as an argument
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run next.go <file_path_to_your_buildManifest.js>")
		os.Exit(1)
	}

	// Get the file path from the command-line arguments
	filePath := os.Args[1]

	// Define the regex pattern for JavaScript paths
	pattern := `static/[^\"]+\.js`

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error: Failed to compile regex pattern '%s'.\n", pattern)
		os.Exit(1)
	}

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error: The file '%s' could not be read.\n", filePath)
		os.Exit(1)
	}

	// Perform the regex search
	matches := re.FindAllString(string(content), -1)

	// Print the extracted JavaScript paths
	for _, match := range matches {
		fmt.Println(match)
	}
}