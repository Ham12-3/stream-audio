# Voice Gateway - Realtime Voice Agent Platform

A production-ready, low-latency voice bot server built with Go, WebRTC, and streaming AI services. Stream user audio from the browser, transcribe in real-time, process with LLMs, and stream synthesized speech back—all in under 300ms per chunk.

## Features

- **WebRTC Audio Pipeline**: Bidirectional audio streaming using pion/webrtc
- **Streaming ASR**: Real-time speech recognition with support for multiple providers
- **LLM Integration**: OpenAI-compatible API support for conversational AI
- **Streaming TTS**: Low-latency text-to-speech synthesis
- **Message Bus**: NATS JetStream for scalable, fault-tolerant message routing
- **Session Management**: Full session lifecycle tracking and recording
- **Plugin System**: Hot-swappable skills (tools) for extending agent capabilities
- **Session Recording**: Save audio (PCM/WAV) and transcripts for analytics
- **Production Ready**: Docker support, observability hooks, configuration management

## Architecture

```
┌─────────────┐
│   Browser   │
│  (WebRTC)   │
└──────┬──────┘
       │ Bidirectional Audio
       ▼
┌─────────────────────────────────────┐
│      Voice Gateway (Go)             │
│  ┌──────────────────────────────┐  │
│  │  WebRTC Handler              │  │
│  │  - Peer Connection Mgmt      │  │
│  │  - Audio Track Processing    │  │
│  └──────────┬───────────────────┘  │
│             │                       │
│  ┌──────────▼───────────────────┐  │
│  │  Audio Ingest Pipeline       │  │
│  │  - Frame Chunker (20-40ms)   │  │
│  │  - VAD (Voice Activity)      │  │
│  │  - Fan-out to Workers        │  │
│  └──────────┬───────────────────┘  │
└─────────────┼───────────────────────┘
              │
              ▼
     ┌────────────────┐
     │ NATS JetStream │
     │   Message Bus  │
     └────────┬───────┘
              │
    ┌─────────┼──────────┐
    │         │          │
    ▼         ▼          ▼
┌───────┐ ┌─────┐  ┌─────────┐
│  ASR  │ │ LLM │  │   TTS   │
│Worker │ │ Svc │  │ Worker  │
└───────┘ └─────┘  └─────────┘
```

## Quick Start

### Prerequisites

- Go 1.21+ (or use Docker)
- NATS Server with JetStream (provided via docker-compose)

### Option 1: Run with Docker (Recommended)

```bash
# Clone the repository
git clone <repo-url>
cd voice-gateway

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f gateway

# Stop services
docker-compose down
```

Access the web UI at: http://localhost:8080

### Option 2: Run Locally

```bash
# Install dependencies
make deps

# Start NATS (in separate terminal)
docker run -p 4222:4222 -p 8222:8222 nats:latest -js

# Build the gateway
make build

# Run the gateway
./bin/gateway

# In separate terminals, run the workers (optional for echo mode)
go run ./cmd/asr-worker &
go run ./cmd/tts-worker &
```

## Project Structure

```
voice-gateway/
├── cmd/
│   ├── gateway/           # Main WebRTC gateway server
│   ├── asr-worker/        # ASR worker (stub)
│   └── tts-worker/        # TTS worker (stub)
├── internal/
│   ├── webrtc/            # WebRTC peer connection handling
│   ├── ingest/            # Audio chunking, VAD, frame processing
│   ├── bus/               # NATS JetStream client
│   ├── llm/               # LLM handler (OpenAI-compatible)
│   ├── skills/            # Plugin/skill system
│   ├── session/           # Session management & recording
│   └── config/            # Configuration management
├── pkg/
│   └── proto/             # gRPC protocol definitions
├── web/
│   └── static/            # Web UI for testing
├── deploy/
│   └── docker/            # Dockerfiles
├── docker-compose.yml     # Full stack orchestration
├── Makefile              # Build automation
└── README.md
```

## Development Milestones

The project is structured in incremental milestones:

### ✅ Milestone 1: Echo MVP (Complete)
- WebRTC bidirectional audio
- Simple echo server (audio in → audio out)
- Web UI for testing

### ✅ Milestone 2: Infrastructure (Complete)
- NATS JetStream integration
- Audio frame chunking (20-40ms)
- Session management
- Worker stubs (ASR/TTS)

### ✅ Milestone 3: Core Components (Complete)
- LLM integration framework
- Skills/plugin system
- Session recording
- Configuration management

### 🚧 Milestone 4: Real ASR/TTS Integration (Next)
**Integrate with real services:**

#### ASR Options:
- **Deepgram**: WebSocket streaming, excellent accuracy
- **AssemblyAI**: Real-time API with interim results
- **Whisper**: Self-hosted via faster-whisper or whisper.cpp
- **Google Cloud Speech**: gRPC streaming
- **Azure Speech**: Streaming SDK

#### TTS Options:
- **ElevenLabs**: Streaming API, best quality, <300ms latency
- **OpenAI TTS**: Streaming support, good quality
- **Google Cloud TTS**: Wide language support
- **Azure Speech**: Neural voices
- **Coqui TTS**: Self-hosted, VITS models

### 🚧 Milestone 5: Production Hardening
- [ ] Observability (Prometheus, OpenTelemetry)
- [ ] Authentication & authorization
- [ ] Rate limiting
- [ ] Retry logic & error handling
- [ ] Load testing & optimization
- [ ] Kubernetes manifests

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
# Server
SERVER_HOST=localhost
SERVER_PORT=8080

