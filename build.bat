@echo off
setlocal enabledelayedexpansion

echo.
echo ============================================================
echo.
echo     LaziTex Multi-Platform Build Script
echo.
echo ============================================================
echo.

REM 创建输出目录
if not exist "bin" mkdir bin

set "total=5"
set "current=0"

REM 步骤 1: 构建 Vue 前端
set /a current+=1
echo [%current%/%total%] Building Vue frontend...
cd frontend\vue
if not exist "node_modules" (
    echo    Installing npm dependencies...
    call npm install
    if %errorlevel% neq 0 (
        echo.
        echo [ERROR] npm install failed!
        pause
        exit /b 1
    )
)

echo    Building Vue frontend...
call npm run build
if %errorlevel% neq 0 (
    echo.
    echo [ERROR] Vue frontend build failed!
    pause
    exit /b 1
)
cd ..\..
echo    [OK] Vue frontend build completed successfully
echo.

echo Starting Go build process...
echo.

REM Windows
set /a current+=1
echo [%current%/%total%] Building Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -o bin/lazitex-windows-amd64.exe ./cmd/lazitex-cli
if %errorlevel% neq 0 (
    echo.
    echo [ERROR] Windows build failed!
    pause
    exit /b 1
)
echo    [OK] Windows build completed successfully
echo.

REM Linux
set /a current+=1
echo [%current%/%total%] Building Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -o bin/lazitex-linux-amd64 ./cmd/lazitex-cli
if %errorlevel% neq 0 (
    echo.
    echo [ERROR] Linux build failed!
    pause
    exit /b 1
)
echo    [OK] Linux build completed successfully
echo.

REM macOS Intel
set /a current+=1
echo [%current%/%total%] Building macOS Intel (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -o bin/lazitex-darwin-amd64 ./cmd/lazitex-cli
if %errorlevel% neq 0 (
    echo.
    echo [ERROR] macOS Intel build failed!
    pause
    exit /b 1
)
echo    [OK] macOS Intel build completed successfully
echo.

REM macOS Apple Silicon
set /a current+=1
echo [%current%/%total%] Building macOS Apple Silicon (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -o bin/lazitex-darwin-arm64 ./cmd/lazitex-cli
if %errorlevel% neq 0 (
    echo.
    echo [ERROR] macOS ARM64 build failed!
    pause
    exit /b 1
)
echo    [OK] macOS ARM64 build completed successfully
echo.

echo ============================================================
echo                    Build Complete!
echo ============================================================
echo.
echo All binaries are in the 'bin' directory:
echo.
echo    - lazitex-windows-amd64.exe
echo    - lazitex-linux-amd64
echo    - lazitex-darwin-amd64
echo    - lazitex-darwin-arm64
echo.
echo All platforms built successfully!
echo.
pause