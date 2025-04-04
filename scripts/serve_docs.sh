#!/bin/bash

# Script to serve the Jekyll site locally for testing

# Navigate to the docs directory
cd "$(dirname "$0")/../docs" || exit

# Install dependencies if needed
if ! bundle check > /dev/null; then
  echo "Installing dependencies..."
  bundle install
fi

# Build and serve the site
echo "Starting Jekyll server..."
echo "Visit http://localhost:4000/noidea/ in your browser"
bundle exec jekyll serve --baseurl "/noidea" 