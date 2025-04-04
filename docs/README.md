# noidea Documentation

This directory contains the source files for the noidea project documentation website, built with Jekyll and deployed to GitHub Pages.

## Local Development

To run the documentation site locally:

1. Make sure you have Ruby and Bundler installed
2. Run the helper script from the project root:
   ```bash
   ./scripts/serve_docs.sh
   ```
3. Open your browser at http://localhost:4000/noidea/

## Adding Content

To add new documentation:

1. Create a new Markdown file (`.md`) in the docs directory
2. Add the front matter at the top:
   ```yaml
   ---
   layout: default
   title: Your Page Title
   ---
   ```
3. Add your content using Markdown
4. Link to your new page from other pages as needed

## Deployment

The documentation site is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The deployment is handled by the GitHub workflow defined in `.github/workflows/pages.yml`.

## Structure

- `_config.yml` - Jekyll configuration
- `Gemfile` - Ruby dependencies
- `.gitignore` - Files to be ignored by Git
- `index.md` - Home page
- `*.md` - Documentation pages
- `assets/` - Images and other static assets 