#!/bin/bash
# migrate_to_secure.sh
# Script to migrate API keys from .env files to secure storage
# Part of the noidea project

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}Noidea API Key Migration Tool${NC}"
echo "This script will help migrate your API keys from .env files to secure storage"
echo

# Check for noidea binary
if ! command -v noidea &> /dev/null; then
    echo -e "${RED}Error: noidea binary not found in PATH${NC}"
    echo "Please make sure noidea is installed and in your PATH"
    exit 1
fi

# Potential .env file locations
HOME_DIR="$HOME"
LOCATIONS=(
    "./.env"
    "./.noidea.env"
    "$HOME_DIR/.noidea/.env"
)

# Check for API keys in environment variables
API_KEY=""
PROVIDER=""

if [ -n "$XAI_API_KEY" ]; then
    API_KEY="$XAI_API_KEY"
    PROVIDER="xai"
    echo -e "${YELLOW}Found XAI API key in environment variables${NC}"
elif [ -n "$OPENAI_API_KEY" ]; then
    API_KEY="$OPENAI_API_KEY"
    PROVIDER="openai"
    echo -e "${YELLOW}Found OpenAI API key in environment variables${NC}"
elif [ -n "$DEEPSEEK_API_KEY" ]; then
    API_KEY="$DEEPSEEK_API_KEY"
    PROVIDER="deepseek"
    echo -e "${YELLOW}Found DeepSeek API key in environment variables${NC}"
elif [ -n "$NOIDEA_API_KEY" ]; then
    API_KEY="$NOIDEA_API_KEY"
    echo -e "${YELLOW}Found generic API key in environment variables${NC}"
    
    # Ask for provider
    echo "What provider is this key for?"
    echo "1. xAI (Grok)"
    echo "2. OpenAI"
    echo "3. DeepSeek"
    read -p "Provider (1-3): " provider_choice
    
    case $provider_choice in
        1) PROVIDER="xai" ;;
        2) PROVIDER="openai" ;;
        3) PROVIDER="deepseek" ;;
        *) PROVIDER="xai" ;; # Default to xAI
    esac
fi

# If no API key found in environment, check .env files
if [ -z "$API_KEY" ]; then
    for location in "${LOCATIONS[@]}"; do
        if [ -f "$location" ]; then
            echo -e "${YELLOW}Found .env file at $location${NC}"
            
            # Check for API keys
            if grep -q "XAI_API_KEY" "$location"; then
                API_KEY=$(grep "XAI_API_KEY" "$location" | cut -d= -f2)
                PROVIDER="xai"
                echo "Found xAI API key"
            elif grep -q "OPENAI_API_KEY" "$location"; then
                API_KEY=$(grep "OPENAI_API_KEY" "$location" | cut -d= -f2)
                PROVIDER="openai"
                echo "Found OpenAI API key"
            elif grep -q "DEEPSEEK_API_KEY" "$location"; then
                API_KEY=$(grep "DEEPSEEK_API_KEY" "$location" | cut -d= -f2)
                PROVIDER="deepseek"
                echo "Found DeepSeek API key"
            elif grep -q "NOIDEA_API_KEY" "$location"; then
                API_KEY=$(grep "NOIDEA_API_KEY" "$location" | cut -d= -f2)
                echo "Found generic API key"
                
                # Ask for provider
                echo "What provider is this key for?"
                echo "1. xAI (Grok)"
                echo "2. OpenAI"
                echo "3. DeepSeek"
                read -p "Provider (1-3): " provider_choice
                
                case $provider_choice in
                    1) PROVIDER="xai" ;;
                    2) PROVIDER="openai" ;;
                    3) PROVIDER="deepseek" ;;
                    *) PROVIDER="xai" ;; # Default to xAI
                esac
            fi
            
            # Remove quotes if present
            API_KEY=$(echo "$API_KEY" | tr -d '"' | tr -d "'" | xargs)
            
            if [ -n "$API_KEY" ]; then
                break
            fi
        fi
    done
fi

