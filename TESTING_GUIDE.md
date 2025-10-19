# Testing Guide - How to See It Work

## Quick Test (Echo Server)

The simplest way to see the project working is to test the **echo server**. This demonstrates the complete WebRTC audio pipeline without needing any external services.

### Prerequisites

You need:
- ✅ Go already installed (we installed it earlier)
- ✅ A web browser with microphone access
- ✅ (Optional) Docker for NATS

### Option 1: Quick Start (Recommended)

```bash
# Navigate to project directory
cd /mnt/c/Users/mobol/Downloads/stream-audio

# Run the quick start script
./start.sh
```

This will:
1. Check if Go is installed
2. Build the gateway binary (if not already built)
3. Start NATS in Docker (if Docker is available)
4. Start the Voice Gateway server
5. Show you the URL to test

### Option 2: Manual Step-by-Step

If the script doesn't work or you want more control:

#### Step 1: Start NATS (Message Bus)

**With Docker:**
```bash
docker run -d --name voice-nats \
  -p 4222:4222 \
  -p 8222:8222 \
  nats:latest -js
```

**Without Docker:**
- Download NATS from: https://nats.io/download/
- Run: `nats-server -js`

**Verify NATS is running:**
```bash
# Should show NATS metrics
curl http://localhost:8222/varz
```

#### Step 2: Build the Gateway

```bash
# From project root
export PATH=$HOME/go/bin:$PATH
export GOROOT=$HOME/go
export GOPATH=$HOME/go-workspace

# Build the gateway
go build -o bin/gateway ./cmd/gateway

# Verify binary exists
ls -lh bin/gateway
```

#### Step 3: Run the Gateway

```bash
# Start the server
./bin/gateway
```

You should see output like:
```
2025/10/19 08:00:00 Voice Gateway starting on localhost:8080
2025/10/19 08:00:00 WebRTC echo server ready
2025/10/19 08:00:00 Open http://localhost:8080 in your browser to test
```

#### Step 4: Test in Browser

1. **Open your browser** (Chrome or Firefox recommended)
2. **Navigate to**: `http://localhost:8080`
3. **You'll see** a beautiful purple gradient interface
4. **Click** the "Start Echo Test" button
5. **Allow** microphone access when prompted
6. **Speak** into your microphone
7. **You should hear** your voice echoed back with a small delay!

### What You're Testing

When you speak, here's what happens:

```
Your Voice (Microphone)
    ↓
Browser captures audio
    ↓
WebRTC sends audio to Gateway
    ↓
Gateway receives RTP packets
    ↓
Gateway echoes packets back
    ↓
Browser plays audio through speakers
    ↓
You hear yourself!
```

---

## Troubleshooting

### "Connection failed" in browser

**Problem**: Gateway can't establish WebRTC connection

**Solutions**:
1. Check gateway is running: `ps aux | grep gateway`
2. Check port 8080 is available: `netstat -an | grep 8080`
3. Check browser console for errors (F12 → Console)
4. Try a different browser

### "Failed to connect to NATS"

**Problem**: NATS server not running or not accessible

**Solutions**:
1. Check NATS is running: `docker ps | grep nats` or `curl localhost:8222`
2. Restart NATS: `docker restart voice-nats`
3. Check logs: `docker logs voice-nats`
4. Verify port 4222 is open: `netstat -an | grep 4222`

### No audio / Can't hear echo

**Problem**: Audio path issue

**Solutions**:
1. Check microphone permissions in browser (click lock icon in address bar)
2. Check system microphone is working (test in another app)
3. Check browser audio output (try playing YouTube)
4. Look at browser console for WebRTC errors
5. Check gateway logs for RTP packet messages

### "go: command not found"

**Problem**: Go not in PATH

**Solution**:
```bash
export PATH=$HOME/go/bin:$PATH
export GOROOT=$HOME/go
export GOPATH=$HOME/go-workspace

# Verify
go version
```

---

## Advanced Testing

### Test with Multiple Sessions

Open multiple browser tabs to `http://localhost:8080` and start echo in each. Each gets its own session!

**What to observe**:
- Gateway logs show multiple session IDs
- Each session works independently
- No audio crosstalk between sessions

### Test NATS Message Flow

While echo is running, check NATS streams:

```bash
# If NATS CLI is installed
nats stream ls
nats stream info AUDIO
nats stream info TEXT
nats stream info TTS

# Watch messages flow
nats stream view AUDIO
```

### Test Session Management

```bash
# While gateway is running, in another terminal
curl http://localhost:8080/offer

# Should return 405 (Method Not Allowed) - that's good!
# Only POST is allowed
```

### Monitor Gateway Logs

The gateway logs show:
```
Created new session: <UUID>
Session <UUID>: Received track: audio (codec: audio/opus)
Session <UUID>: Connection state changed: connected
```

This shows:
- Session creation
- Audio track reception
- Connection lifecycle

---

## Testing with Docker Compose

For a more production-like test:

