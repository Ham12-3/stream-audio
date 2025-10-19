package bus

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Client wraps NATS JetStream client
type Client struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	ctx    context.Context
	cancel context.CancelFunc
}

// Message represents a message on the bus
type Message struct {
	SessionID string
	Data      []byte
	Timestamp time.Time
}

// NewClient creates a new NATS client
func NewClient(url string) (*Client, error) {
	// Connect to NATS
	nc, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		nc:     nc,
		js:     js,
		ctx:    ctx,
		cancel: cancel,
	}

	// Initialize streams
	if err := client.initStreams(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to initialize streams: %w", err)
	}

	log.Println("Connected to NATS JetStream")
	return client, nil
}

// initStreams creates the required JetStream streams
func (c *Client) initStreams() error {
	// Create AUDIO stream for audio frames
	_, err := c.js.CreateOrUpdateStream(c.ctx, jetstream.StreamConfig{
		Name:        "AUDIO",
		Subjects:    []string{"voice.audio.>"},
		Retention:   jetstream.WorkQueuePolicy,
		MaxAge:      time.Hour,
		Storage:     jetstream.MemoryStorage,
		Replicas:    1,
		Description: "Audio frames stream",
	})
	if err != nil {
		return fmt.Errorf("failed to create AUDIO stream: %w", err)
	}

	// Create TEXT stream for transcripts
	_, err = c.js.CreateOrUpdateStream(c.ctx, jetstream.StreamConfig{
		Name:        "TEXT",
		Subjects:    []string{"voice.text.>"},
		Retention:   jetstream.WorkQueuePolicy,
		MaxAge:      time.Hour,
		Storage:     jetstream.MemoryStorage,
		Replicas:    1,
		Description: "Transcripts stream",
	})
	if err != nil {
		return fmt.Errorf("failed to create TEXT stream: %w", err)
	}

	// Create TTS stream for synthesized audio
	_, err = c.js.CreateOrUpdateStream(c.ctx, jetstream.StreamConfig{
		Name:        "TTS",
		Subjects:    []string{"voice.tts.>"},
		Retention:   jetstream.WorkQueuePolicy,
		MaxAge:      time.Hour,
		Storage:     jetstream.MemoryStorage,
		Replicas:    1,
		Description: "TTS audio stream",
	})
	if err != nil {
		return fmt.Errorf("failed to create TTS stream: %w", err)
	}

	return nil
}

// PublishAudio publishes an audio frame to the bus
func (c *Client) PublishAudio(sessionID string, data []byte) error {
	subject := fmt.Sprintf("voice.audio.%s", sessionID)
	_, err := c.js.Publish(c.ctx, subject, data)
	return err
}

// PublishText publishes a transcript to the bus
func (c *Client) PublishText(sessionID string, text []byte) error {
	subject := fmt.Sprintf("voice.text.%s", sessionID)
	_, err := c.js.Publish(c.ctx, subject, text)
	return err
}

// PublishTTS publishes synthesized audio to the bus
func (c *Client) PublishTTS(sessionID string, data []byte) error {
	subject := fmt.Sprintf("voice.tts.%s", sessionID)
	_, err := c.js.Publish(c.ctx, subject, data)
	return err
}

// SubscribeAudio subscribes to audio frames for a session
func (c *Client) SubscribeAudio(sessionID string, handler func(*Message)) error {
	subject := fmt.Sprintf("voice.audio.%s", sessionID)

	cons, err := c.js.CreateOrUpdateConsumer(c.ctx, "AUDIO", jetstream.ConsumerConfig{
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		handler(&Message{
			SessionID: sessionID,
			Data:      msg.Data(),
			Timestamp: time.Now(),
		})
		msg.Ack()
	})

	return err
}

// SubscribeText subscribes to transcripts for a session
func (c *Client) SubscribeText(sessionID string, handler func(*Message)) error {
	subject := fmt.Sprintf("voice.text.%s", sessionID)

	cons, err := c.js.CreateOrUpdateConsumer(c.ctx, "TEXT", jetstream.ConsumerConfig{
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		handler(&Message{
			SessionID: sessionID,
			Data:      msg.Data(),
			Timestamp: time.Now(),
		})
		msg.Ack()
	})

	return err
}

// SubscribeTTS subscribes to synthesized audio for a session
func (c *Client) SubscribeTTS(sessionID string, handler func(*Message)) error {
	subject := fmt.Sprintf("voice.tts.%s", sessionID)

	cons, err := c.js.CreateOrUpdateConsumer(c.ctx, "TTS", jetstream.ConsumerConfig{
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		handler(&Message{
			SessionID: sessionID,
			Data:      msg.Data(),
			Timestamp: time.Now(),
		})
		msg.Ack()
	})

	return err
}

// Close closes the NATS connection
func (c *Client) Close() {
	c.cancel()
	if c.nc != nil {
		c.nc.Close()
	}
}