# If still no API key found, prompt the user
if [ -z "$API_KEY" ]; then
    echo -e "${YELLOW}No API key found in environment variables or .env files${NC}"
    echo "Please provide your API key manually:"
    
    # Ask for provider
    echo "What provider is this key for?"
    echo "1. xAI (Grok)"
    echo "2. OpenAI"
    echo "3. DeepSeek"
    read -p "Provider (1-3): " provider_choice
    
    case $provider_choice in
        1) PROVIDER="xai" ;;
        2) PROVIDER="openai" ;;
        3) PROVIDER="deepseek" ;;
        *) PROVIDER="xai" ;; # Default to xAI
    esac
    
    # Read API key securely
    read -sp "Enter your API key: " API_KEY
    echo
    
    if [ -z "$API_KEY" ]; then
        echo -e "${RED}No API key provided. Exiting.${NC}"
        exit 1
    fi
fi

# Ensure we have both provider and API key
if [ -z "$PROVIDER" ] || [ -z "$API_KEY" ]; then
    echo -e "${RED}Error: Missing provider or API key${NC}"
    exit 1
fi

# Before storing the key, validate it
echo -e "${BLUE}Validating API key...${NC}"

# Validate the API key against the provider
validate_api_key() {
    local provider=$1
    local api_key=$2
    local url=""
    
    case "$provider" in
        xai)
            url="https://api.groq.com/v1/models"
            ;;
        openai)
            url="https://api.openai.com/v1/models"
            ;;
        deepseek)
            url="https://api.deepseek.com/v1/models"
            ;;
        *)
            echo "Unknown provider: $provider"
            return 1
            ;;
    esac
    
    # Make a request to check if the key is valid
    response=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $api_key" $url)
    
    if [[ "$response" =~ ^2[0-9][0-9]$ ]]; then
        # 2xx status code means success
        return 0
    else
        # Non-2xx status code means failure
        return 1
    fi
}

if validate_api_key "$PROVIDER" "$API_KEY"; then
    echo -e "${GREEN}API key is valid!${NC}"
else
    echo -e "${RED}Warning: API key appears to be invalid${NC}"
    echo "The key was rejected by the provider. It may be incorrect or expired."
    
    # Ask if user wants to continue
    read -p "Continue with migration anyway? [Y/n]: " continue_anyway
    if [[ "$continue_anyway" =~ ^[Nn]$ ]]; then
        echo "Migration cancelled."
        exit 1
    fi
fi

echo -e "${BLUE}Storing API key securely...${NC}"

# Use noidea command to store the key securely
# First export to environment for the command to pick up
export NOIDEA_API_KEY="$API_KEY"
export NOIDEA_LLM_PROVIDER="$PROVIDER"

# Call noidea to store the key securely
if noidea config apikey; then
    echo -e "${GREEN}API key stored securely!${NC}"
    
    # Ask if user wants to remove keys from .env files
    echo
    read -p "Would you like to remove API keys from .env files? [y/N]: " remove_env
    
    if [[ "$remove_env" =~ ^[Yy]$ ]]; then
        for location in "${LOCATIONS[@]}"; do
            if [ -f "$location" ]; then
                # Create backup
                cp "$location" "${location}.bak"
                echo "Created backup at ${location}.bak"
                
                # Remove API key lines
                sed -i '/XAI_API_KEY=/d' "$location"
                sed -i '/OPENAI_API_KEY=/d' "$location"
                sed -i '/DEEPSEEK_API_KEY=/d' "$location"
                sed -i '/NOIDEA_API_KEY=/d' "$location"
                
                echo "Removed API keys from $location"
            fi
        done
        echo -e "${GREEN}API keys removed from .env files${NC}"
    fi
    
    echo
    echo -e "${GREEN}Migration complete!${NC}"
    echo "Your API key is now securely stored and can be managed with:"
    echo "  noidea config apikey        # Update API key"
    echo "  noidea config apikey-status # Check status"
    echo "  noidea config apikey-remove # Remove API key"
    echo
    echo "See docs/api-key-management.md for more information"
else
    echo -e "${RED}Failed to store API key securely${NC}"
    echo "You can try running the command manually:"
    echo "  noidea config apikey"
    exit 1
fi 