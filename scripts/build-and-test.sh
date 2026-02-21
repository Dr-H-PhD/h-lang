#!/bin/bash
# build-and-test.sh - Rebuild project and run all tests
# Increments build number on each execution

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BUILD_FILE="$PROJECT_ROOT/.build-number"

cd "$PROJECT_ROOT"

# Colours for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Read current build number or initialise
if [ -f "$BUILD_FILE" ]; then
    BUILD_NUM=$(cat "$BUILD_FILE")
else
    BUILD_NUM=0
fi

# Increment build number
BUILD_NUM=$((BUILD_NUM + 1))
echo "$BUILD_NUM" > "$BUILD_FILE"

echo -e "${BLUE}[INFO]${NC} Build #$BUILD_NUM started at $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Step 1: Clean
echo -e "${BLUE}[INFO]${NC} Cleaning previous build artefacts..."
rm -f hlc 2>/dev/null || true

# Step 2: Format code
echo -e "${BLUE}[INFO]${NC} Formatting code..."
go fmt ./... > /dev/null 2>&1

# Step 3: Vet code
echo -e "${BLUE}[INFO]${NC} Running go vet..."
go vet ./... 2>&1
echo -e "${GREEN}[PASS]${NC} Code passed vet checks"
echo ""

# Step 4: Build
echo -e "${BLUE}[INFO]${NC} Building compiler..."
go build -o hlc ./cmd/hlc
echo -e "${GREEN}[PASS]${NC} Compiler built successfully"
echo ""

# Step 5: Run unit tests
echo -e "${BLUE}[INFO]${NC} Running unit tests..."
go test ./... 2>&1
echo -e "${GREEN}[PASS]${NC} All unit tests passed"
echo ""

# Step 6: Compile all examples
echo -e "${BLUE}[INFO]${NC} Compiling example programs..."
EXAMPLES_PASS=0

for f in examples/*.hl; do
    name=$(basename "$f" .hl)
    if ./hlc "$f" > /dev/null 2>&1 && [ -f "./$name" ]; then
        rm -f "./$name"
        EXAMPLES_PASS=$((EXAMPLES_PASS + 1))
        echo -e "  ${GREEN}✓${NC} $name"
    else
        echo -e "  ${RED}✗${NC} $name"
        exit 1
    fi
done

echo -e "${GREEN}[PASS]${NC} All examples compiled ($EXAMPLES_PASS examples)"
echo ""

# Step 7: Show version
VERSION=$(./hlc --version 2>/dev/null | head -1)
echo -e "${BLUE}[INFO]${NC} Compiler version: $VERSION"
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}[PASS]${NC} Build #$BUILD_NUM completed successfully"
echo "  Examples:  $EXAMPLES_PASS compiled"
echo "  Binary:    $(ls -lh hlc | awk '{print $5}')"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
