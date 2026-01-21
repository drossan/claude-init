package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client es un cliente para la API de Google Gemini.
type Client struct {
	apiKey    string
	baseURL   string
	model     string
	maxTokens int
	client    *http.Client
}

// NewClient crea un nuevo cliente de Gemini.
func NewClient(apiKey, baseURL, model string, maxTokens int) *Client {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta/models"
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	if maxTokens == 0 {
		maxTokens = 1000000 // Gemini 2.5 Flash tiene 1M tokens
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

// generateContentRequest representa una solicitud a la API de Gemini.
type generateContentRequest struct {
	Contents          []content         `json:"contents"`
	SystemInstruction *content          `json:"systemInstruction,omitempty"`
	GenerationConfig  *generationConfig `json:"generationConfig,omitempty"`
}

type content struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text,omitempty"`
}

type generationConfig struct {
	MaxOutputTokens int `json:"maxOutputTokens,omitempty"`
}

// generateContentResponse representa la respuesta de la API de Gemini.
type generateContentResponse struct {
	Candidates []candidate `json:"candidates"`
	Error      *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

type candidate struct {
	Content      content `json:"content"`
	FinishReason string  `json:"finishReason,omitempty"`
}

// SendMessage envía un mensaje a Gemini y retorna la respuesta.
func (c *Client) SendMessage(systemPrompt, userMessage string) (string, error) {
	// Construir contents array
	contents := []content{}

	// En Gemini, el system prompt va en systemInstruction
	var systemInstruction *content
	if systemPrompt != "" {
		systemInstruction = &content{
			Parts: []part{{Text: systemPrompt}},
		}
	}

	// Añadir mensaje del usuario
	contents = append(contents, content{
		Role:  "user",
		Parts: []part{{Text: userMessage}},
	})

	req := generateContentRequest{
		Contents:          contents,
		SystemInstruction: systemInstruction,
		GenerationConfig: &generationConfig{
			MaxOutputTokens: c.maxTokens,
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	// El endpoint incluye el modelo: {model}:generateContent
	url := fmt.Sprintf("%s/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

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

	var geminiResp generateContentResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	candidate := geminiResp.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in response")
	}

	return candidate.Content.Parts[0].Text, nil
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *Client) SendSimpleMessage(message string) (string, error) {
	return c.SendMessage("", message)
}

// Close cierra el cliente y libera recursos.
func (c *Client) Close() error {
	return nil
}
