package session

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Recorder handles recording of audio and transcripts
type Recorder struct {
	sessionID    string
	recordingDir string
	audioFile    *os.File
	metaFile     *os.File
	transcripts  []TranscriptEntry
	mu           sync.Mutex
}

// TranscriptEntry represents a transcript with timing
type TranscriptEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	IsFinal   bool      `json:"is_final"`
	Speaker   string    `json:"speaker"` // "user" or "assistant"
}

// RecordingMetadata contains session metadata
type RecordingMetadata struct {
	SessionID   string            `json:"session_id"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Duration    float64           `json:"duration_seconds"`
	SampleRate  int               `json:"sample_rate"`
	Channels    int               `json:"channels"`
	Transcripts []TranscriptEntry `json:"transcripts"`
	Stats       RecordingStats    `json:"stats"`
}

// RecordingStats contains statistics about the recording
type RecordingStats struct {
	TotalAudioBytes  int64 `json:"total_audio_bytes"`
	TranscriptCount  int   `json:"transcript_count"`
	SpeechDuration   float64 `json:"speech_duration_seconds"`
}

// NewRecorder creates a new session recorder
func NewRecorder(sessionID string, recordingDir string) (*Recorder, error) {
	// Create recording directory if it doesn't exist
	if err := os.MkdirAll(recordingDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create recording directory: %w", err)
	}

	sessionDir := filepath.Join(recordingDir, sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	// Create audio file (PCM format)
	audioPath := filepath.Join(sessionDir, "audio.pcm")
	audioFile, err := os.Create(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio file: %w", err)
	}

	// Create metadata file
	metaPath := filepath.Join(sessionDir, "metadata.json")
	metaFile, err := os.Create(metaPath)
	if err != nil {
		audioFile.Close()
		return nil, fmt.Errorf("failed to create metadata file: %w", err)
	}

	return &Recorder{
		sessionID:    sessionID,
		recordingDir: sessionDir,
		audioFile:    audioFile,
		metaFile:     metaFile,
		transcripts:  []TranscriptEntry{},
	}, nil
}

// WriteAudio writes audio data to the recording
func (r *Recorder) WriteAudio(data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.audioFile == nil {
		return fmt.Errorf("audio file not open")
	}

	_, err := r.audioFile.Write(data)
	return err
}

// AddTranscript adds a transcript entry
func (r *Recorder) AddTranscript(text string, isFinal bool, speaker string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.transcripts = append(r.transcripts, TranscriptEntry{
		Timestamp: time.Now(),
		Text:      text,
		IsFinal:   isFinal,
		Speaker:   speaker,
	})
}

// Close finalizes the recording and writes metadata
func (r *Recorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Close audio file
	if r.audioFile != nil {
		if err := r.audioFile.Close(); err != nil {
			return fmt.Errorf("failed to close audio file: %w", err)
		}
	}

	// Write metadata
	if err := r.writeMetadata(); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Close metadata file
	if r.metaFile != nil {
		if err := r.metaFile.Close(); err != nil {
			return fmt.Errorf("failed to close metadata file: %w", err)
		}
	}

	return nil
}

// writeMetadata writes session metadata to file
func (r *Recorder) writeMetadata() error {
	if r.metaFile == nil {
		return fmt.Errorf("metadata file not open")
	}

	// Get audio file size
	audioInfo, err := os.Stat(filepath.Join(r.recordingDir, "audio.pcm"))
	if err != nil {
		return fmt.Errorf("failed to stat audio file: %w", err)
	}

	// Calculate duration (assuming 16kHz, 16-bit, mono)
	sampleRate := 16000
	bytesPerSample := 2
	totalSamples := audioInfo.Size() / int64(bytesPerSample)
	duration := float64(totalSamples) / float64(sampleRate)

	metadata := RecordingMetadata{
		SessionID:   r.sessionID,
		StartTime:   time.Now().Add(-time.Duration(duration) * time.Second),
		EndTime:     time.Now(),
		Duration:    duration,
		SampleRate:  sampleRate,
		Channels:    1,
		Transcripts: r.transcripts,
		Stats: RecordingStats{
			TotalAudioBytes: audioInfo.Size(),
			TranscriptCount: len(r.transcripts),
			SpeechDuration:  duration,
		},
	}

	encoder := json.NewEncoder(r.metaFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metadata)
}

// ExportToWAV converts the PCM recording to WAV format
func (r *Recorder) ExportToWAV(outputPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pcmPath := filepath.Join(r.recordingDir, "audio.pcm")
	pcmFile, err := os.Open(pcmPath)
	if err != nil {
		return fmt.Errorf("failed to open PCM file: %w", err)
	}
	defer pcmFile.Close()

	wavFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create WAV file: %w", err)
	}
	defer wavFile.Close()

	// Get PCM file size
	pcmInfo, err := pcmFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat PCM file: %w", err)
	}

	// WAV header parameters
	sampleRate := uint32(16000)
	bitsPerSample := uint16(16)
	numChannels := uint16(1)
	dataSize := uint32(pcmInfo.Size())

	// Write WAV header
	if err := writeWAVHeader(wavFile, sampleRate, bitsPerSample, numChannels, dataSize); err != nil {
		return fmt.Errorf("failed to write WAV header: %w", err)
	}

	// Copy PCM data
	if _, err := io.Copy(wavFile, pcmFile); err != nil {
		return fmt.Errorf("failed to copy audio data: %w", err)
	}

	return nil
}

// writeWAVHeader writes a WAV file header
func writeWAVHeader(w io.Writer, sampleRate uint32, bitsPerSample, numChannels uint16, dataSize uint32) error {
	byteRate := sampleRate * uint32(numChannels) * uint32(bitsPerSample/8)
	blockAlign := numChannels * (bitsPerSample / 8)

	// RIFF header
	if _, err := w.Write([]byte("RIFF")); err != nil {
		return err
	}
	if err := writeUint32(w, dataSize+36); err != nil {
		return err
	}
	if _, err := w.Write([]byte("WAVE")); err != nil {
		return err
	}

	// fmt chunk
	if _, err := w.Write([]byte("fmt ")); err != nil {
		return err
	}
	if err := writeUint32(w, 16); err != nil { // chunk size
		return err
	}
	if err := writeUint16(w, 1); err != nil { // audio format (PCM)
		return err
	}
	if err := writeUint16(w, numChannels); err != nil {
		return err
	}
	if err := writeUint32(w, sampleRate); err != nil {
		return err
	}
	if err := writeUint32(w, byteRate); err != nil {
		return err
	}
	if err := writeUint16(w, blockAlign); err != nil {
		return err
	}
	if err := writeUint16(w, bitsPerSample); err != nil {
		return err
	}

	// data chunk
	if _, err := w.Write([]byte("data")); err != nil {
		return err
	}
	if err := writeUint32(w, dataSize); err != nil {
		return err
	}

	return nil
}

func writeUint16(w io.Writer, val uint16) error {
	bytes := []byte{byte(val), byte(val >> 8)}
	_, err := w.Write(bytes)
	return err
}

func writeUint32(w io.Writer, val uint32) error {
	bytes := []byte{byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24)}
	_, err := w.Write(bytes)
	return err
}
