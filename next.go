package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Check if the file path or URL is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run next.go [-u] <file_path_or_url_to_your_buildManifest.js>")
		fmt.Println("  -u : Optional flag to specify a URL instead of a local file")
		os.Exit(1)
	}

	var filePath string
	var isURL bool

	// Check if -u flag is provided
	if len(os.Args) >= 3 && os.Args[1] == "-u" {
		isURL = true
		filePath = os.Args[2]
	} else {
		filePath = os.Args[1]
	}

	// Define the regex pattern for JavaScript paths
	pattern := `static/[^\"]+\.js`

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error: Failed to compile regex pattern '%s'.\n", pattern)
		os.Exit(1)
	}

	var content string

	if isURL {
		// Download content from URL
		content, err = downloadContent(filePath)
		if err != nil {
			fmt.Printf("Error: Failed to download content from URL '%s'.\n", filePath)
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// Read the file content
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error: The file '%s' could not be read.\n", filePath)
			os.Exit(1)
		}
		content = string(fileContent)
	}

	// Perform the regex search
	matches := re.FindAllString(content, -1)

	// Print the extracted JavaScript paths
	for _, match := range matches {
		fmt.Println(match)
	}
}

func downloadContent(url string) (string, error) {
	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert to string and return
	return string(body), nil
}