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
		// Download content from URL and extract build manifest
		content, err = downloadAndExtractManifest(filePath)
		if err != nil {
			fmt.Printf("Error: Failed to process URL '%s'.\n", filePath)
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

func downloadAndExtractManifest(url string) (string, error) {
	// Create HTTP request to get HTML content
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the HTML response body
	htmlBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert to string for regex processing
	htmlContent := string(htmlBody)

	// Extract build manifest URL from HTML using regex
	// Look for patterns like: <script src="/_next/static/.../build-manifest.json" />
	buildManifestPattern := `<script[^>]*src=["']([^"']*build-manifest\.json[^"']*)["'][^>]*>`
	buildManifestRe, err := regexp.Compile(buildManifestPattern)
	if err != nil {
		return "", fmt.Errorf("failed to compile build manifest regex: %v", err)
	}

	buildManifestMatches := buildManifestRe.FindAllStringSubmatch(htmlContent, -1)

	var buildManifestURL string

	if len(buildManifestMatches) > 0 {
		// Use the first match
		buildManifestURL = buildManifestMatches[0][1]
	} else {
		// Try alternative pattern: look for script tags with build-manifest in the content
		alternativePattern := `<script[^>]*>([^<]*build-manifest\.json[^<]*)</script>`
		altRe, err := regexp.Compile(alternativePattern)
		if err != nil {
			return "", fmt.Errorf("failed to compile alternative regex: %v", err)
		}

		altMatches := altRe.FindAllStringSubmatch(htmlContent, -1)
		if len(altMatches) > 0 {
			// Extract the URL from the content
			for _, match := range altMatches {
				urlMatch := regexp.MustCompile(`["']([^"']*build-manifest\.json[^"']*)["']`)
				urlSubmatches := urlMatch.FindAllStringSubmatch(match[1], -1)
				if len(urlSubmatches) > 0 {
					buildManifestURL = urlSubmatches[0][1]
					break
				}
			}
		}
	}

	// If we found a build manifest URL, download it
	if buildManifestURL != "" {
		// Handle relative URLs by prepending the base URL
		if !strings.HasPrefix(buildManifestURL, "http") {
			baseURL := extractBaseURL(url)
			buildManifestURL = baseURL + buildManifestURL
		}

		// Download build manifest content
		manifestResp, err := http.Get(buildManifestURL)
		if err != nil {
			return "", fmt.Errorf("failed to download build manifest: %v", err)
		}
		defer manifestResp.Body.Close()

		if manifestResp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTP error downloading manifest: %d %s", manifestResp.StatusCode, manifestResp.Status)
		}

		// Read the manifest content
		manifestBody, err := io.ReadAll(manifestResp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read build manifest body: %v", err)
		}

		return string(manifestBody), nil
	} else {
		// If no build manifest found, return the original HTML content
		// This allows users to extract JS paths from HTML directly if needed
		fmt.Println("Warning: No build manifest URL found in HTML. Returning HTML content for JS path extraction.")
		return htmlContent, nil
	}
}

func extractBaseURL(url string) string {
	// Extract base URL (protocol + host + port)
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		return parts[0] + "//" + parts[2]
	}
	return url
}