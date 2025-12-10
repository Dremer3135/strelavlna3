#!/bin/bash

# This script takes a base path for asset files and a string of LaTeX code,
# compiles it, and returns the path to the generated PNG.

# Check for input
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <base_files_path> \"<latex_content>\""
    exit 1
fi

# Assign arguments to variables for clarity
BASE_FILES_PATH=$1
LATEX_CONTENT=$2

# Create a temporary directory to work in. This is better than mktemp for a file
# as it gives us a clean space for all intermediate files.
TEMP_DIR=$(mktemp -d)

# Define the full paths for our files within the temp directory.
TEX_FILE="${TEMP_DIR}/problem.tex"
DVI_FILE="${TEMP_DIR}/problem.dvi"
PNG_FILE="${TEMP_DIR}/problem.png"

# Ensure the temporary directory and its contents are cleaned up when the script exits.
# trap 'rm -rf "$TEMP_DIR"' EXIT

# If the base files path exists and is a directory, copy its contents.
if [ -d "$BASE_FILES_PATH" ]; then
    # Copy all files from the base path to the temporary directory.
    # This makes images available for the \includegraphics command in LaTeX.
    cp "$BASE_FILES_PATH"/* "$TEMP_DIR/" 2>/dev/null || true
fi

# Write the LaTeX content passed as an argument into the .tex file.
echo "$LATEX_CONTENT" > "$TEX_FILE"

# Run the latex compiler. We run this inside the temporary directory
# so that latex can find the image files.
(cd "$TEMP_DIR" && latex -interaction=nonstopmode "$TEX_FILE" > /dev/null 2>&1)

# Check if the DVI file was created. If not, latex failed.
if [ ! -f "$DVI_FILE" ]; then
    echo "Error: latex compilation failed. No DVI file produced." >&2
    # For debugging, you can check the log file in the temp dir.
    # cat "${TEMP_DIR}/problem.log" >&2
    exit 1
fi

# Convert the DVI file to a tight-fitting PNG.
# We also run this from inside the temp directory.
(cd "$TEMP_DIR" && dvipng -D 149 -T tight -o "$PNG_FILE" "$DVI_FILE" > /dev/null 2>&1)

# Check if the PNG file was created.
if [ ! -f "$PNG_FILE" ]; then
    echo "Error: dvipng failed. No PNG file produced." >&2
    exit 1
fi

# The final action is to print the path of the generated PNG to standard output.
# The Go program will read this from the script's stdout.
echo "$PNG_FILE"
