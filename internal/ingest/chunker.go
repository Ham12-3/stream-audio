package ingest

import (
	"encoding/binary"
	"io"
	"log"
	"math"
	"time"

	"github.com/pion/rtp"
)

const (
	// Standard frame sizes for audio processing
	FrameDuration20ms = 20 * time.Millisecond
	FrameDuration40ms = 40 * time.Millisecond

	// Opus sample rate (WebRTC default)
	OpusSampleRate = 48000

	// PCM sample rate (common for ASR)
	PCMSampleRate = 16000
)

// Chunker processes audio frames and chunks them for downstream processing
type Chunker struct {
	sampleRate    int
	frameDuration time.Duration
	samplesPerFrame int
	buffer        []byte
	onChunk       func([]byte)
}

// NewChunker creates a new audio chunker
func NewChunker(sampleRate int, frameDuration time.Duration, onChunk func([]byte)) *Chunker {
	samplesPerFrame := int(float64(sampleRate) * frameDuration.Seconds())

	return &Chunker{
		sampleRate:      sampleRate,
		frameDuration:   frameDuration,
		samplesPerFrame: samplesPerFrame,
		buffer:          make([]byte, 0, samplesPerFrame*2), // 16-bit samples
		onChunk:         onChunk,
	}
}

// ProcessRTP processes an RTP packet and chunks the audio
func (c *Chunker) ProcessRTP(packet *rtp.Packet) error {
	// Add payload to buffer
	c.buffer = append(c.buffer, packet.Payload...)

	// Calculate bytes per frame (16-bit samples = 2 bytes per sample)
	bytesPerFrame := c.samplesPerFrame * 2

	// Process complete frames
	for len(c.buffer) >= bytesPerFrame {
		chunk := c.buffer[:bytesPerFrame]
		c.buffer = c.buffer[bytesPerFrame:]

		if c.onChunk != nil {
			c.onChunk(chunk)
		}
	}

	return nil
}

// VAD implements simple Voice Activity Detection
type VAD struct {
	threshold     float64
	minSilenceDur time.Duration
	silenceStart  time.Time
	isSpeaking    bool
	onSpeechStart func()
	onSpeechEnd   func()
}

// NewVAD creates a new Voice Activity Detector
func NewVAD(threshold float64, minSilenceDur time.Duration) *VAD {
	return &VAD{
		threshold:     threshold,
		minSilenceDur: minSilenceDur,
	}
}

// SetCallbacks sets the speech event callbacks
func (v *VAD) SetCallbacks(onStart, onEnd func()) {
	v.onSpeechStart = onStart
	v.onSpeechEnd = onEnd
}

// Process analyzes an audio chunk for voice activity
func (v *VAD) Process(chunk []byte) bool {
	energy := calculateEnergy(chunk)

	if energy > v.threshold {
		// Speech detected
		if !v.isSpeaking {
			v.isSpeaking = true
			v.silenceStart = time.Time{}
			if v.onSpeechStart != nil {
				v.onSpeechStart()
			}
		}
		return true
	}

	// Silence detected
	if v.isSpeaking {
		if v.silenceStart.IsZero() {
			v.silenceStart = time.Now()
		} else if time.Since(v.silenceStart) > v.minSilenceDur {
			v.isSpeaking = false
			if v.onSpeechEnd != nil {
				v.onSpeechEnd()
			}
		}
	}

	return false
}

// IsSpeaking returns whether speech is currently detected
func (v *VAD) IsSpeaking() bool {
	return v.isSpeaking
}

// calculateEnergy computes the RMS energy of an audio chunk
func calculateEnergy(pcm []byte) float64 {
	if len(pcm) < 2 {
		return 0
	}

	var sum float64
	samples := len(pcm) / 2

	for i := 0; i < samples; i++ {
		sample := int16(binary.LittleEndian.Uint16(pcm[i*2:]))
		normalized := float64(sample) / 32768.0
		sum += normalized * normalized
	}

	rms := math.Sqrt(sum / float64(samples))
	return rms
}

// FrameReader reads audio frames from an io.Reader
type FrameReader struct {
	reader       io.Reader
	frameSize    int
	buffer       []byte
}

// NewFrameReader creates a new frame reader
func NewFrameReader(reader io.Reader, frameSize int) *FrameReader {
	return &FrameReader{
		reader:    reader,
		frameSize: frameSize,
		buffer:    make([]byte, frameSize),
	}
}

// ReadFrame reads a single frame
func (fr *FrameReader) ReadFrame() ([]byte, error) {
	n, err := io.ReadFull(fr.reader, fr.buffer)
	if err != nil {
		return nil, err
	}

	if n != fr.frameSize {
		return nil, io.ErrUnexpectedEOF
	}

	// Return a copy to avoid buffer reuse issues
	frame := make([]byte, n)
	copy(frame, fr.buffer[:n])
	return frame, nil
}

// FanOut duplicates audio chunks to multiple channels
func FanOut(in <-chan []byte, outs ...chan<- []byte) {
	for chunk := range in {
		for _, out := range outs {
			select {
			case out <- chunk:
			default:
				log.Printf("Warning: dropping audio chunk (channel full)")
			}
		}
	}

	// Close all output channels when input closes
	for _, out := range outs {
		close(out)
	}
}
