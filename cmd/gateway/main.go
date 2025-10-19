package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"voice-gateway/internal/config"
	"voice-gateway/internal/session"
	"voice-gateway/internal/webrtc"
)

type OfferRequest struct {
	SDP string `json:"sdp"`
}

type AnswerResponse struct {
	SDP string `json:"sdp"`
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Create session manager
	sessionMgr := session.NewManager()

	// Create WebRTC handler
	webrtcHandler := webrtc.NewHandler(cfg.WebRTC.ICEServers, sessionMgr)

	// Set up HTTP handlers
	http.HandleFunc("/offer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req OfferRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
			return
		}

		answer, err := webrtcHandler.HandleOffer(req.SDP)
		if err != nil {
			log.Printf("Error handling offer: %v", err)
			http.Error(w, fmt.Sprintf("Failed to process offer: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AnswerResponse{SDP: answer})
	})

	// Serve static files (web UI)
	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Voice Gateway starting on %s", addr)
	log.Printf("WebRTC echo server ready")
	log.Printf("Open http://%s in your browser to test", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
