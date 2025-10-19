@echo off
echo ========================================
echo Voice Gateway - Simple Test
echo ========================================
echo.

echo Checking Go installation...
go version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ERROR: Go is not installed on Windows!
    echo.
    echo Please install Go first by running:
    echo   INSTALL_GO.bat
    echo.
    echo Or manually download from: https://go.dev/dl/
    echo.
    pause
    exit /b 1
)

echo [OK] Go is installed:
go version
echo.

echo.
echo Building gateway (this may take 30 seconds)...
go build -o bin\gateway.exe .\cmd\gateway
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ERROR: Build failed!
    echo Check the error messages above
    pause
    exit /b 1
)

echo.
echo ========================================
echo Gateway built successfully!
echo.
echo NOTE: This will work WITHOUT Docker/NATS
echo The echo server will still function!
echo.
echo Starting gateway...
echo Open http://localhost:8080 in your browser
echo Press Ctrl+C to stop
echo ========================================
echo.

bin\gateway.exe
