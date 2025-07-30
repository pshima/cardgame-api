#!/usr/bin/env bash

# Simple Cross-Platform Build Script for Card Game API
# Usage: ./scripts/build-simple.sh [platform] [architecture]
# Example: ./scripts/build-simple.sh linux amd64
# Example: ./scripts/build-simple.sh all

set -e

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Configuration
readonly APP_NAME="cardgame-api"
readonly DIST_DIR="dist"

# Build information
readonly VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
readonly BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
readonly COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
readonly LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitHash=${COMMIT_HASH} -w -s"

# Print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Create distribution directory
create_dist_dir() {
    if [[ ! -d "${DIST_DIR}" ]]; then
        mkdir -p "${DIST_DIR}"
        print_status "${BLUE}" "Created distribution directory: ${DIST_DIR}"
    fi
}

# Build for a specific platform
build_platform() {
    local goos=$1
    local goarch=$2
    local suffix=${3:-""}
    
    local output_file="${DIST_DIR}/${APP_NAME}-${goos}-${goarch}${suffix}"
    
    print_status "${YELLOW}" "Building for ${goos}/${goarch}..."
    
    if GOOS="${goos}" GOARCH="${goarch}" CGO_ENABLED=0 \
        go build -trimpath -ldflags "${LDFLAGS}" -o "${output_file}" .; then
        print_status "${GREEN}" "✅ Built ${output_file}"
    else
        print_status "${RED}" "❌ Failed to build for ${goos}/${goarch}"
        return 1
    fi
}

# Build all platforms
build_all() {
    print_status "${BLUE}" "Building for all supported platforms..."
    
    build_platform "linux" "amd64" ""
    build_platform "linux" "arm64" ""
    build_platform "windows" "amd64" ".exe" 
    build_platform "windows" "arm64" ".exe"
    build_platform "darwin" "amd64" ""
    build_platform "darwin" "arm64" ""
    
    print_status "${GREEN}" "🎉 All builds completed!"
    list_builds
}

# List built binaries
list_builds() {
    if [[ -d "${DIST_DIR}" ]] && [[ -n "$(ls -A "${DIST_DIR}" 2>/dev/null)" ]]; then
        print_status "${BLUE}" "Built binaries:"
        ls -lah "${DIST_DIR}/"
    else
        print_status "${YELLOW}" "No binaries found in ${DIST_DIR}/"
    fi
}

# Clean build artifacts
clean_builds() {
    if [[ -d "${DIST_DIR}" ]]; then
        rm -rf "${DIST_DIR}"
        print_status "${GREEN}" "Cleaned build artifacts"
    else
        print_status "${YELLOW}" "No build artifacts to clean"
    fi
}

# Print usage
print_usage() {
    echo "Simple Cross-Platform Build Script"
    echo ""
    echo "Usage: $0 [platform] [architecture]"
    echo "       $0 all"
    echo "       $0 clean"
    echo "       $0 list"
    echo ""
    echo "Examples:"
    echo "  $0 linux amd64"
    echo "  $0 darwin arm64"
    echo "  $0 windows amd64"
    echo "  $0 all"
    echo ""
    echo "Build info: ${VERSION} (${COMMIT_HASH}) at ${BUILD_TIME}"
}

# Verify Go is installed
if ! command -v go >/dev/null 2>&1; then
    print_status "${RED}" "Error: Go is not installed or not in PATH"
    exit 1
fi

# Main script logic
case "${1:-}" in
    "")
        print_usage
        ;;
    "all")
        create_dist_dir
        build_all
        ;;
    "clean")
        clean_builds
        ;;
    "list")
        list_builds
        ;;
    "help"|"-h"|"--help")
        print_usage
        ;;
    *)
        if [[ $# -lt 2 ]]; then
            print_status "${RED}" "Error: Platform and architecture required"
            print_usage
            exit 1
        fi
        
        goos=$1
        goarch=$2
        suffix=""
        
        # Set suffix for Windows
        if [[ "$goos" == "windows" ]]; then
            suffix=".exe"
        fi
        
        create_dist_dir
        build_platform "${goos}" "${goarch}" "${suffix}"
        ;;
esac