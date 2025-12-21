#!/bin/bash

echo ""
echo "============================================================"
echo ""
echo "     LaziTex Multi-Platform Build Script"
echo ""
echo "============================================================"
echo ""

# 创建输出目录
mkdir -p bin

total=5
current=0

# 步骤 1: 构建 Vue 前端
current=$((current + 1))
echo "[$current/$total] Building Vue frontend..."
cd frontend/vue
if [ ! -d "node_modules" ]; then
    echo "    Installing npm dependencies..."
    npm install
    if [ $? -ne 0 ]; then
        echo ""
        echo "[ERROR] npm install failed!"
        exit 1
    fi
fi

echo "    Building Vue frontend..."
npm run build
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] Vue frontend build failed!"
    exit 1
fi
cd ../..
echo "    [OK] Vue frontend build completed successfully"
echo ""

echo "Starting Go build process..."
echo ""

# Windows
current=$((current + 1))
echo "[$current/$total] Building Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o bin/lazitex-windows-amd64.exe ./cmd/lazitex-cli
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] Windows build failed!"
    exit 1
fi
echo "    [OK] Windows build completed successfully"
echo ""

# Linux
current=$((current + 1))
echo "[$current/$total] Building Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o bin/lazitex-linux-amd64 ./cmd/lazitex-cli
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] Linux build failed!"
    exit 1
fi
echo "    [OK] Linux build completed successfully"
echo ""

# macOS Intel
current=$((current + 1))
echo "[$current/$total] Building macOS Intel (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o bin/lazitex-darwin-amd64 ./cmd/lazitex-cli
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] macOS Intel build failed!"
    exit 1
fi
echo "    [OK] macOS Intel build completed successfully"
echo ""

# macOS Apple Silicon
current=$((current + 1))
echo "[$current/$total] Building macOS Apple Silicon (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o bin/lazitex-darwin-arm64 ./cmd/lazitex-cli
if [ $? -ne 0 ]; then
    echo ""
    echo "[ERROR] macOS ARM64 build failed!"
    exit 1
fi
echo "    [OK] macOS ARM64 build completed successfully"
echo ""

echo "============================================================"
echo "                    Build Complete!"
echo "============================================================"
echo ""
echo "All binaries are in the 'bin' directory:"
echo ""
echo "    - lazitex-windows-amd64.exe"
echo "    - lazitex-linux-amd64"
echo "    - lazitex-darwin-amd64"
echo "    - lazitex-darwin-arm64"
echo ""
echo "All platforms built successfully!"
echo ""