```bash
# Build and start all services
docker-compose up --build

# In another terminal, check services
docker-compose ps

# View logs
docker-compose logs -f gateway
docker-compose logs -f nats
docker-compose logs -f asr-worker
docker-compose logs -f tts-worker

# Test
# Open http://localhost:8080

# Stop everything
docker-compose down
```

**What this tests**:
- Multi-container orchestration
- Service discovery
- Network communication
- Container builds

---

## Expected Behavior

### ✅ Working Correctly

**Browser UI:**
- Status shows "Connected - Speak now!"
- Logs show connection progress
- You hear your voice echoed back
- Small delay (~100-300ms) is normal

**Gateway Logs:**
- "Created new session"
- "Received track"
- "Connection state changed: connected"
- No errors

**NATS:**
- Streams created (AUDIO, TEXT, TTS)
- No connection errors

### ❌ Problems

**If you see:**
- "Connection failed" → Check gateway logs
- "Disconnected" → Check NATS and gateway
- No echo → Check audio permissions
- Error messages → Check browser console

---

## Video Demo Script

Want to record a demo? Here's a script:

1. **Show the code structure**
   ```bash
   tree -L 2 -d  # or ls -R
   ```

2. **Start the services**
   ```bash
   ./start.sh
   # Or docker-compose up
   ```

3. **Open browser**
   - Show the beautiful UI
   - Point out the status indicator
   - Show the logs area

4. **Start echo test**
   - Click button
   - Allow microphone
   - Speak: "Hello, this is a real-time voice gateway built with Go and WebRTC"
   - Show the echo working

5. **Show the logs**
   - Terminal: session creation, connection state
   - Browser: WebRTC signaling, ICE state

6. **Stop and explain**
   - Explain the audio flow
   - Show the code for WebRTC handler
   - Explain next steps (ASR/TTS integration)

---

## Performance Testing

### Measure Latency

```javascript
// In browser console while echo is running
let startTime = Date.now();
// Speak a sharp sound (like clap)
// When you hear it back:
let latency = Date.now() - startTime;
console.log('Echo latency:', latency, 'ms');
```

**Expected**: 100-300ms (depending on network, processing)

### Stress Test

Open 10 browser tabs, start echo in all:

```bash
# Monitor gateway resource usage
top -p $(pgrep gateway)

# Expected:
# CPU: 5-20% (moderate)
# Memory: ~50MB per connection
```

---

## Next Level Testing (After ASR/TTS Integration)

Once you integrate real ASR/TTS services:

### 1. Test Speech Recognition
- Speak: "What is the weather today?"
- Check ASR worker logs for transcript
- Verify text published to NATS

### 2. Test Text-to-Speech
- Send text via NATS
- Hear synthesized voice response

### 3. Test Full Loop
- Speak a question
- LLM generates response
- TTS speaks the answer
- Measure end-to-end latency

---

## Demo Talking Points

When showing this project:

**For Technical Audience:**
- "This is a WebRTC gateway handling bidirectional audio streams"
- "Using pion/webrtc for Go, NATS JetStream for message bus"
- "The echo demonstrates the full RTP packet pipeline"
- "Real-time audio processing with 20-40ms frame chunking"
- "Ready to integrate with Deepgram ASR and ElevenLabs TTS"

**For Non-Technical Audience:**
- "This is a voice bot server, like Siri or Alexa"
- "You speak, it hears you, processes your speech, and responds"
- "Built from scratch in Go, handling real-time audio"
- "Currently echoing to demonstrate the audio pipeline works"
- "Next step is adding AI for real conversations"

---

## Quick Reference

### Start Everything
```bash
# Option 1: Script
./start.sh

# Option 2: Docker
docker-compose up

# Option 3: Manual
docker run -d -p 4222:4222 nats:latest -js
./bin/gateway
```

### Stop Everything
```bash
# Ctrl+C to stop gateway

# Stop NATS
docker stop voice-nats

# Or with Docker Compose
docker-compose down
```

### Check Status
```bash
# Gateway running?
ps aux | grep gateway

# NATS running?
curl http://localhost:8222/varz

# Ports open?
netstat -an | grep -E '8080|4222'
```

### Clean Up
```bash
# Remove NATS container
docker rm -f voice-nats

# Remove Docker Compose stack
docker-compose down -v

# Clean builds
make clean
```

---

## Success Checklist

Before demoing, verify:

- [ ] NATS is running (check port 4222)
- [ ] Gateway binary built (ls bin/gateway)
- [ ] Gateway is running (ps aux | grep gateway)
- [ ] Browser can access http://localhost:8080
- [ ] Microphone permissions granted
- [ ] Audio input/output working
- [ ] Echo works (you can hear yourself)
- [ ] Logs show session creation
- [ ] No errors in console

---

## Recording the Demo

If making a video:

1. **Terminal window**: Show gateway logs
2. **Browser window**: Show the UI
3. **Code editor**: Show key files (handler.go)
4. **Flow**: Start → Test → Explain → Show code → Stop

**Duration**: 2-3 minutes is perfect

**Key message**: "Real-time voice processing with WebRTC in Go, production-ready foundation, ready for AI integration"

---

You're all set to test! Start with the Quick Start and work your way through. The echo server is the best way to show "it works!"
