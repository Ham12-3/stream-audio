@echo off
echo ========================================
echo    Voice Gateway - Quick Start
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go from https://go.dev/dl/
    pause
    exit /b 1
)

REM Check if binary exists, build if not
if not exist "bin\gateway.exe" (
    echo Building gateway...
    go build -o bin\gateway.exe .\cmd\gateway
    if %ERRORLEVEL% NEQ 0 (
        echo ERROR: Build failed
        pause
        exit /b 1
    )
    echo Build complete!
    echo.
)

REM Check if Docker is available
where docker >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo Starting NATS JetStream...
    docker ps -a --filter "name=voice-gateway-nats" --format "{{.Names}}" | findstr voice-gateway-nats >nul
    if %ERRORLEVEL% NEQ 0 (
        docker run -d --name voice-gateway-nats -p 4222:4222 -p 8222:8222 nats:latest -js
        echo NATS started on ports 4222 and 8222
    ) else (
        docker start voice-gateway-nats >nul 2>nul
        echo NATS already running
    )
    echo.
) else (
    echo WARNING: Docker not found. NATS will not be started.
    echo Install Docker Desktop or run NATS manually.
    echo.
)

echo ========================================
echo Starting Voice Gateway...
echo.
echo Server will be available at:
echo   http://localhost:8080
echo.
echo To test:
echo   1. Open http://localhost:8080 in your browser
echo   2. Click 'Start Echo Test'
echo   3. Allow microphone access
echo   4. Speak and hear yourself back!
echo.
echo Press Ctrl+C to stop
echo ========================================
echo.

REM Run the gateway
bin\gateway.exe
