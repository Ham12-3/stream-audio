# Windows Setup Guide

## Quick Start (Windows)

### Option 1: Using the Batch File (Easiest)

1. **Open Command Prompt** (cmd) or PowerShell
2. **Navigate to the project**:
   ```cmd
   cd C:\Users\mobol\Downloads\stream-audio
   ```
3. **Run the batch file**:
   ```cmd
   start.bat
   ```
4. **Open your browser** to http://localhost:8080
5. **Click "Start Echo Test"** and speak!

### Option 2: Using WSL (What we used before)

Since your path was `/mnt/c/...` earlier, you have WSL (Windows Subsystem for Linux) installed:

1. **Open WSL/Ubuntu** (search for "Ubuntu" or "WSL" in Start menu)
2. **Navigate to project**:
   ```bash
   cd /mnt/c/Users/mobol/Downloads/stream-audio
   ```
3. **Run the Linux script**:
   ```bash
   ./start.sh
   ```
4. **Open browser** to http://localhost:8080

### Option 3: Manual Steps (Most Control)

#### Step 1: Set up Go Path
```cmd
set PATH=%USERPROFILE%\go\bin;%PATH%
set GOROOT=%USERPROFILE%\go
set GOPATH=%USERPROFILE%\go-workspace
```

#### Step 2: Build the Gateway
```cmd
cd C:\Users\mobol\Downloads\stream-audio
go build -o bin\gateway.exe .\cmd\gateway
```

#### Step 3: Start NATS (if you have Docker Desktop)
```cmd
docker run -d --name voice-nats -p 4222:4222 -p 8222:8222 nats:latest -js
```

Or download NATS for Windows from: https://nats.io/download/

#### Step 4: Run the Gateway
```cmd
bin\gateway.exe
```

#### Step 5: Test in Browser
Open http://localhost:8080

---

## Troubleshooting Windows Issues

### "go: command not found" or "'go' is not recognized"

**Problem**: Go is not in your PATH

**Solution**:
```cmd
set PATH=%USERPROFILE%\go\bin;%PATH%
set GOROOT=%USERPROFILE%\go
```

Or permanently add to Windows Environment Variables:
1. Press Win+R, type `sysdm.cpl`, press Enter
2. Click "Advanced" → "Environment Variables"
3. Add to PATH: `C:\Users\mobol\go\bin`

### "Docker is not running"

**Problem**: Docker Desktop not installed or not running

**Solutions**:
1. **Install Docker Desktop**: https://www.docker.com/products/docker-desktop/
2. **Or skip NATS** (echo will still work, but workers won't)
3. **Or download NATS** for Windows: https://nats.io/download/

### "Cannot connect to NATS"

**Problem**: NATS not running

**Solution**:
Check if NATS container is running:
```cmd
docker ps | findstr nats
```

If not, start it:
```cmd
docker start voice-nats
```

Or restart:
```cmd
docker restart voice-nats
```

### "Port 8080 is already in use"

**Problem**: Another application is using port 8080

**Solution**:
Find what's using it:
```cmd
netstat -ano | findstr :8080
```

Kill that process or change the port:
```cmd
set SERVER_PORT=8081
bin\gateway.exe
```

Then open http://localhost:8081

---

## Building on Windows

### Build All Components
```cmd
REM Build gateway
go build -o bin\gateway.exe .\cmd\gateway

REM Build ASR worker
go build -o bin\asr-worker.exe .\cmd\asr-worker

REM Build TTS worker
go build -o bin\tts-worker.exe .\cmd\tts-worker
```

### Clean Build
```cmd
rmdir /s /q bin
go build -o bin\gateway.exe .\cmd\gateway
```

---

## Running Without Docker (Windows Native)

If you don't want to use Docker, you can run NATS natively on Windows:

### Download NATS
1. Go to: https://github.com/nats-io/nats-server/releases
2. Download `nats-server-v*-windows-amd64.zip`
3. Extract to `C:\nats`

### Run NATS
```cmd
C:\nats\nats-server.exe -js
```

Leave this window open, open a new Command Prompt for the gateway.

---

## Using PowerShell Instead of CMD

If you prefer PowerShell:

```powershell
# Navigate to project
cd C:\Users\mobol\Downloads\stream-audio

# Set Go environment
$env:PATH = "$env:USERPROFILE\go\bin;$env:PATH"

# Build
go build -o bin\gateway.exe .\cmd\gateway

# Run
.\bin\gateway.exe
```

---

## Quick Test Commands (Windows)

```cmd
REM Check Go is installed
go version

REM Check Docker is running
docker ps

REM Check if NATS is accessible
curl http://localhost:8222/varz

REM Check if gateway is running
netstat -ano | findstr :8080

REM Build and run
go build -o bin\gateway.exe .\cmd\gateway && bin\gateway.exe
```

---

## Windows Firewall

If you get a Windows Firewall popup when starting the gateway:
1. **Click "Allow access"** on both Private and Public networks
2. This allows the web browser to connect to the gateway

---

## Recommended: Use WSL for Best Experience

Since you already have WSL (we used it during setup), it's the easiest way:

```bash
# Open WSL (Ubuntu)
wsl

# Navigate to project
cd /mnt/c/Users/mobol/Downloads/stream-audio

# Everything works like Linux now
./start.sh
```

---

## Testing Checklist (Windows)

- [ ] Go is installed: `go version` shows version
- [ ] Project built: `bin\gateway.exe` exists
- [ ] Docker running (optional): `docker ps` works
- [ ] NATS running: `curl http://localhost:8222` responds
- [ ] Gateway running: Terminal shows "Voice Gateway starting"
- [ ] Browser works: http://localhost:8080 loads
- [ ] Microphone allowed: Browser shows permission granted
- [ ] Echo works: You can hear yourself!

---

## Still Having Issues?

### Quick Debug
1. Check Go installation:
   ```cmd
   where go
   go version
   ```

2. Check if binary exists:
   ```cmd
   dir bin
   ```

3. Try running directly:
   ```cmd
   go run .\cmd\gateway\main.go
   ```

4. Check logs for errors

### Common Windows-Specific Issues

**Issue**: `\r\n` line endings in Go files
**Fix**:
```cmd
git config --global core.autocrlf input
```

**Issue**: Antivirus blocking binary
**Fix**: Add `bin\` folder to antivirus exclusions

**Issue**: Path with spaces
**Fix**: Use quotes:
```cmd
"C:\Users\My Name\Downloads\stream-audio\bin\gateway.exe"
```

---

## Alternative: Use Git Bash

If you have Git for Windows installed:

1. Open **Git Bash** (right-click in folder → "Git Bash Here")
2. Run:
   ```bash
   ./start.sh
   ```

Git Bash understands Linux commands!

---

## Summary

**Easiest way on Windows**:
```cmd
cd C:\Users\mobol\Downloads\stream-audio
start.bat
```

Then open http://localhost:8080 in your browser!
