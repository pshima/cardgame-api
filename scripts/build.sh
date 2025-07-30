#!/usr/bin/env bash

# Card Game API Cross-Platform Build Script
# Usage: ./scripts/build.sh [platform] [architecture]
# Example: ./scripts/build.sh linux amd64
# Example: ./scripts/build.sh all (builds all platforms)

set -eo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Configuration
readonly APP_NAME="cardgame-api"
readonly DIST_DIR="dist"
readonly BUILD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Build information
readonly VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
readonly BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
readonly COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
readonly LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitHash=${COMMIT_HASH} -w -s"
readonly BUILD_FLAGS="-trimpath -ldflags \"${LDFLAGS}\""

# Supported platforms (platform:suffix)
readonly SUPPORTED_PLATFORMS="
linux-amd64:
linux-arm64:
windows-amd64:.exe
windows-arm64:.exe
darwin-amd64:
darwin-arm64:
"

# Print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Print usage information
print_usage() {
    echo "Card Game API Cross-Platform Build Script"
    echo ""
    echo "Usage:"
    echo "  $0 [platform] [architecture]"
    echo "  $0 all"
    echo ""
    echo "Examples:"
    echo "  $0 linux amd64          # Build for Linux x86_64"
    echo "  $0 darwin arm64         # Build for macOS Apple Silicon"
    echo "  $0 windows amd64        # Build for Windows x86_64"
    echo "  $0 all                  # Build for all supported platforms"
    echo ""
    echo "Supported platforms:"
    for platform in "${!PLATFORMS[@]}"; do
        IFS=' ' read -r goos goarch suffix <<< "${PLATFORMS[$platform]}"
        echo "  - ${platform} (${goos}/${goarch})"
    done
    echo ""
    echo "Build information:"
    echo "  Version: ${VERSION}"
    echo "  Build time: ${BUILD_TIME}"
    echo "  Commit: ${COMMIT_HASH}"
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
    
    if ! GOOS="${goos}" GOARCH="${goarch}" CGO_ENABLED=0 \
        go build ${BUILD_FLAGS} -o "${output_file}" .; then
        print_status "${RED}" "❌ Failed to build for ${goos}/${goarch}"
        return 1
    fi
    
    # Get file size for display
    if command -v stat >/dev/null 2>&1; then
        if [[ "$(uname)" == "Darwin" ]]; then
            local size=$(stat -f%z "${output_file}")
        else
            local size=$(stat -c%s "${output_file}")
        fi
        local size_mb=$(awk "BEGIN {printf \"%.1f\", ${size}/1024/1024}")
        print_status "${GREEN}" "✅ Built ${output_file} (${size_mb} MB)"
    else
        print_status "${GREEN}" "✅ Built ${output_file}"
    fi
}

# Build all platforms
build_all() {
    print_status "${BLUE}" "Building for all supported platforms..."
    local failed_builds=()
    
    for platform in "${!PLATFORMS[@]}"; do
        IFS=' ' read -r goos goarch suffix <<< "${PLATFORMS[$platform]}"
        if ! build_platform "${goos}" "${goarch}" "${suffix}"; then
            failed_builds+=("${platform}")
        fi
    done
    
    if [[ ${#failed_builds[@]} -eq 0 ]]; then
        print_status "${GREEN}" "🎉 All builds completed successfully!"
        list_builds
    else
        print_status "${RED}" "❌ Some builds failed: ${failed_builds[*]}"
        return 1
    fi
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

# Main script logic
main() {
    cd "${BUILD_DIR}"
    
    case "${1:-}" in
        "")
            print_usage
            exit 0
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
            exit 0
            ;;
        *)
            if [[ $# -lt 2 ]]; then
                print_status "${RED}" "Error: Platform and architecture required"
                print_usage
                exit 1
            fi
            
            local goos=$1
            local goarch=$2
            local platform_key="${goos}-${goarch}"
            
            if [[ -z "${PLATFORMS[$platform_key]:-}" ]]; then
                print_status "${RED}" "Error: Unsupported platform ${platform_key}"
                print_usage
                exit 1
            fi
            
            create_dist_dir
            IFS=' ' read -r _ _ suffix <<< "${PLATFORMS[$platform_key]}"
            build_platform "${goos}" "${goarch}" "${suffix}"
            ;;
    esac
}

# Verify Go is installed
if ! command -v go >/dev/null 2>&1; then
    print_status "${RED}" "Error: Go is not installed or not in PATH"
    exit 1
fi

# Verify we're in a Go module
if [[ ! -f "go.mod" ]]; then
    print_status "${RED}" "Error: Not in a Go module directory"
    exit 1
fi

main "$@"