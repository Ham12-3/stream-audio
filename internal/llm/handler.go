package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the LLM
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// ChatResponse represents a response from the LLM
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int     `json:"index"`
		Message Message `json:"message"`
		Delta   struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// Handler manages LLM interactions
type Handler struct {
	apiURL     string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewHandler creates a new LLM handler
// Supports OpenAI-compatible APIs (OpenAI, Anthropic via proxy, local LLMs)
func NewHandler(apiURL, apiKey, model string) *Handler {
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1/chat/completions"
	}
	if model == "" {
		model = "gpt-4o-mini"
	}

	return &Handler{
		apiURL: apiURL,
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Chat sends a chat request to the LLM
func (h *Handler) Chat(messages []Message) (string, error) {
	req := ChatRequest{
		Model:    h.model,
		Messages: messages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", h.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if h.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)
	}

	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// StreamChat sends a streaming chat request to the LLM
func (h *Handler) StreamChat(messages []Message, onChunk func(string)) error {
	req := ChatRequest{
		Model:    h.model,
		Messages: messages,
		Stream:   true,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", h.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if h.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)
	}

	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse SSE stream
	reader := io.Reader(resp.Body)
	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read stream: %w", err)
		}

		chunk := string(buf[:n])
		lines := strings.Split(chunk, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return nil
			}

			var chatResp ChatResponse
			if err := json.Unmarshal([]byte(data), &chatResp); err != nil {
				continue // Skip malformed chunks
			}

			if len(chatResp.Choices) > 0 && chatResp.Choices[0].Delta.Content != "" {
				onChunk(chatResp.Choices[0].Delta.Content)
			}
		}
	}

	return nil
}

// ConversationContext maintains conversation state
type ConversationContext struct {
	SessionID string
	Messages  []Message
	MaxTokens int
}

// NewConversationContext creates a new conversation context
func NewConversationContext(sessionID string, systemPrompt string) *ConversationContext {
	messages := []Message{}
	if systemPrompt != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	return &ConversationContext{
		SessionID: sessionID,
		Messages:  messages,
		MaxTokens: 4000, // Conservative limit for context
	}
}

// AddUserMessage adds a user message to the context
func (c *ConversationContext) AddUserMessage(content string) {
	c.Messages = append(c.Messages, Message{
		Role:    "user",
		Content: content,
	})
}

// AddAssistantMessage adds an assistant message to the context
func (c *ConversationContext) AddAssistantMessage(content string) {
	c.Messages = append(c.Messages, Message{
		Role:    "assistant",
		Content: content,
	})
}

// GetMessages returns all messages
func (c *ConversationContext) GetMessages() []Message {
	return c.Messages
}

// Clear clears the conversation (keeping system prompt)
func (c *ConversationContext) Clear() {
	if len(c.Messages) > 0 && c.Messages[0].Role == "system" {
		c.Messages = c.Messages[:1]
	} else {
		c.Messages = []Message{}
	}
}
