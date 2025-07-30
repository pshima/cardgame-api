#!/bin/bash

# Docker Multi-Architecture Build Script for Card Game API
# Usage: ./scripts/docker-build.sh [options]

set -euo pipefail

# Colors for output  
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Configuration
readonly APP_NAME="cardgame-api"
readonly REGISTRY=${DOCKER_REGISTRY:-""}
readonly DOCKERFILE=${DOCKERFILE:-"Dockerfile"}

# Build information
readonly VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
readonly BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
readonly COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Default values
PUSH=false
PLATFORMS="linux/amd64,linux/arm64"
TAGS=""
BUILD_ARGS=""

# Print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Print usage information
print_usage() {
    cat << EOF
Docker Multi-Architecture Build Script for Card Game API

Usage: $0 [OPTIONS]

Options:
    -p, --push              Push images to registry after building
    --platforms PLATFORMS   Comma-separated list of platforms (default: linux/amd64,linux/arm64)
    -t, --tag TAG          Additional tag for the image (can be used multiple times)
    --registry REGISTRY    Docker registry to use (default: none)
    --build-arg ARG        Pass build argument to docker build (can be used multiple times)
    -h, --help             Show this help message

Examples:
    $0                                                    # Build for default platforms
    $0 --push                                            # Build and push to registry
    $0 --platforms linux/amd64                          # Build only for amd64
    $0 -t latest -t v1.0.0 --push                      # Build with multiple tags and push
    $0 --registry myregistry.com/myuser                 # Use custom registry
    $0 --build-arg VERSION=1.2.3                       # Pass custom build arguments

Supported Platforms:
    - linux/amd64 (x86_64)
    - linux/arm64 (aarch64)
    - linux/arm/v7 (ARMv7)

Build Information:
    Version: ${VERSION}
    Build Time: ${BUILD_TIME}
    Commit: ${COMMIT_HASH}
EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -p|--push)
                PUSH=true
                shift
                ;;
            --platforms)
                PLATFORMS="$2"
                shift 2
                ;;
            -t|--tag)
                if [[ -n "$TAGS" ]]; then
                    TAGS="$TAGS,$2"
                else
                    TAGS="$2"
                fi
                shift 2
                ;;
            --registry)
                REGISTRY="$2"
                shift 2
                ;;
            --build-arg)
                BUILD_ARGS="$BUILD_ARGS --build-arg $2"
                shift 2
                ;;
            -h|--help)
                print_usage
                exit 0
                ;;
            *)
                print_status "$RED" "Unknown option: $1"
                print_usage
                exit 1
                ;;
        esac
    done
}

# Generate image names with tags
generate_image_names() {
    local base_name
    if [[ -n "$REGISTRY" ]]; then
        base_name="${REGISTRY}/${APP_NAME}"
    else
        base_name="$APP_NAME"
    fi
    
    local image_names=""
    if [[ -n "$TAGS" ]]; then
        IFS=',' read -ra TAG_ARRAY <<< "$TAGS"
        for tag in "${TAG_ARRAY[@]}"; do
            if [[ -n "$image_names" ]]; then
                image_names="$image_names -t ${base_name}:${tag}"
            else
                image_names="-t ${base_name}:${tag}"
            fi
        done
    else
        # Default tag
        image_names="-t ${base_name}:${VERSION}"
    fi
    
    echo "$image_names"
}

# Check prerequisites
check_prerequisites() {
    # Check if Docker is installed
    if ! command -v docker >/dev/null 2>&1; then
        print_status "$RED" "Error: Docker is not installed or not in PATH"
        exit 1
    fi
    
    # Check if Docker buildx is available
    if ! docker buildx version >/dev/null 2>&1; then
        print_status "$RED" "Error: Docker Buildx is not available"
        print_status "$YELLOW" "Please install Docker Desktop or enable buildx plugin"
        exit 1
    fi
    
    # Check if we're in the right directory
    if [[ ! -f "$DOCKERFILE" ]]; then
        print_status "$RED" "Error: $DOCKERFILE not found in current directory"
        exit 1
    fi
    
    # Check if we're in a git repository for version info
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        print_status "$YELLOW" "Warning: Not in a git repository, version info may be limited"
    fi
}

# Set up Docker buildx
setup_buildx() {
    print_status "$BLUE" "Setting up Docker Buildx..."
    
    # Create a new buildx instance if it doesn't exist
    if ! docker buildx inspect multiarch >/dev/null 2>&1; then
        print_status "$YELLOW" "Creating new buildx instance 'multiarch'..."
        docker buildx create --name multiarch --driver docker-container --use
        docker buildx inspect --bootstrap
    else
        docker buildx use multiarch
    fi
}

# Build Docker images
build_images() {
    local image_names
    image_names=$(generate_image_names)
    
    print_status "$BLUE" "Building Docker images..."
    print_status "$YELLOW" "Platforms: $PLATFORMS"
    print_status "$YELLOW" "Tags: $image_names"
    print_status "$YELLOW" "Push: $PUSH"
    
    local push_flag=""
    if [[ "$PUSH" == "true" ]]; then
        push_flag="--push"
        print_status "$YELLOW" "Images will be pushed to registry"
    else
        push_flag="--load"
        print_status "$YELLOW" "Images will be loaded locally (single platform only)"
        # If not pushing, we can only load single platform
        if [[ "$PLATFORMS" == *","* ]]; then
            print_status "$YELLOW" "Multi-platform build detected, switching to --push mode"
            push_flag="--push"
            PUSH=true
        fi
    fi
    
    # Build command
    local build_cmd="docker buildx build"
    build_cmd="$build_cmd --platform $PLATFORMS"
    build_cmd="$build_cmd $image_names"
    build_cmd="$build_cmd --build-arg VERSION=$VERSION"
    build_cmd="$build_cmd --build-arg BUILD_TIME=$BUILD_TIME"
    build_cmd="$build_cmd --build-arg COMMIT_HASH=$COMMIT_HASH"
    build_cmd="$build_cmd $BUILD_ARGS"
    build_cmd="$build_cmd $push_flag"
    build_cmd="$build_cmd -f $DOCKERFILE"
    build_cmd="$build_cmd ."
    
    print_status "$BLUE" "Executing: $build_cmd"
    
    if eval "$build_cmd"; then
        print_status "$GREEN" "✅ Docker build completed successfully!"
        
        if [[ "$PUSH" == "true" ]]; then
            print_status "$GREEN" "Images pushed to registry"
        else
            print_status "$GREEN" "Images loaded locally"
            print_status "$BLUE" "Available images:"
            docker images "$APP_NAME" --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
        fi
    else
        print_status "$RED" "❌ Docker build failed"
        exit 1
    fi
}

# Main function
main() {
    parse_args "$@"
    check_prerequisites
    setup_buildx
    build_images
    
    print_status "$GREEN" "🎉 Multi-architecture build complete!"
    print_status "$BLUE" "Build info: $VERSION ($COMMIT_HASH) built at $BUILD_TIME"
}

main "$@"