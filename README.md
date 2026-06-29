# next

A python regex based script to help out extract all js paths from buildmanifest in NextJs based apps.

## Go Version

This project now includes a Go version of the tool for improved performance and portability.

### Usage

#### Python Version
```bash
python3 next.py <file_path_to_your_buildManifest.js>
```

#### Go Version
```bash
# For local file
go run next.go <file_path_to_your_buildManifest.js>

# For URL (downloads HTML and extracts build manifest)
go run next.go -u <url_to_your_nextjs_app>
```

## Features

- Extracts JavaScript file paths from Next.js build manifest files
- Uses regex pattern matching to find static JS files
- Handles file reading errors gracefully
- Command-line interface for easy usage
- Supports both local files and remote URLs (downloads HTML content)

## Example Output
```
static/abc123.js
static/def456.js
static/ghi789.js
```

## Requirements

- Python 3.x (for Python version)
- Go 1.16+ (for Go version)
