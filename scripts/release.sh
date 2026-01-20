#!/bin/bash
# Release script for ia-start CLI
#
# This script automates the release process:
# 1. Runs all tests
# 2. Builds for all platforms
# 3. Creates release archives
# 4. Generates checksums
# 5. Creates a git tag
# 6. Pushes the tag to trigger GitHub release

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if we're in a git repo
if ! git rev-parse --git-head > /dev/null 2>&1; then
    log_error "Not in a git repository"
    exit 1
fi

# Check if working tree is clean
if [ -n "$(git status --porcelain)" ]; then
    log_error "Working tree is not clean. Please commit or stash changes first."
    git status --short
    exit 1
fi

# Get version argument
if [ -z "$1" ]; then
    log_error "Usage: $0 <version> [remote]"
    echo ""
    echo "Example:"
    echo "  $0 v0.1.0       # Create and push v0.1.0 tag"
    echo "  $0 v0.1.0 origin # Create and push to origin remote"
    exit 1
fi

VERSION="$1"
REMOTE="${2:-origin}"

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
    log_error "Invalid version format: $VERSION"
    log_error "Version must match: v<major>.<minor>.<patch>[-<prerelease>]"
    echo ""
    echo "Examples:"
    echo "  v1.0.0"
    echo "  v1.2.3-beta.1"
    echo "  v2.0.0-rc.1"
    exit 1
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    log_error "Tag $VERSION already exists"
    log_info "To delete existing tag: git tag -d $VERSION && git push $REMOTE :refs/tags/$VERSION"
    exit 1
fi

log_info "Starting release process for $VERSION"
echo ""

# Step 1: Run tests
log_info "Step 1/6: Running tests..."
if ! make test-race; then
    log_error "Tests failed"
    exit 1
fi
log_success "Tests passed"
echo ""

# Step 2: Run linting
log_info "Step 2/6: Running linters..."
if ! make lint; then
    log_error "Linting failed"
    exit 1
fi
log_success "Linting passed"
echo ""

# Step 3: Build for all platforms
log_info "Step 3/6: Building for all platforms..."
if ! make build-all VERSION="$VERSION"; then
    log_error "Build failed"
    exit 1
fi
log_success "Build completed"
echo ""

# Step 4: Create release archives
log_info "Step 4/6: Creating release archives..."
if ! make release VERSION="$VERSION"; then
    log_error "Release creation failed"
    exit 1
fi
log_success "Release archives created"
echo ""

# Step 5: Create git tag
log_info "Step 5/6: Creating git tag..."
if ! git tag -a "$VERSION" -m "Release $VERSION"; then
    log_error "Failed to create tag"
    exit 1
fi
log_success "Tag $VERSION created"
echo ""

# Step 6: Push tag
log_info "Step 6/6: Pushing tag to $REMOTE..."
read -p "Push tag to $REMOTE? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if ! git push "$REMOTE" "$VERSION"; then
        log_error "Failed to push tag"
        exit 1
    fi
    log_success "Tag pushed successfully"
    echo ""
    log_info "GitHub Actions will now:"
    log_info "  1. Run all tests"
    log_info "  2. Build for all platforms"
    log_info "  3. Create GitHub release"
    log_info "  4. Upload release artifacts"
    echo ""
    log_info "Watch the release at:"
    log_info "  https://github.com/drossan/claude-init/actions"
else
    log_warning "Tag not pushed. To push manually:"
    log_warning "  git push $REMOTE $VERSION"
fi

echo ""
log_success "Release $VERSION ready!"
echo ""
