package claudeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client es un cliente para la API de Anthropic Claude.
type Client struct {
	apiKey    string
	baseURL   string
	model     string
	maxTokens int
	client    *http.Client
}

// NewClient crea un nuevo cliente de Claude API.
func NewClient(apiKey, baseURL, model string, maxTokens int) *Client {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1/messages"
	}
	if model == "" {
		model = "claude-opus-4"
	}
	if maxTokens == 0 {
		maxTokens = 200000 // Claude Opus 4.5 tiene 200K tokens
	}

	return &Client{
		apiKey:    apiKey,
		baseURL:   baseURL,
		model:     model,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// messageRequest representa una solicitud a la API de mensajes de Claude.
type messageRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
	System    string    `json:"system,omitempty"`
	Stream    bool      `json:"stream,omitempty"`
}

// message representa un mensaje en la conversación.
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// messageResponse representa la respuesta de la API.
type messageResponse struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Role       string         `json:"role"`
	Content    []contentBlock `json:"content"`
	Model      string         `json:"model"`
	StopReason string         `json:"stop_reason"`
	Error      *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// contentBlock representa un bloque de contenido en la respuesta.
type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// SendMessage envía un mensaje a Claude y retorna la respuesta.
func (c *Client) SendMessage(systemPrompt, userMessage string) (string, error) {
	messages := []message{
		{
			Role:    "user",
			Content: userMessage,
		},
	}

	req := messageRequest{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		Messages:  messages,
		System:    systemPrompt,
		Stream:    false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var msgResp messageResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if msgResp.Error != nil {
		return "", fmt.Errorf("API error: %s - %s", msgResp.Error.Type, msgResp.Error.Message)
	}

	if len(msgResp.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return msgResp.Content[0].Text, nil
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *Client) SendSimpleMessage(message string) (string, error) {
	return c.SendMessage("", message)
}

// Close cierra el cliente y libera recursos.
func (c *Client) Close() error {
	return nil
}
