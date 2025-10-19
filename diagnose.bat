@echo off
echo ==========================================
echo Voice Gateway Diagnostic Tool
echo ==========================================
echo.

echo [1/6] Checking Go in Windows...
%USERPROFILE%\go\bin\go version >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [OK] Go is installed in Windows
    %USERPROFILE%\go\bin\go version
) else (
    echo [FAIL] Go not found in Windows at %USERPROFILE%\go\bin
)
echo.

echo [2/6] Checking Go in PATH...
go version >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [OK] Go is in PATH
    go version
) else (
    echo [FAIL] Go not in PATH
    echo SOLUTION: Run this before building:
    echo   set PATH=%USERPROFILE%\go\bin;%%PATH%%
)
echo.

echo [3/6] Checking project location...
if exist "go.mod" (
    echo [OK] Found go.mod - you're in the right directory
) else (
    echo [FAIL] No go.mod found
    echo Are you in the stream-audio directory?
    echo Run: cd C:\Users\mobol\Downloads\stream-audio
)
echo.

echo [4/6] Checking if gateway is already built...
if exist "bin\gateway.exe" (
    echo [OK] Gateway binary exists
    dir bin\gateway.exe
) else (
    echo [INFO] Gateway not built yet - will need to build
)
echo.

echo [5/6] Checking port 8080...
netstat -ano | findstr :8080 >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [WARN] Something is using port 8080:
    netstat -ano | findstr :8080
    echo SOLUTION: Use different port or stop other service
) else (
    echo [OK] Port 8080 is available
)
echo.

echo [6/6] Checking WSL...
wsl echo WSL works >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [OK] WSL is available
    echo You can use WSL as backup option!
) else (
    echo [INFO] WSL not available
)
echo.

echo ==========================================
echo Summary and Recommendations
echo ==========================================
echo.

REM Check all conditions and give recommendation
%USERPROFILE%\go\bin\go version >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    if exist "go.mod" (
        echo RECOMMENDATION: Try running SIMPLE_START.bat
        echo.
        echo Run this command:
        echo   SIMPLE_START.bat
        echo.
    ) else (
        echo ERROR: You're not in the project directory!
        echo.
        echo Run this:
        echo   cd C:\Users\mobol\Downloads\stream-audio
        echo   diagnose.bat
    )
) else (
    wsl ~/go/bin/go version >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        echo RECOMMENDATION: Use WSL (Go is installed there)
        echo.
        echo Run these commands:
        echo   wsl
        echo   cd /mnt/c/Users/mobol/Downloads/stream-audio
        echo   ./start.sh
        echo.
    ) else (
        echo ERROR: Go is not installed!
        echo.
        echo Please install Go from:
        echo   https://go.dev/dl/
        echo.
    )
)

echo ==========================================
pause
