#!/bin/bash
#
# Release script for webex-mcp-server
# This script implements the release stage of the 12-factor app methodology
# by combining the build artifact with configuration
#

set -e

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
VERSION=${1:-latest}
BUILD_DIR="releases/${VERSION}"
DOCKER_REGISTRY=${DOCKER_REGISTRY:-""}
DOCKER_IMAGE_NAME="webex-mcp-server"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Validate environment
validate_environment() {
    log_info "Validating environment..."
    
    if ! command -v git &> /dev/null; then
        log_error "git is not installed"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        log_error "docker is not installed"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        log_error "go is not installed"
        exit 1
    fi
}

# Build application
build_application() {
    log_info "Building application..."
    make clean
    make build
    
    if [ ! -f "build/webex-mcp-server" ]; then
        log_error "Build failed: binary not found"
        exit 1
    fi
}

# Create release directory
create_release() {
    log_info "Creating release ${VERSION}..."
    
    # Create release directory
    mkdir -p "${BUILD_DIR}"
    
    # Copy build artifact
    cp build/webex-mcp-server "${BUILD_DIR}/"
    
    # Copy configuration template
    cp .env.example "${BUILD_DIR}/.env.template"
    
    # Create release metadata
    cat > "${BUILD_DIR}/release.json" <<EOF
{
  "version": "${VERSION}",
  "build_date": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "git_commit": "$(git rev-parse HEAD)",
  "git_branch": "$(git branch --show-current)",
  "git_tag": "$(git describe --tags --always)",
  "go_version": "$(go version | awk '{print $3}')",
  "build_user": "${USER}",
  "build_host": "$(hostname)"
}
EOF
    
    # Create README for the release
    cat > "${BUILD_DIR}/README.md" <<EOF
# Webex MCP Server Release ${VERSION}

## Release Information
- Version: ${VERSION}
- Build Date: $(date -u +%Y-%m-%dT%H:%M:%SZ)
- Git Commit: $(git rev-parse HEAD)

## Configuration
1. Copy \`.env.template\` to \`.env\`
2. Fill in your Webex API credentials
3. Run the server:
   \`\`\`bash
   ./webex-mcp-server -http :3001
   \`\`\`

## Environment Variables
- WEBEX_PUBLIC_WORKSPACE_API_KEY: Your Webex API token
- WEBEX_API_BASE_URL: Webex API base URL (default: https://webexapis.com/v1)
- PORT: Server port (default: 3001)
- USE_FASTHTTP: Enable fast HTTP client (default: true)
EOF
    
    log_info "Release directory created at ${BUILD_DIR}"
}

# Build Docker images
build_docker_images() {
    log_info "Building Docker images..."
    
    # Copy .dockerignore temporarily for build
    cp deployment/docker/.dockerignore .dockerignore 2>/dev/null || true
    
    # Build main image
    docker build -f deployment/docker/Dockerfile -t ${DOCKER_IMAGE_NAME}:${VERSION} .
    docker tag ${DOCKER_IMAGE_NAME}:${VERSION} ${DOCKER_IMAGE_NAME}:latest
    
    # Clean up temporary .dockerignore
    rm -f .dockerignore
    
    # Tag for registry if specified
    if [ -n "${DOCKER_REGISTRY}" ]; then
        docker tag ${DOCKER_IMAGE_NAME}:${VERSION} ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${VERSION}
        docker tag ${DOCKER_IMAGE_NAME}:latest ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:latest
        log_info "Tagged images for registry: ${DOCKER_REGISTRY}"
    fi
    
    # Save image to release directory
    log_info "Saving Docker image to release directory..."
    docker save ${DOCKER_IMAGE_NAME}:${VERSION} | gzip > "${BUILD_DIR}/webex-mcp-server-${VERSION}.tar.gz"
}

# Create release archive
create_archive() {
    log_info "Creating release archive..."
    
    # Create tarball excluding Docker image
    tar -czf "releases/webex-mcp-server-${VERSION}.tar.gz" \
        -C "${BUILD_DIR}" \
        --exclude="*.tar.gz" \
        .
    
    log_info "Release archive created: releases/webex-mcp-server-${VERSION}.tar.gz"
}

# Generate checksums
generate_checksums() {
    log_info "Generating checksums..."
    
    cd "${BUILD_DIR}"
    
    # Generate SHA256 checksums
    if command -v sha256sum &> /dev/null; then
        sha256sum webex-mcp-server > SHA256SUMS
        sha256sum *.tar.gz >> SHA256SUMS 2>/dev/null || true
    elif command -v shasum &> /dev/null; then
        shasum -a 256 webex-mcp-server > SHA256SUMS
        shasum -a 256 *.tar.gz >> SHA256SUMS 2>/dev/null || true
    fi
    
    cd - > /dev/null
}

# Main execution
main() {
    log_info "Starting release process for version ${VERSION}"
    
    validate_environment
    build_application
    create_release
    build_docker_images
    create_archive
    generate_checksums
    
    log_info "Release ${VERSION} completed successfully!"
    log_info "Release artifacts:"
    echo "  - Binary: ${BUILD_DIR}/webex-mcp-server"
    echo "  - Docker image: ${DOCKER_IMAGE_NAME}:${VERSION}"
    echo "  - Archive: releases/webex-mcp-server-${VERSION}.tar.gz"
    echo "  - Docker export: ${BUILD_DIR}/webex-mcp-server-${VERSION}.tar.gz"
    
    if [ -n "${DOCKER_REGISTRY}" ]; then
        log_warn "Don't forget to push Docker images:"
        echo "  docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${VERSION}"
        echo "  docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:latest"
    fi
}

# Run main function
main