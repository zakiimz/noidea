#!/bin/bash
# setup-blog.sh - Set up blog dependencies for NoIdea documentation
# This script installs the necessary Python packages for the blog

echo "Setting up blog dependencies for NoIdea documentation..."

# Check if pip is installed
if ! command -v pip &> /dev/null; then
    echo "Error: pip is not installed. Please install Python and pip first."
    exit 1
fi

# Install MkDocs and required plugins
echo "Installing MkDocs and required plugins..."
pip install mkdocs-material
pip install mkdocs-rss-plugin

# Check if installation was successful
if [ $? -eq 0 ]; then
    echo "Successfully installed blog dependencies!"
    echo ""
    echo "You can now build the documentation with:"
    echo "  mkdocs build"
    echo ""
    echo "Or serve it locally with:"
    echo "  mkdocs serve"
    echo ""
    echo "Your blog is available at: /blog/"
else
    echo "Error: Failed to install dependencies."
    exit 1
fi

# Make script executable
chmod +x "$(dirname "$0")/setup-blog.sh" 