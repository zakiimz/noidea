#!/bin/bash
#
# noidea - Version management script
#
# This script helps manage versioning for noidea
# It can display the current version or bump the version according to semantic versioning
#

set -e

# Text formatting
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Current script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Show help message
usage() {
    echo -e "${BLUE}noidea Version Manager${NC}"
    echo "============================="
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  show       Display current version"
    echo "  major      Bump major version (x.0.0)"
    echo "  minor      Bump minor version (0.x.0)"
    echo "  patch      Bump patch version (0.0.x)"
    echo "  help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 show    # Show current version"
    echo "  $0 patch   # Bump patch version"
    exit 1
}

# Get current version from git tag
get_current_version() {
    latest_tag=$(git describe --tags $(git rev-list --tags --max-count=1) 2>/dev/null || echo "v0.0.0")
    echo "$latest_tag"
}

# Update version in cmd/root.go
update_version_in_file() {
    local new_version=$1
    sed -i "s/Version   = \"[^\"]*\"/Version   = \"$new_version\"/" "$ROOT_DIR/cmd/root.go"
    echo -e "${GREEN}✓${NC} Updated version in cmd/root.go to $new_version"
}

# Bump version according to semantic versioning
bump_version() {
    local bump_type=$1
    local current_version=$(get_current_version)
    
    # Remove 'v' prefix
    version=${current_version#v}
    
    # Split version into components
    IFS='.' read -r -a version_parts <<< "$version"
    major=${version_parts[0]:-0}
    minor=${version_parts[1]:-0}
    patch=${version_parts[2]:-0}
    
    # Bump version according to type
    case "$bump_type" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo -e "${RED}Error: Invalid bump type: $bump_type${NC}"
            usage
            ;;
    esac
    
    # Create new version
    new_version="v$major.$minor.$patch"
    echo -e "${CYAN}Bumping version: $current_version → $new_version${NC}"
    
    # Update version in files
    update_version_in_file "$new_version"
    
    # Prompt for git commit and tag
    echo ""
    read -p "Create git commit and tag for this version? (y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git add "$ROOT_DIR/cmd/root.go"
        git commit -m "Bump version to $new_version"
        git tag -a "$new_version" -m "Release $new_version"
        echo -e "${GREEN}✓${NC} Created commit and tag $new_version"
        
        echo ""
        read -p "Push changes and tag to remote repository? (y/n): " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git push
            git push origin "$new_version"
            echo -e "${GREEN}✓${NC} Pushed changes to remote repository"
            
            # Check if noidea is installed and GitHub integration is available
            if command -v noidea >/dev/null 2>&1; then
                # Check if user has GitHub token
                if noidea github status >/dev/null 2>&1; then
                    echo ""
                    echo -e "${CYAN}Generating enhanced release notes...${NC}"
                    noidea github release notes --tag="$new_version"
                else
                    echo ""
                    echo -e "${YELLOW}Note:${NC} To generate enhanced AI release notes, run:"
                    echo "  noidea github auth        # Authenticate with GitHub"
                    echo "  noidea github release notes --tag=$new_version  # Generate release notes"
                fi
            else
                echo ""
                echo -e "${YELLOW}Note:${NC} To generate enhanced AI release notes, make sure noidea is installed and run:"
                echo "  noidea github release notes --tag=$new_version"
            fi
        else
            echo -e "${YELLOW}Remember to push your changes:${NC}"
            echo "  git push && git push origin $new_version"
            echo ""
            echo -e "${YELLOW}After pushing, you can generate enhanced release notes:${NC}"
            echo "  noidea github release notes --tag=$new_version"
        fi
    else
        echo -e "${YELLOW}Changes made locally. Don't forget to commit and create a tag.${NC}"
    fi
}

# Main execution
if [ $# -eq 0 ]; then
    usage
fi

command=$1

case "$command" in
    show)
        current_version=$(get_current_version)
        echo -e "${CYAN}Current version: $current_version${NC}"
        ;;
    major|minor|patch)
        bump_version "$command"
        ;;
    help|--help|-h)
        usage
        ;;
    *)
        echo -e "${RED}Error: Unknown command: $command${NC}"
        usage
        ;;
esac 