#!/bin/bash
# Security scanning script for Card Game API container

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

IMAGE_NAME="${1:-cardgame-api:latest}"

echo "🔒 Security Scan for ${IMAGE_NAME}"
echo "=================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Build the image if it doesn't exist
if ! docker image inspect "${IMAGE_NAME}" > /dev/null 2>&1; then
    echo -e "${YELLOW}Building image ${IMAGE_NAME}...${NC}"
    docker build -t "${IMAGE_NAME}" .
fi

# 1. Scan with Trivy for vulnerabilities
echo -e "\n${GREEN}1. Vulnerability Scan (Trivy)${NC}"
if command -v trivy &> /dev/null; then
    trivy image --severity HIGH,CRITICAL "${IMAGE_NAME}"
else
    docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
        aquasec/trivy:latest image --severity HIGH,CRITICAL "${IMAGE_NAME}"
fi

# 2. Check for secrets with SecretScanner
echo -e "\n${GREEN}2. Secret Scan${NC}"
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
    -v "$(pwd)":/src \
    trufflesecurity/trufflehog:latest \
    filesystem /src --no-update

# 3. Static analysis with hadolint for Dockerfile
echo -e "\n${GREEN}3. Dockerfile Lint (Hadolint)${NC}"
if command -v hadolint &> /dev/null; then
    hadolint Dockerfile
else
    docker run --rm -i hadolint/hadolint < Dockerfile
fi

# 4. Check image configuration
echo -e "\n${GREEN}4. Image Configuration Check${NC}"
echo "Checking for security best practices..."

# Check if running as non-root
USER=$(docker inspect "${IMAGE_NAME}" --format='{{.Config.User}}')
if [ -z "$USER" ] || [ "$USER" = "root" ] || [ "$USER" = "0" ]; then
    echo -e "${RED}❌ Image runs as root user${NC}"
else
    echo -e "${GREEN}✓ Image runs as non-root user: $USER${NC}"
fi

# Check for HEALTHCHECK
HEALTHCHECK=$(docker inspect "${IMAGE_NAME}" --format='{{.Config.Healthcheck}}')
if [ "$HEALTHCHECK" = "<nil>" ]; then
    echo -e "${RED}❌ No HEALTHCHECK defined${NC}"
else
    echo -e "${GREEN}✓ HEALTHCHECK is defined${NC}"
fi

# Check exposed ports
PORTS=$(docker inspect "${IMAGE_NAME}" --format='{{range $p, $conf := .Config.ExposedPorts}}{{$p}} {{end}}')
echo -e "${GREEN}✓ Exposed ports: $PORTS${NC}"

# 5. Check for outdated dependencies
echo -e "\n${GREEN}5. Dependency Check${NC}"
docker run --rm -v "$(pwd)":/src -w /src golang:1.24.4-alpine sh -c "
    go list -m -u all 2>/dev/null | grep -E '\[.*\]' || echo 'All dependencies are up to date'
"

# 6. SBOM Generation
echo -e "\n${GREEN}6. Software Bill of Materials (SBOM)${NC}"
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
    anchore/syft:latest "${IMAGE_NAME}" -o cyclonedx-json > sbom.json
echo -e "${GREEN}✓ SBOM generated: sbom.json${NC}"

# Summary
echo -e "\n${GREEN}Security Scan Complete!${NC}"
echo "========================"
echo "Review the output above for any security issues."
echo "Fix any HIGH or CRITICAL vulnerabilities before deploying to production."