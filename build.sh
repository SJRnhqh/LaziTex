#!/bin/bash

echo ""
echo "============================================================"
echo ""
echo "     LaziTex Multi-Platform Build Script"
echo ""
echo "============================================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}[ERROR]${NC} $1 is not installed or not in PATH"
        echo ""
        case $1 in
            npm|node)
                echo "Please install Node.js:"
                echo "  macOS:   brew install node"
                echo "  Linux:   Use your distribution's package manager"
                echo "  Windows: Download from https://nodejs.org/"
                ;;
            go)
                echo "Please install Go:"
                echo "  macOS:   brew install go"
                echo "  Linux:   Use your distribution's package manager"
                echo "  Windows: Download from https://go.dev/dl/"
                ;;
        esac
        exit 1
    fi
}

# 检查 Go 版本
check_go_version() {
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    local major=$(echo $go_version | cut -d. -f1)
    local minor=$(echo $go_version | cut -d. -f2)
    
    # 检查是否满足 Go 1.25+ 要求
    if [ "$major" -lt 1 ] || ([ "$major" -eq 1 ] && [ "$minor" -lt 25 ]); then
        echo -e "${YELLOW}[WARNING]${NC} Go version may be too old"
        echo "  Required: Go 1.25+"
        echo "  Current:  $(go version)"
        echo ""
        read -p "Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# 检查 Go 代理配置（针对中国大陆用户）
check_go_proxy() {
    local proxy=$(go env GOPROXY 2>/dev/null)
    if [ -z "$proxy" ]; then
        proxy="direct"
    fi
    
    # 如果使用默认代理且可能无法访问，给出提示
    if [[ "$proxy" == *"proxy.golang.org"* ]] || [[ "$proxy" == "direct" ]]; then
        # 尝试快速检测代理是否可达（使用 curl 的 --max-time 参数，兼容性更好）
        if command -v curl &> /dev/null; then
            if ! curl -s --max-time 2 https://proxy.golang.org > /dev/null 2>&1; then
                echo -e "${YELLOW}[WARNING]${NC} Go proxy may be unreachable"
                echo "  Current GOPROXY: $proxy"
                echo "  If you're in China, consider setting:"
                echo "    go env -w GOPROXY=https://goproxy.cn,direct"
                echo ""
                read -p "Continue anyway? (y/N) " -n 1 -r
                echo
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    exit 1
                fi
            fi
        fi
    fi
}

# 创建输出目录
mkdir -p bin

total=6
current=0

# 步骤 1: 环境检查
current=$((current + 1))
echo "[$current/$total] Checking build environment..."
check_command node
check_command npm
check_command go

# 显示版本信息
echo "    Node.js: $(node --version)"
echo "    npm: $(npm --version)"
echo "    Go: $(go version | awk '{print $3}')"

check_go_version
check_go_proxy
echo -e "${GREEN}[OK]${NC} Build environment check passed"
echo ""

# 步骤 2: 构建 Vue 前端
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