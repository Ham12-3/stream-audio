# Windows Installation Guide (No WSL)

## Step 1: Install Go for Windows

### Download Go
1. Open your browser
2. Go to: **https://go.dev/dl/**
3. Download: **go1.21.5.windows-amd64.msi** (or latest version)
4. Run the installer
5. Click "Next" through the installation (use default location: `C:\Program Files\Go`)

### Verify Installation
Open a **NEW** Command Prompt (important - must be new!) and run:
```cmd
go version
```

**Expected output:**
```
go version go1.21.5 windows/amd64
```

If you see this, Go is installed! ✅

---

## Step 2: Build the Project

### Navigate to Project
```cmd
cd C:\Users\mobol\Downloads\stream-audio
```

### Build Gateway
```cmd
go build -o bin\gateway.exe .\cmd\gateway
```

This will take 30-60 seconds the first time (downloading dependencies).

### Verify Build
```cmd
dir bin\gateway.exe
```

You should see the file listed!

---

## Step 3: Run the Gateway

### Start Gateway
```cmd
bin\gateway.exe
```

**Expected output:**
```
2025/10/19 08:00:00 Voice Gateway starting on localhost:8080
2025/10/19 08:00:00 WebRTC echo server ready
2025/10/19 08:00:00 Open http://localhost:8080 in your browser to test
```

### Test in Browser
1. Open your browser
2. Go to: **http://localhost:8080**
3. Click "Start Echo Test"
4. Allow microphone access
5. Speak and hear yourself!

---

## Step 4: (Optional) Start NATS

The echo server works WITHOUT NATS, but if you want the full system:

### Install Docker Desktop
1. Download from: **https://www.docker.com/products/docker-desktop/**
2. Install and start Docker Desktop
3. Wait for it to say "Docker is running"

### Start NATS
```cmd
docker run -d --name voice-nats -p 4222:4222 -p 8222:8222 nats:latest -js
```

### Verify NATS
```cmd
curl http://localhost:8222/varz
```

Or open in browser: http://localhost:8222/varz

---

## Troubleshooting

### "go: command not found" after installing

**Problem:** Command Prompt was open before installing Go

**Solution:**
1. Close ALL Command Prompt windows
2. Open a NEW Command Prompt
3. Try again: `go version`

### "curl: command not found"

**Problem:** Windows doesn't have curl in older versions

**Solution:** Just open http://localhost:8222/varz in your browser instead

### Build errors about dependencies

**Problem:** Network issue or first-time setup

**Solution:**
```cmd
go mod download
go mod tidy
go build -o bin\gateway.exe .\cmd\gateway
```

### Port 8080 in use

**Problem:** Another program using port 8080

**Solution:** Use different port
```cmd
set SERVER_PORT=8081
bin\gateway.exe
```

Then open http://localhost:8081

---

## Complete Start-to-Finish Commands

Once Go is installed, here's everything in order:

```cmd
REM 1. Navigate to project
cd C:\Users\mobol\Downloads\stream-audio

REM 2. Build gateway
go build -o bin\gateway.exe .\cmd\gateway

REM 3. Run gateway
bin\gateway.exe

REM 4. Open browser to http://localhost:8080
```

---

## Quick Start Batch File (After Go is Installed)

Once Go is installed, you can use:
```cmd
cd C:\Users\mobol\Downloads\stream-audio
SIMPLE_START.bat
```

This will build and run automatically!
