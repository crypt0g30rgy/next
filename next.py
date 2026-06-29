#!/usr/bin/env python3

import os
import re
import sys
import urllib.request
import urllib.parse
from urllib.error import URLError, HTTPError

def main():
    # Check if the file path or URL is provided as an argument
    if len(sys.argv) < 2:
        print("Usage: python3 next.py [-u] <file_path_or_url_to_your_buildManifest.js>")
        print("  -u : Optional flag to specify a URL instead of a local file")
        sys.exit(1)

    # Check if -u flag is provided
    is_url = False
    if len(sys.argv) >= 3 and sys.argv[1] == "-u":
        is_url = True
        file_path = sys.argv[2]
    else:
        file_path = sys.argv[1]

    # Define the regex pattern for JavaScript paths
    pattern = r'static/[^\"]+\.js'

    # Compile the regex
    try:
        compiled_pattern = re.compile(pattern)
    except re.error as e:
        print(f"Error: Failed to compile regex pattern '{pattern}'.")
        sys.exit(1)

    content = ""

    if is_url:
        # Download content from URL and extract build manifest
        try:
            content = download_and_extract_manifest(file_path)
        except Exception as e:
            print(f"Error: Failed to process URL '{file_path}'.")
            print(e)
            sys.exit(1)
    else:
        # Get the file path from the command-line arguments
        # Check if the file exists
        if not os.path.isfile(file_path):
            print(f"Error: The file '{file_path}' does not exist.")
            sys.exit(1)

        try:
            # Open the file in read mode
            with open(file_path, 'r') as file:
                # Read the entire content of the file
                content = file.read()
        except Exception as e:
            print(f"Error: The file '{file_path}' could not be read.")
            print(e)
            sys.exit(1)

    # Perform the regex search
    matches = compiled_pattern.findall(content)

    # Print the extracted JavaScript paths
    for match in matches:
        print(match)

def download_and_extract_manifest(url):
    """Download HTML from URL, extract build manifest URL, and return manifest content"""

    # Download HTML content
    try:
        response = urllib.request.urlopen(url)
        html_content = response.read().decode('utf-8')
    except (URLError, HTTPError) as e:
        raise Exception(f"Failed to download HTML from URL '{url}': {e}")

    # Extract build manifest URL from HTML using regex
    # Look for patterns like: <script src="/_next/static/.../build-manifest.json" />
    build_manifest_pattern = r'<script[^>]*src=["\']([^"\']*build-manifest\.json[^"\']*)["\'][^>]*>'

    try:
        build_manifest_matches = re.findall(build_manifest_pattern, html_content)
    except re.error as e:
        raise Exception(f"Failed to compile build manifest regex: {e}")

    if build_manifest_matches:
        # Use the first match
        build_manifest_url = build_manifest_matches[0]

        # Handle relative URLs by prepending the base URL
        if not build_manifest_url.startswith('http'):
            base_url = extract_base_url(url)
            build_manifest_url = base_url + build_manifest_url

        # Download build manifest content
        try:
            manifest_response = urllib.request.urlopen(build_manifest_url)
            manifest_content = manifest_response.read().decode('utf-8')
            return manifest_content
        except (URLError, HTTPError) as e:
            raise Exception(f"Failed to download build manifest: {e}")
    else:
        # If no build manifest found, return the original HTML content
        # This allows users to extract JS paths from HTML directly if needed
        print("Warning: No build manifest URL found in HTML. Returning HTML content for JS path extraction.")
        return html_content

def extract_base_url(url):
    """Extract base URL (protocol + host + port) from a full URL"""
    parsed_url = urllib.parse.urlparse(url)
    return f"{parsed_url.scheme}://{parsed_url.netloc}"

if __name__ == "__main__":
    main()
