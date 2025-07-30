#!/bin/bash
# Generate secure secrets for the application

set -euo pipefail

# Colors
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

SECRETS_DIR="./secrets"

echo -e "${GREEN}Generating secure secrets...${NC}"

# Create secrets directory
mkdir -p "$SECRETS_DIR"

# Generate API key (32 bytes, base64 encoded)
if [ ! -f "$SECRETS_DIR/api_key.txt" ]; then
    openssl rand -base64 32 > "$SECRETS_DIR/api_key.txt"
    echo -e "${GREEN}✓ Generated API key${NC}"
else
    echo -e "${YELLOW}⚠ API key already exists${NC}"
fi

# Generate database password (24 characters)
if [ ! -f "$SECRETS_DIR/db_password.txt" ]; then
    openssl rand -base64 24 | tr -d '\n' > "$SECRETS_DIR/db_password.txt"
    echo -e "${GREEN}✓ Generated database password${NC}"
else
    echo -e "${YELLOW}⚠ Database password already exists${NC}"
fi

# Generate JWT secret (64 bytes, hex encoded)
if [ ! -f "$SECRETS_DIR/jwt_secret.txt" ]; then
    openssl rand -hex 64 > "$SECRETS_DIR/jwt_secret.txt"
    echo -e "${GREEN}✓ Generated JWT secret${NC}"
else
    echo -e "${YELLOW}⚠ JWT secret already exists${NC}"
fi

# Set proper permissions (read-only for owner)
chmod 600 "$SECRETS_DIR"/*.txt
chmod 700 "$SECRETS_DIR"

echo -e "\n${GREEN}Secrets generated successfully!${NC}"
echo -e "Location: $SECRETS_DIR/"
echo -e "\n${YELLOW}⚠ WARNING: Never commit these files to version control!${NC}"
echo -e "Add the following to .gitignore:"
echo -e "  secrets/*.txt"