#!/bin/bash
#
# serve_docs.sh - Run the Jekyll docs site locally for testing
#

set -e

# Change to the docs directory
cd "$(dirname "$0")/../docs"

echo "🌐 Setting up Jekyll environment..."

# Check if Ruby is installed
if ! command -v ruby &> /dev/null; then
    echo "❌ Ruby is not installed. Please install Ruby to run Jekyll."
    exit 1
fi

# Check if Bundler is installed
if ! command -v bundle &> /dev/null; then
    echo "Installing Bundler..."
    gem install bundler >/dev/null
fi

# Install dependencies
echo "Installing dependencies..."
bundle install --quiet

# Run the Jekyll server
echo "🚀 Starting Jekyll server at http://localhost:4000/noidea/"
bundle exec jekyll serve --livereload --baseurl '/noidea' 