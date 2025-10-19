# Next Steps - Voice Gateway Implementation Guide

This document outlines the immediate next steps to evolve this project from an MVP to a production-ready voice agent platform.

## Immediate Priorities (Week 1-2)

### 1. Test the Echo Server
```bash
# Build and run
make build
docker run -p 4222:4222 nats:latest -js &
./bin/gateway

# Open browser to http://localhost:8080
# Verify audio echo works
```

### 2. Integrate Real ASR (Deepgram Recommended)

**Why Deepgram?**
- WebSocket streaming API (easy integration)
- Low latency (~200ms first partial)
- Excellent accuracy
- Generous free tier

**Implementation Steps:**

1. **Sign up**: https://deepgram.com
2. **Get API key**: Save to `.env` as `DEEPGRAM_API_KEY`
3. **Update ASR worker** (`cmd/asr-worker/main.go`):

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/gorilla/websocket"
    "voice-gateway/internal/bus"
    "voice-gateway/internal/config"
)

func main() {
    cfg := config.Load()
    busClient, _ := bus.NewClient(cfg.NATS.URL)
    defer busClient.Close()

    // Connect to Deepgram
    apiKey := os.Getenv("DEEPGRAM_API_KEY")
    url := "wss://api.deepgram.com/v1/listen?encoding=linear16&sample_rate=16000&channels=1"

    header := http.Header{}
    header.Set("Authorization", "Token "+apiKey)

    conn, _, err := websocket.DefaultDialer.Dial(url, header)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Subscribe to audio from NATS
    busClient.SubscribeAudio("*", func(msg *bus.Message) {
        // Send audio to Deepgram
        conn.WriteMessage(websocket.BinaryMessage, msg.Data)
    })

    // Receive transcripts from Deepgram
    for {
        var result DeepgramResponse
        if err := conn.ReadJSON(&result); err != nil {
            log.Printf("Read error: %v", err)
            break
        }

        if len(result.Channel.Alternatives) > 0 {
            transcript := result.Channel.Alternatives[0].Transcript
            if transcript != "" {
                // Publish to NATS
                data, _ := json.Marshal(TranscriptMessage{
                    SessionID: msg.SessionID,
                    Text:      transcript,
                    IsFinal:   result.IsFinal,
                })
                busClient.PublishText(msg.SessionID, data)
            }
        }
    }
}

type DeepgramResponse struct {
    Channel struct {
        Alternatives []struct {
            Transcript string  `json:"transcript"`
            Confidence float64 `json:"confidence"`
        } `json:"alternatives"`
    } `json:"channel"`
    IsFinal bool `json:"is_final"`
}
```

### 3. Integrate Real TTS (ElevenLabs Recommended)

**Why ElevenLabs?**
- Best voice quality
- Streaming API with <300ms latency
- `eleven_turbo_v2` model is perfect for real-time

**Implementation Steps:**

1. **Sign up**: https://elevenlabs.io
2. **Get API key**: Save to `.env` as `ELEVENLABS_API_KEY`
3. **Choose voice**: Browse voices and copy voice ID
4. **Update TTS worker** (`cmd/tts-worker/main.go`):

```go
package main

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "os"
)

func main() {
    cfg := config.Load()
    busClient, _ := bus.NewClient(cfg.NATS.URL)
    defer busClient.Close()

    // Subscribe to text from NATS
    busClient.SubscribeText("*", func(msg *bus.Message) {
        var textMsg TextMessage
        json.Unmarshal(msg.Data, &textMsg)

        synthesizeAndPublish(textMsg, busClient)
    })

    select {}
}

func synthesizeAndPublish(textMsg TextMessage, busClient *bus.Client) {
    apiKey := os.Getenv("ELEVENLABS_API_KEY")
    voiceID := os.Getenv("ELEVENLABS_VOICE_ID") // e.g., "21m00Tcm4TlvDq8ikWAM"

    url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s/stream", voiceID)

    body := map[string]interface{}{
        "text":     textMsg.Text,
        "model_id": "eleven_turbo_v2",
    }

    bodyBytes, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
    req.Header.Set("xi-api-key", apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Printf("ElevenLabs error: %v", err)
        return
    }
    defer resp.Body.Close()

    // Stream audio chunks to NATS
    buffer := make([]byte, 4096)
    for {
        n, err := resp.Body.Read(buffer)
        if n > 0 {
            busClient.PublishTTS(textMsg.SessionID, buffer[:n])
        }
        if err == io.EOF {
            break
        }
    }
}
```

### 4. Wire Up Full Pipeline in Gateway

Update `internal/webrtc/handler.go` to:
1. Send audio to NATS instead of echoing
2. Subscribe to TTS stream from NATS
3. Play TTS audio back to user

```go
// In HandleOffer function, replace echo logic:

peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
    log.Printf("Session %s: Received track", sess.ID)

    // Send audio to NATS for ASR processing
    go func() {
        for {
            rtp, _, err := track.ReadRTP()
            if err != nil {
                return
            }

            // Publish to NATS
            busClient.PublishAudio(sess.ID, rtp.Payload)
        }
    }()

    // Subscribe to TTS output from NATS
    busClient.SubscribeTTS(sess.ID, func(msg *bus.Message) {
        // Write TTS audio to local track
        localTrack.Write(msg.Data)
    })
})
```

## Advanced Features (Week 3-4)

### 5. Add LLM Conversation Loop

Create `cmd/llm-worker/main.go`:

```go
package main

