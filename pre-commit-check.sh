#!/bin/bash

# Pre-commit check script
# Run this before every commit to ensure code quality

set -e

echo "🔍 Running pre-commit checks..."

# 1. Format check
echo ""
echo "1️⃣  Checking code formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "❌ Code is not formatted. Running go fmt..."
    go fmt ./...
    echo "✅ Code formatted"
else
    echo "✅ Code is properly formatted"
fi

# 2. Go vet
echo ""
echo "2️⃣  Running go vet..."
if go vet ./...; then
    echo "✅ Go vet passed"
else
    echo "❌ Go vet failed"
    exit 1
fi

# 3. Linting
echo ""
echo "3️⃣  Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run; then
        echo "✅ Linting passed"
    else
        echo "❌ Linting failed"
        exit 1
    fi
else
    echo "⚠️  golangci-lint not found. Install with:"
    echo "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    echo "   Skipping lint check..."
fi

# 4. Tests
echo ""
echo "4️⃣  Running tests..."
if go test -v ./...; then
    echo "✅ All tests passed"
else
    echo "❌ Tests failed"
    exit 1
fi

# 5. Coverage check
echo ""
echo "5️⃣  Checking test coverage..."
go test -coverprofile=coverage.out ./... > /dev/null 2>&1
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "   Test coverage: ${coverage}%"

if (( $(echo "$coverage < 70.0" | bc -l) )); then
    echo "⚠️  Warning: Test coverage is below 70%"
else
    echo "✅ Test coverage is good"
fi

# 6. Build check
echo ""
echo "6️⃣  Checking build..."
if go build -o markdown-viewer-editor .; then
    echo "✅ Build successful"
    rm -f markdown-viewer-editor
else
    echo "❌ Build failed"
    exit 1
fi

echo ""
echo "🎉 All pre-commit checks passed!"
echo "   You can now commit your changes."
