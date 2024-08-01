#!/usr/bin/env python3

import os
import re
import sys

def main():
    # Check if the file path is provided as an argument
    if len(sys.argv) != 2:
        print("Usage: python3 next.py <file_path_to_your_`buildManifest.js`>")
        sys.exit(1)
    
    # Get the file path from the command-line arguments
    file_path = sys.argv[1]
    
    # Define the regex pattern for JavaScript paths
    pattern = r'static/[^\"]+\.js'

    # Check if the file exists
    if not os.path.isfile(file_path):
        print(f"Error: The file '{file_path}' does not exist.")
        sys.exit(1)
    
    try:
        # Open the file in read mode
        with open(file_path, 'r') as file:
            # Read the entire content of the file
            content = file.read()
        
        # Perform the regex search
        matches = re.findall(pattern, content)
        
        # Print the extracted JavaScript paths
        for match in matches:
            print(match)
    
    except Exception as e:
        print(f"Error: An error occurred while reading the file '{file_path}'.")
        print(e)
        sys.exit(1)

if __name__ == "__main__":
    main()
