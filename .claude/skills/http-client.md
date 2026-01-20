# Skill: HTTP Client (http-client)

## Propósito
Especialidad en crear clientes HTTP robustos en Go que se comuniquen con APIs externas (Claude, OpenAI, z.ai).

## Responsabilidades
- Crear clientes HTTP que sigan las mejores prácticas
- Manejar timeouts y cancelación con context
- Manejar retries con backoff exponencial
- Manejar rate limiting apropiadamente
- Parsear respuestas JSON correctamente

## Cliente HTTP Base

### Estructura del Cliente

```go
package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
    return &Client{
        apiKey:  apiKey,
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}
```

### POST con JSON

```go
func (c *Client) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
    // Construir request
    body, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    // Crear request con context
    url := fmt.Sprintf("%s/v1/recommendations", c.baseURL)
    httpRequest, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    httpRequest.Header.Set("Content-Type", "application/json")
    httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
    httpRequest.Header.Set("X-API-Version", "2024-01-01")

    // Ejecutar request
    httpResponse, err := c.httpClient.Do(httpRequest)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer httpResponse.Body.Close()

    // Leer respuesta
    responseBody, err := io.ReadAll(httpResponse.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Verificar status code
    if httpResponse.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d: %s", httpResponse.StatusCode, string(responseBody))
    }

    // Parsear respuesta
    var response RecommendationResponse
    if err := json.Unmarshal(responseBody, &response); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return &response, nil
}
```

## Manejo de Errores HTTP

### Errores con Contexto

```go
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
    resp, err := c.httpClient.Do(req)
    if err != nil {
        // Error de red o timeout
        if ctx.Err() == context.DeadlineExceeded {
            return nil, fmt.Errorf("request timeout: %w", err)
        }
        return nil, fmt.Errorf("request failed: %w", err)
    }

    // Errores HTTP
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body) // best effort

        // Status codes específicos
        switch resp.StatusCode {
        case http.StatusUnauthorized:
            return nil, fmt.Errorf("unauthorized: invalid API key")
        case http.StatusTooManyRequests:
            return nil, fmt.Errorf("rate limit exceeded: %s", string(body))
        case http.StatusInternalServerError:
            return nil, fmt.Errorf("internal server error: %s", string(body))
        default:
            return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
        }
    }

    return resp, nil
}
```

### Retry con Backoff Exponencial

```go
package retry

import (
    "context"
    "fmt"
    "time"
)

type RetryConfig struct {
    MaxRetries    int
    InitialDelay  time.Duration
    MaxDelay      time.Duration
    BackoffFactor float64
}

func DoWithRetry(ctx context.Context, config RetryConfig, fn func() error) error {
    var lastErr error
    delay := config.InitialDelay

    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        if attempt > 0 {
            // Esperar antes de reintentar
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return ctx.Err()
            }
        }

        err := fn()
        if err == nil {
            return nil
        }

        lastErr = err

        // No reintentar para ciertos errores
        if !shouldRetry(err) {
            return err
        }

        // Calcular próximo delay con backoff exponencial
        delay = time.Duration(float64(delay) * config.BackoffFactor)
        if delay > config.MaxDelay {
            delay = config.MaxDelay
        }
    }

    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func shouldRetry(err error) bool {
    // Reintentar solo para errores transitorios
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        return true
    }
    // Añadir lógica específica según el error
    return false
}
```

## Streaming Responses

### Lectura de Streaming (Server-Sent Events)

```go
func (c *Client) StreamChat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, <-chan error) {
    chunks := make(chan ChatChunk)
    errs := make(chan error, 1)

    go func() {
        defer close(chunks)
        defer close(errs)

        // Crear request
        body, _ := json.Marshal(req)
        httpRequest, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/stream", bytes.NewReader(body))
        httpRequest.Header.Set("Content-Type", "application/json")
        httpRequest.Header.Set("Authorization", "Bearer "+c.apiKey)

        // Ejecutar request
        resp, err := c.httpClient.Do(httpRequest)
        if err != nil {
            errs <- err
            return
        }
        defer resp.Body.Close()

        // Leer stream línea por línea
        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := scanner.Text()

            // Parsear línea (formato SSE: data: {...})
            if strings.HasPrefix(line, "data: ") {
                jsonStr := strings.TrimPrefix(line, "data: ")

                var chunk ChatChunk
                if err := json.Unmarshal([]byte(jsonStr), &chunk); err != nil {
                    errs <- fmt.Errorf("failed to parse chunk: %w", err)
                    return
                }

                select {
                case chunks <- chunk:
                case <-ctx.Done():
                    errs <- ctx.Err()
                    return
                }
            }
        }

        if err := scanner.Err(); err != nil {
            errs <- err
        }
    }()

    return chunks, errs
}
```

## Rate Limiting

### Rate Limiter Simple

```go
package ratelimit

import (
    "context"
    "time"
)

type RateLimiter struct {
    tokens chan struct{}
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, requestsPerSecond),
    }

    // Llenar tokens inicialmente
    for i := 0; i < requestsPerSecond; i++ {
        rl.tokens <- struct{}{}
    }

    // Refill tokens
    go func() {
        ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
        defer ticker.Stop()
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}{}:
            default: // channel lleno, descartar token
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Integrar Rate Limiter en el Cliente

```go
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
    // Esperar si es necesario por rate limit
    if c.rateLimiter != nil {
        if err := c.rateLimiter.Wait(ctx); err != nil {
            return nil, err
        }
    }

    return c.httpClient.Do(req)
}
```

## Testing de Clientes HTTP

### Mock del HTTP Client

```go
package ai_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestClient_GenerateRecommendations(t *testing.T) {
    // Crear servidor de prueba
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verificar request
        if r.Method != "POST" {
            t.Errorf("expected POST, got %s", r.Method)
        }
        if r.Header.Get("Authorization") != "Bearer test-key" {
            t.Errorf("invalid authorization header")
        }

        // Responder
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"recommendations":["typescript","react"]}`))
    }))
    defer server.Close()

    // Crear cliente con URL del servidor de prueba
    client := NewClient("test-key", server.URL)

    // Ejecutar
    resp, err := client.GenerateRecommendations(context.Background(), RecommendationRequest{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Verificar respuesta
    if len(resp.Recommendations) != 2 {
        t.Errorf("expected 2 recommendations, got %d", len(resp.Recommendations))
    }
}
```

## Checklist de HTTP Client

- [ ] Los requests usan `context.Context`
- [ ] Los responses se cierran con `defer`
- [ ] Los timeouts están configurados
- [ ] Los errores tienen contexto suficiente
- [ ] El rate limiting está implementado si es necesario
- [ ] Los retries están implementados para errores transitorios
- [ ] Los headers están correctamente configurados
- [ ] Las respuestas se validan (status code)
- [ ] El JSON parsing maneja errores
- [ ] Los tests usan mocks (httptest.Server)
