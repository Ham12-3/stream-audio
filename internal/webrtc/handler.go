package webrtc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/pion/webrtc/v4"
	"voice-gateway/internal/session"
)

// Handler manages WebRTC peer connections
type Handler struct {
	config         *webrtc.Configuration
	sessionManager *session.Manager
	mu             sync.RWMutex
}

// NewHandler creates a new WebRTC handler
func NewHandler(iceServers []string, sessionMgr *session.Manager) *Handler {
	config := &webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{},
	}

	for _, server := range iceServers {
		config.ICEServers = append(config.ICEServers, webrtc.ICEServer{
			URLs: []string{server},
		})
	}

	return &Handler{
		config:         config,
		sessionManager: sessionMgr,
	}
}

// HandleOffer processes a WebRTC offer and returns an answer
func (h *Handler) HandleOffer(offerJSON string) (string, error) {
	// Create a new session
	sess := h.sessionManager.Create()
	log.Printf("Created new session: %s", sess.ID)

	// Create a new peer connection
	peerConnection, err := webrtc.NewPeerConnection(*h.config)
	if err != nil {
		return "", fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Create a local audio track for echo
	localTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		return "", fmt.Errorf("failed to create local track: %w", err)
	}

	// Add the local track to the peer connection
	rtpSender, err := peerConnection.AddTrack(localTrack)
	if err != nil {
		return "", fmt.Errorf("failed to add track: %w", err)
	}

	// Read RTCP packets (keep connection alive)
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Handle incoming tracks (echo logic)
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Session %s: Received track: %s (codec: %s)", sess.ID, track.ID(), track.Codec().MimeType)
		sess.UpdateState(session.StateListening)

		// Echo: read RTP packets and write them to the local track
		go func() {
			defer func() {
				sess.UpdateState(session.StateDisconnected)
				log.Printf("Session %s: Track ended", sess.ID)
			}()

			for {
				// Read RTP packet
				rtp, _, readErr := track.ReadRTP()
				if readErr != nil {
					if readErr == io.EOF {
						return
					}
					log.Printf("Session %s: Error reading RTP: %v", sess.ID, readErr)
					return
				}

				// Echo back: write the same packet to the local track
				if writeErr := localTrack.WriteRTP(rtp); writeErr != nil {
					if writeErr == io.ErrClosedPipe {
						return
					}
					log.Printf("Session %s: Error writing RTP: %v", sess.ID, writeErr)
				}
			}
		}()
	})

	// Handle connection state changes
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("Session %s: Connection state changed: %s", sess.ID, s.String())

		switch s {
		case webrtc.PeerConnectionStateConnected:
			sess.UpdateState(session.StateConnected)
		case webrtc.PeerConnectionStateDisconnected, webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateClosed:
			sess.UpdateState(session.StateDisconnected)
			h.sessionManager.Delete(sess.ID)
		}
	})

	// Set the remote description (offer)
	var offer webrtc.SessionDescription
	if err := json.Unmarshal([]byte(offerJSON), &offer); err != nil {
		return "", fmt.Errorf("failed to unmarshal offer: %w", err)
	}

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		return "", fmt.Errorf("failed to set remote description: %w", err)
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create answer: %w", err)
	}

	// Set the local description
	if err := peerConnection.SetLocalDescription(answer); err != nil {
		return "", fmt.Errorf("failed to set local description: %w", err)
	}

	// Marshal the answer to JSON
	answerJSON, err := json.Marshal(answer)
	if err != nil {
		return "", fmt.Errorf("failed to marshal answer: %w", err)
	}

	return string(answerJSON), nil
}
