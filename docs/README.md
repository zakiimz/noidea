# noidea Documentation

This directory contains the GitHub Pages documentation site for the noidea project.

## Local Development

To run this site locally:

1. Install Ruby and Bundler
2. Run `bundle install` to install dependencies
3. Run `bundle exec jekyll serve` to start the local server
4. Visit `http://localhost:4000/noidea/` in your browser

## Structure

- `_config.yml`: Jekyll configuration
- `Gemfile`: Ruby dependencies
- `index.md`: Home page
- `docs/`: Documentation pages
  - `getting-started.md`: Getting started guide
  - `usage.md`: Usage documentation
  - `faq.md`: Frequently asked questions
  - `configuration.md`: Advanced configuration
- `assets/`: Images and other assets

## Deployment

The site is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The deployment is handled by the GitHub workflow defined in `.github/workflows/docs.yml`. 