# WebRTC
STUN_SERVER=stun:stun.l.google.com:19302

# NATS
NATS_URL=nats://localhost:4222

# Services
ASR_URL=localhost:50051
TTS_URL=localhost:50052

# LLM (OpenAI-compatible)
LLM_API_URL=https://api.openai.com/v1/chat/completions
LLM_API_KEY=your-key-here
LLM_MODEL=gpt-4o-mini

# Recording
RECORDING_ENABLED=true
RECORDING_DIR=./recordings
```

## Usage Examples

### Testing the Echo Server

1. Start the gateway: `make run`
2. Open http://localhost:8080
3. Click "Start Echo Test"
4. Allow microphone access
5. Speak and hear yourself echoed back

### Integrating Real ASR

Replace the stub in `cmd/asr-worker/main.go`:

```go
// Example: Deepgram integration
func processAudioStream(sessionID string, busClient *bus.Client) {
    // Connect to Deepgram
    conn, _ := websocket.Dial("wss://api.deepgram.com/v1/listen")

    // Subscribe to audio frames from NATS
    busClient.SubscribeAudio(sessionID, func(msg *bus.Message) {
        // Send audio to Deepgram
        conn.WriteMessage(websocket.BinaryMessage, msg.Data)
    })

    // Receive transcripts
    for {
        var result DeepgramResult
        conn.ReadJSON(&result)

        // Publish to NATS
        transcript := TranscriptMessage{
            SessionID: sessionID,
            Text:      result.Channel.Alternatives[0].Transcript,
            IsFinal:   result.IsFinal,
        }

        data, _ := json.Marshal(transcript)
        busClient.PublishText(sessionID, data)
    }
}
```

### Integrating Real TTS

Replace the stub in `cmd/tts-worker/main.go`:

```go
// Example: ElevenLabs integration
func synthesizeText(text string, sessionID string, busClient *bus.Client) {
    url := "https://api.elevenlabs.io/v1/text-to-speech/VOICE_ID/stream"

    body := map[string]interface{}{
        "text":     text,
        "model_id": "eleven_turbo_v2",
    }

    // Stream audio chunks back
    resp, _ := http.Post(url, "application/json", bodyReader)

    buffer := make([]byte, 4096)
    for {
        n, _ := resp.Body.Read(buffer)
        if n > 0 {
            // Publish audio chunk to NATS
            busClient.PublishTTS(sessionID, buffer[:n])
        }
    }
}
```

## Skills/Plugins

Add custom skills to extend the voice agent:

```go
type WeatherSkill struct{}

func (s *WeatherSkill) Name() string {
    return "get_weather"
}

func (s *WeatherSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    location := params["location"].(string)

    // Call weather API
    weather := fetchWeather(location)

    return map[string]interface{}{
        "temperature": weather.Temp,
        "conditions":  weather.Description,
    }, nil
}

// Register in main.go
skillRegistry := skills.NewRegistry()
skillRegistry.Register(&WeatherSkill{})
```

## API Endpoints

### POST /offer
Creates a WebRTC peer connection.

**Request:**
```json
{
  "sdp": "{\"type\":\"offer\",\"sdp\":\"...\"}"
}
```

**Response:**
```json
{
  "sdp": "{\"type\":\"answer\",\"sdp\":\"...\"}"
}
```

### GET /
Serves the web UI.

## Message Bus Topics

| Subject | Purpose | Producer | Consumer |
|---------|---------|----------|----------|
| `voice.audio.<session>` | Audio frames | Gateway | ASR Worker |
| `voice.text.<session>` | Transcripts | ASR Worker | LLM/Gateway |
| `voice.tts.<session>` | Synthesized audio | TTS Worker | Gateway |

## Performance

Target latencies:
- **WebRTC RTT**: <50ms
- **ASR First Partial**: <200ms
- **LLM First Token**: <300ms
- **TTS First Byte**: <300ms
- **End-to-End**: <600ms (user speech → bot response)

## Testing

```bash
# Run tests
make test

# Build all components
make build

# Clean build artifacts
make clean
```

## Deployment

### Docker

```bash
docker-compose up -d
```

### Kubernetes (Coming Soon)

```bash
kubectl apply -f deploy/k8s/
```

## Roadmap

- [ ] Real ASR integration (Deepgram/Whisper)
- [ ] Real TTS integration (ElevenLabs/OpenAI)
- [ ] Full LLM conversation loop
- [ ] Tool calling / function execution
- [ ] Multi-party rooms (SFU)
- [ ] WASM plugin sandboxing
- [ ] Analytics dashboard
- [ ] Edge deployment (Fly.io, Cloudflare)
- [ ] Mobile SDK

## Contributing

This is a portfolio/learning project showcasing:
- Go concurrency patterns (goroutines, channels)
- WebRTC media handling
- Streaming architecture
- Microservices with message bus
- Production Go project structure

## License

MIT

## Resources

- [Pion WebRTC](https://github.com/pion/webrtc)
- [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
- [Deepgram API](https://developers.deepgram.com/)
- [ElevenLabs API](https://elevenlabs.io/docs/api-reference)

---

**Status**: 🚀 MVP Complete - Ready for ASR/TTS integration

Built with Go, WebRTC, NATS, and streaming AI services.
