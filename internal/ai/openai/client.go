package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client es un cliente para la API de OpenAI.
type Client struct {
	apiKey      string
	baseURL     string
	model       string
	maxTokens   int
	temperature float32
	client      *http.Client
}

// NewClient crea un nuevo cliente de OpenAI.
func NewClient(apiKey, baseURL, model string, maxTokens int) *Client {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "gpt-5.1"
	}
	if maxTokens == 0 {
		maxTokens = 100000 // GPT-5.1 soporta tokens extendidos
	}
	temperature := float32(0.7) // Balance entre creatividad y coherencia

	return &Client{
		apiKey:      apiKey,
		baseURL:     baseURL,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// chatRequest representa una solicitud a la API de chat de OpenAI.
type chatRequest struct {
	Model               string        `json:"model"`
	Messages            []chatMessage `json:"messages"`
	Temperature         float32       `json:"temperature,omitempty"`
	MaxTokens           int           `json:"max_tokens,omitempty"`
	MaxCompletionTokens int           `json:"max_completion_tokens,omitempty"`
}

// chatMessage representa un mensaje en la conversación.
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse representa la respuesta de la API.
type chatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// SendMessage envía un mensaje a OpenAI y retorna la respuesta.
func (c *Client) SendMessage(systemPrompt, userMessage string) (string, error) {
	messages := []chatMessage{}

	if systemPrompt != "" {
		messages = append(messages, chatMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	messages = append(messages, chatMessage{
		Role:    "user",
		Content: userMessage,
	})

	// GPT-5 models require max_completion_tokens instead of max_tokens
	reqBody := chatRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: c.temperature,
	}

	// Usar max_completion_tokens para modelos GPT-5, max_tokens para otros
	if isGPT5Model(c.model) {
		reqBody.MaxCompletionTokens = c.maxTokens
	} else {
		reqBody.MaxTokens = c.maxTokens
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
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

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *Client) SendSimpleMessage(message string) (string, error) {
	return c.SendMessage("", message)
}

// Close cierra el cliente y libera recursos.
func (c *Client) Close() error {
	// Nada que cerrar para el cliente HTTP básico
	return nil
}

// isGPT5Model detecta si el modelo es un modelo GPT-5 que requiere max_completion_tokens.
func isGPT5Model(model string) bool {
	// Los modelos GPT-5 requieren max_completion_tokens en lugar de max_tokens
	return contains(model, "gpt-5") && !contains(model, "gpt-4") && !contains(model, "gpt-3")
}

// contains verifica si una cadena contiene un substring (case-insensitive).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

// findSubstring busca un substring en una cadena de forma case-insensitive.
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