import (
    "voice-gateway/internal/bus"
    "voice-gateway/internal/llm"
)

func main() {
    busClient, _ := bus.NewClient(cfg.NATS.URL)
    llmHandler := llm.NewHandler(
        os.Getenv("LLM_API_URL"),
        os.Getenv("LLM_API_KEY"),
        "gpt-4o-mini",
    )

    conversations := make(map[string]*llm.ConversationContext)

    // Subscribe to transcripts
    busClient.SubscribeText("*", func(msg *bus.Message) {
        var transcript TranscriptMessage
        json.Unmarshal(msg.Data, &transcript)

        if !transcript.IsFinal {
            return // Only process final transcripts
        }

        // Get or create conversation
        conv, ok := conversations[transcript.SessionID]
        if !ok {
            conv = llm.NewConversationContext(
                transcript.SessionID,
                "You are a helpful voice assistant.",
            )
            conversations[transcript.SessionID] = conv
        }

        conv.AddUserMessage(transcript.Text)

        // Stream LLM response
        var fullResponse string
        llmHandler.StreamChat(conv.GetMessages(), func(chunk string) {
            fullResponse += chunk

            // Publish chunk to TTS (sentence by sentence)
            if strings.HasSuffix(chunk, ".") || strings.HasSuffix(chunk, "?") {
                textMsg := TextMessage{
                    SessionID: transcript.SessionID,
                    Text:      fullResponse,
                }
                data, _ := json.Marshal(textMsg)
                busClient.PublishText(transcript.SessionID, data)
                fullResponse = ""
            }
        })

        conv.AddAssistantMessage(fullResponse)
    })

    select {}
}
```

### 6. Add Tool Calling

In LLM worker, integrate skills:

```go
skillRegistry := skills.NewRegistry()
skills.InitDefaultSkills(skillRegistry)

// When LLM returns a tool call:
if toolCall := extractToolCall(llmResponse); toolCall != nil {
    result, err := skillRegistry.Execute(
        context.Background(),
        toolCall.Name,
        toolCall.Params,
    )

    // Send result back to LLM for final response
    conv.AddUserMessage(fmt.Sprintf("Tool result: %v", result))
    // Get LLM response and send to TTS
}
```

### 7. Production Hardening

**Observability:**
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    audioPacketsProcessed = prometheus.NewCounter(...)
    transcriptLatency = prometheus.NewHistogram(...)
    activeConnections = prometheus.NewGauge(...)
)

// In WebRTC handler:
audioPacketsProcessed.Inc()
```

**Authentication:**
```go
// JWT middleware for /offer endpoint
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !validateJWT(token) {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Testing & Optimization

### Performance Testing
```bash
# Install hey for load testing
go install github.com/rakyll/hey@latest

# Test WebRTC signaling
hey -n 1000 -c 10 -m POST http://localhost:8080/offer

# Monitor metrics
curl http://localhost:8222/metrics  # NATS
```

### Audio Quality Optimization
- Implement Opus codec (better compression than PCM)
- Add jitter buffer for network resilience
- Tune VAD thresholds based on environment

### Latency Optimization
- Reduce frame size to 20ms (from 40ms)
- Use WebSocket for signaling (faster than HTTP)
- Deploy ASR/TTS workers closer to users (edge)

## Deployment to Production

### Fly.io (Recommended for Edge)
```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Deploy gateway
fly launch --name voice-gateway
fly deploy

# Deploy workers
cd cmd/asr-worker && fly launch
cd cmd/tts-worker && fly launch
```

### Kubernetes
```yaml
# deploy/k8s/gateway.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: voice-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: voice-gateway
  template:
    metadata:
      labels:
        app: voice-gateway
    spec:
      containers:
      - name: gateway
        image: voice-gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: NATS_URL
          value: "nats://nats:4222"
```

## Portfolio Enhancement Tips

**Make it stand out:**
1. **Live demo**: Deploy to fly.io, add link to README
2. **Demo video**: Record yourself using the voice agent
3. **Metrics dashboard**: Add Grafana dashboard showing latency metrics
4. **Blog post**: Write about WebRTC + streaming architecture
5. **Benchmarks**: Compare latency with commercial solutions

**Technical highlights to mention:**
- Real-time bidirectional audio streaming
- Microservices architecture with message bus
- Goroutines for concurrent audio processing
- Streaming APIs (ASR, LLM, TTS) orchestration
- Production patterns (retry, backpressure, monitoring)

## Resources

**Learning:**
- [WebRTC for the Curious](https://webrtcforthecurious.com/)
- [NATS Patterns](https://docs.nats.io/nats-concepts/core-nats/patterns)
- [Pion WebRTC Examples](https://github.com/pion/webrtc/tree/master/examples)

**APIs:**
- Deepgram Docs: https://developers.deepgram.com/docs
- ElevenLabs Docs: https://elevenlabs.io/docs
- OpenAI Streaming: https://platform.openai.com/docs/api-reference/streaming

**Community:**
- Pion Slack: https://pion.ly/slack
- NATS Slack: https://slack.nats.io

---

**Current Status**: ✅ Echo MVP Complete
**Next Milestone**: 🚧 ASR/TTS Integration
**Time Estimate**: 2-4 weeks to production-ready

Good luck! This is a genuinely impressive portfolio project.
