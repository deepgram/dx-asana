#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get version from user
if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 v0.2.0"
  exit 1
fi

VERSION=$1

# Validate version format
if ! [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo -e "${RED}✗ Invalid version format. Use: v0.0.0${NC}"
  exit 1
fi

echo -e "${GREEN}📦 Preparing release for $VERSION${NC}"

# Build and test
echo -e "${YELLOW}🔨 Building...${NC}"
if ! make clean; then
  echo -e "${RED}✗ Build clean failed${NC}"
  exit 1
fi

if ! make build; then
  echo -e "${RED}✗ Build failed${NC}"
  exit 1
fi

echo -e "${YELLOW}✅ Running tests...${NC}"
if ! make test; then
  echo -e "${RED}✗ Tests failed${NC}"
  exit 1
fi

# Check if there are changes
if [ -z "$(git status --porcelain)" ]; then
  echo -e "${YELLOW}📝 No uncommitted changes${NC}"
else
  echo -e "${YELLOW}📝 Committing changes...${NC}"
  git add -A
  git commit -m "Release $VERSION" || true
fi

# Create tag
echo -e "${YELLOW}🏷️  Creating git tag...${NC}"
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo -e "${YELLOW}Tag $VERSION already exists. Deleting old tag...${NC}"
  git tag -d "$VERSION"
  git push origin --delete "$VERSION" || true
fi

git tag -a "$VERSION" -m "Release $VERSION"
echo -e "${GREEN}✓ Tag created: $VERSION${NC}"

# Push
echo -e "${YELLOW}📤 Pushing to GitHub...${NC}"
if ! git push origin main; then
  echo -e "${RED}✗ Push failed${NC}"
  exit 1
fi

if ! git push origin "$VERSION"; then
  echo -e "${RED}✗ Tag push failed${NC}"
  exit 1
fi

echo -e "${GREEN}✨ Release $VERSION initiated!${NC}"
echo ""
echo -e "${GREEN}GitHub Actions will now build and attach binaries to the release.${NC}"
echo ""
echo "Check progress at: https://github.com/deepgram/dx-asana/actions"
echo ""
echo -e "${YELLOW}Homebrew distribution is paused while the fork settles — install from source for now:${NC}"
echo "  git clone https://github.com/deepgram/dx-asana.git && cd dx-asana && make install"
