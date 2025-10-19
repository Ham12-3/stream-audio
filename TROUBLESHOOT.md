# Troubleshooting Guide

## Step-by-Step Diagnosis

### Step 1: Check if Go is Working

Open Command Prompt and run:
```cmd
%USERPROFILE%\go\bin\go version
```

**Expected output:**
```
go version go1.25.3 linux/amd64
```

**If you get an error:**
- Go is not installed correctly
- Run the commands below to reinstall

---

### Step 2: Verify Go Installation Path

```cmd
dir %USERPROFILE%\go\bin
```

**Expected output:**
```
go.exe
gofmt.exe
```

**If folder doesn't exist:**
Go is not installed. See "Installing Go" section below.

---

### Step 3: Try Building Manually

```cmd
cd C:\Users\mobol\Downloads\stream-audio
set PATH=%USERPROFILE%\go\bin;%PATH%
go build -o test.exe .\cmd\gateway\main.go
```

**Copy the EXACT error message** and let me know what it says!

---

## Common Errors & Solutions

### Error: "go: command not found" or "'go' is not recognized"

**Problem:** Go is not in your PATH

**Solution:**
```cmd
set PATH=%USERPROFILE%\go\bin;%PATH%
set GOROOT=%USERPROFILE%\go
set GOPATH=%USERPROFILE%\go-workspace

REM Test again
go version
```

If still fails, Go is not installed correctly.

---

### Error: "package XXX is not in GOROOT"

**Problem:** Dependencies not downloaded

**Solution:**
```cmd
cd C:\Users\mobol\Downloads\stream-audio
go mod download
go mod tidy
go build -o bin\gateway.exe .\cmd\gateway
```

---

### Error: "cannot find module providing package"

**Problem:** go.mod or dependencies corrupted

**Solution:**
```cmd
cd C:\Users\mobol\Downloads\stream-audio
go clean -modcache
go mod download
go build -o bin\gateway.exe .\cmd\gateway
```

---

### Error: "listen tcp :8080: bind: address already in use"

**Problem:** Something else is using port 8080

**Find what's using it:**
```cmd
netstat -ano | findstr :8080
```

**Solution 1 - Kill the process:**
Note the PID number (last column), then:
```cmd
taskkill /PID <number> /F
```

**Solution 2 - Use different port:**
```cmd
set SERVER_PORT=8081
bin\gateway.exe
```
Then open http://localhost:8081

---

### Error: "Failed to connect to NATS"

**This is OK!** The echo server works without NATS.

**If you want NATS:**
```cmd
docker run -d -p 4222:4222 nats:latest -js
```

Or ignore it - echo still works!

---

## Nuclear Option: Complete Fresh Start

If nothing works, let's start completely fresh:

### 1. Check Go Installation
```cmd
dir %USERPROFILE%\go
```

If this folder doesn't exist, Go isn't installed.

### 2. Navigate to Project
```cmd
cd C:\Users\mobol\Downloads\stream-audio
```

### 3. Set Environment
```cmd
set PATH=%USERPROFILE%\go\bin;%PATH%
set GOROOT=%USERPROFILE%\go
set GOPATH=%USERPROFILE%\go-workspace
```

### 4. Clean Everything
```cmd
rmdir /s /q bin
go clean -cache
go clean -modcache
```

### 5. Download Dependencies
```cmd
go mod download
```

### 6. Build
```cmd
go build -o bin\gateway.exe .\cmd\gateway
```

### 7. Run
```cmd
bin\gateway.exe
```

---

## Alternative: Use WSL (Easiest!)

Since we installed Go in WSL earlier:

### 1. Open WSL
Search for "Ubuntu" in Windows Start menu, or:
```cmd
wsl
```

### 2. Navigate to Project
```bash
cd /mnt/c/Users/mobol/Downloads/stream-audio
```

### 3. Run
```bash
./start.sh
```

**This should just work!**

---

## Testing Without Building

Try running directly:
```cmd
cd C:\Users\mobol\Downloads\stream-audio
set PATH=%USERPROFILE%\go\bin;%PATH%
go run .\cmd\gateway\main.go
```

This runs without building. Slower but works if build fails.

---

## Collecting Debug Info

If still failing, run these commands and **send me the output**:

```cmd
echo === System Info ===
systeminfo | findstr /C:"OS Name" /C:"OS Version"

echo === Go Version ===
%USERPROFILE%\go\bin\go version

echo === Go Location ===
where go

echo === Project Files ===
dir C:\Users\mobol\Downloads\stream-audio

echo === Go Mod ===
type go.mod

echo === Try Build ===
cd C:\Users\mobol\Downloads\stream-audio
set PATH=%USERPROFILE%\go\bin;%PATH%
go build -o test.exe .\cmd\gateway 2>&1
```

Copy all the output and send it to me!

---

## Absolute Simplest Test

Let's just test if Go works at all:

### Create test.go
Create a file `C:\Users\mobol\Downloads\test.go`:
```go
package main
import "fmt"
func main() {
    fmt.Println("Go works!")
}
```

### Run it
```cmd
cd C:\Users\mobol\Downloads
%USERPROFILE%\go\bin\go run test.go
```

**Expected:** Prints "Go works!"

If this works, Go is fine. If not, Go is broken.

---

## If Go Is Broken: Reinstall

### 1. Check if Go is in WSL or Windows

Go might be installed in WSL (Linux) but not Windows.

**Check Windows:**
```cmd
dir %USERPROFILE%\go
```

**Check WSL:**
```cmd
wsl ls -l /home/*/go
```

### 2. If Go is only in WSL

**Option A: Use WSL** (Recommended)
```cmd
wsl
cd /mnt/c/Users/mobol/Downloads/stream-audio
./start.sh
```

**Option B: Install Go for Windows**
1. Download: https://go.dev/dl/go1.21.5.windows-amd64.msi
2. Run installer
3. Open new Command Prompt
4. Test: `go version`

---

## Quick Check Script

Save this as `check.bat` and run it:

```batch
@echo off
echo Checking Go in Windows...
%USERPROFILE%\go\bin\go version
if %ERRORLEVEL% EQU 0 (
    echo SUCCESS: Go works in Windows!
) else (
    echo FAILED: Go not in Windows
    echo Checking WSL...
    wsl ~/go/bin/go version
    if %ERRORLEVEL% EQU 0 (
        echo SUCCESS: Go works in WSL!
        echo SOLUTION: Use WSL to run the project
        echo Run: wsl
        echo Then: cd /mnt/c/Users/mobol/Downloads/stream-audio
        echo Then: ./start.sh
    ) else (
        echo FAILED: Go not found anywhere
        echo You need to install Go
    )
)
pause
```

---

## What Error Did You Get?

Please tell me:

1. **What command did you run?**
   - `start.bat`?
   - `./start.sh`?
   - Something else?

2. **What was the error message?**
   - Copy the full error text

3. **Which terminal?**
   - Command Prompt (cmd)?
   - PowerShell?
   - WSL/Ubuntu?
   - Git Bash?

With this info, I can give you the exact fix!

---

## Meanwhile: Try WSL (Should Work Immediately)

We installed Go in WSL during setup, so this should work:

```cmd
wsl
cd /mnt/c/Users/mobol/Downloads/stream-audio
./start.sh
```

Then open http://localhost:8080 in Windows browser!
