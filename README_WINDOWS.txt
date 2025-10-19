==========================================
Voice Gateway - Windows Quick Start
==========================================

STEP 1: Install Go for Windows
--------------------------------
Run this file:
  INSTALL_GO.bat

This will:
  1. Check if Go is installed
  2. Open browser to download page
  3. Guide you through installation

OR manually:
  1. Go to: https://go.dev/dl/
  2. Download: go1.21.5.windows-amd64.msi
  3. Run the installer
  4. Use default settings (Next, Next, Install)
  5. Close and open a NEW Command Prompt

To verify Go is installed:
  go version

You should see: "go version go1.21.x windows/amd64"


STEP 2: Build and Run
----------------------
Run this file:
  SIMPLE_START.bat

This will:
  1. Build the gateway (takes 30-60 seconds first time)
  2. Start the server
  3. Show you the URL to test


STEP 3: Test in Browser
------------------------
  1. Open browser
  2. Go to: http://localhost:8080
  3. Click "Start Echo Test"
  4. Allow microphone
  5. Speak and hear yourself!


==========================================
Troubleshooting
==========================================

Problem: "go: command not found" after installing
Solution: Close ALL Command Prompts, open NEW one

Problem: "Port 8080 in use"
Solution:
  set SERVER_PORT=8081
  bin\gateway.exe
  (Then open http://localhost:8081)

Problem: Build errors
Solution:
  go mod download
  go mod tidy
  go build -o bin\gateway.exe .\cmd\gateway


==========================================
Manual Commands (if batch files don't work)
==========================================

1. Navigate to project:
   cd C:\Users\mobol\Downloads\stream-audio

2. Build:
   go build -o bin\gateway.exe .\cmd\gateway

3. Run:
   bin\gateway.exe

4. Open browser to: http://localhost:8080


==========================================
Need Help?
==========================================

Read these files:
  - WINDOWS_INSTALL.md - Detailed setup guide
  - TROUBLESHOOT.md - Common problems
  - TESTING_GUIDE.md - How to test

Or run diagnostic:
  diagnose.bat
