#!/bin/bash
# build-and-test.sh - Rebuild project and run all tests
# Increments patch version by the number of passing tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
VERSION_FILE="$PROJECT_ROOT/pkg/version/version.go"

cd "$PROJECT_ROOT"

# Colours for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}[INFO]${NC} Build started at $(date '+%Y-%m-%d %H:%M:%S')"
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

# Step 5: Run unit tests and count them
echo -e "${BLUE}[INFO]${NC} Running unit tests..."
TEST_OUTPUT=$(go test -v ./... 2>&1)
echo "$TEST_OUTPUT" | grep -E "^(ok|PASS|FAIL|\?)"

# Count passing tests
UNIT_TESTS=$(echo "$TEST_OUTPUT" | grep -c "^--- PASS" || echo "0")
echo -e "${GREEN}[PASS]${NC} Unit tests passed: $UNIT_TESTS"
echo ""

# Step 6: Compile all examples (integration tests)
echo -e "${BLUE}[INFO]${NC} Compiling example programs (integration tests)..."
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

echo -e "${GREEN}[PASS]${NC} Integration tests passed: $EXAMPLES_PASS"
echo ""

# Calculate total tests
TOTAL_TESTS=$((UNIT_TESTS + EXAMPLES_PASS))

# Step 7: Increment version by total number of tests
CURRENT_PATCH=$(grep -oP 'Patch = \K[0-9]+' "$VERSION_FILE")
NEW_PATCH=$((CURRENT_PATCH + TOTAL_TESTS))

# Format version string with leading zeros (e.g., 0.0.004)
NEW_VERSION_STR=$(printf "0.0.%03d" "$NEW_PATCH")

# Update version.go
sed -i "s/Patch = $CURRENT_PATCH/Patch = $NEW_PATCH/" "$VERSION_FILE"
sed -i "s/return \"0.0.[0-9]*\"/return \"$NEW_VERSION_STR\"/" "$VERSION_FILE"

# Rebuild with new version
echo -e "${BLUE}[INFO]${NC} Rebuilding with new version..."
go build -o hlc ./cmd/hlc

# Step 8: Show version
VERSION=$(./hlc --version 2>/dev/null | head -1)
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}[PASS]${NC} Build completed successfully"
echo ""
echo "  Unit tests:        $UNIT_TESTS"
echo "  Integration tests: $EXAMPLES_PASS"
echo "  Total tests:       $TOTAL_TESTS"
echo ""
echo "  Version:   0.0.$(printf "%03d" $CURRENT_PATCH) -> $NEW_VERSION_STR (+$TOTAL_TESTS)"
echo "  Binary:    $(ls -lh hlc | awk '{print $5}')"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
