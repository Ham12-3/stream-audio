@echo off
echo ==========================================
echo Go Installation Helper for Windows
echo ==========================================
echo.

echo Checking if Go is already installed...
go version >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo.
    echo [SUCCESS] Go is already installed!
    go version
    echo.
    echo You can now run: SIMPLE_START.bat
    pause
    exit /b 0
)

echo.
echo [INFO] Go is not installed on Windows
echo.
echo ==========================================
echo Please follow these steps:
echo ==========================================
echo.
echo 1. Your browser will open to Go downloads page
echo 2. Download: go1.21.5.windows-amd64.msi
echo 3. Run the installer (use default settings)
echo 4. Wait for installation to complete
echo 5. Come back here and press any key
echo.
echo Press any key to open browser...
pause >nul

REM Open Go download page
start https://go.dev/dl/

echo.
echo ==========================================
echo Installing Go...
echo ==========================================
echo.
echo After downloading:
echo   1. Double-click the .msi file
echo   2. Click "Next" through installation
echo   3. Wait for "Completed" message
echo   4. Close the installer
echo.
echo When finished, press any key to test...
pause >nul

echo.
echo Testing Go installation...
echo Please close this window and open a NEW Command Prompt
echo Then run: go version
echo.
echo If you see "go version go1.21.x", Go is installed!
echo Then you can run: SIMPLE_START.bat
echo.
pause
