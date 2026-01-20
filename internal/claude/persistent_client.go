package claude

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

// PersistentClient mantiene un proceso de Claude CLI corriendo
// y reutiliza la conexión para múltiples requests.
type PersistentClient struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	scanner *bufio.Scanner
	mu      sync.Mutex
	running bool
	timeout time.Duration
}

// NewPersistentClient crea un nuevo cliente persistente.
func NewPersistentClient() *PersistentClient {
	return &PersistentClient{
		timeout: 120 * time.Second,
	}
}

// Start inicia el proceso persistente de Claude CLI.
func (c *PersistentClient) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil // Ya está corriendo
	}

	// Iniciar Claude CLI sin -p para modo interactivo
	// Usamos --permission-mode=dontAsk para evitar prompts interactivos
	c.cmd = exec.Command("claude", "--permission-mode=dontAsk")

	// Crear pipes para stdin/stdout/stderr
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %w", err)
	}

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %w", err)
	}

	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %w", err)
	}

	c.stdin = stdin
	c.stdout = stdout
	c.stderr = stderr
	c.scanner = bufio.NewScanner(stdout)

	// Iniciar el proceso
	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("error starting claude CLI: %w", err)
	}

	c.running = true

	// Esperar a que el prompt esté listo
	// Leemos hasta ver el prompt de Claude (normalmente ">")
	time.Sleep(500 * time.Millisecond)

	return nil
}

// Stop detiene el proceso persistente.
func (c *PersistentClient) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return nil
	}

	// Cerrar stdin para señalar EOF al proceso
	if c.stdin != nil {
		c.stdin.Close()
	}

	// Esperar a que el proceso termine
	done := make(chan error, 1)
	go func() {
		done <- c.cmd.Wait()
	}()

	select {
	case err := <-done:
		c.running = false
		return err
	case <-time.After(5 * time.Second):
		// Timeout, matar el proceso
		c.cmd.Process.Kill()
		c.running = false
		return fmt.Errorf("timeout waiting for process to exit")
	}
}

// IsRunning retorna true si el proceso está corriendo.
func (c *PersistentClient) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

// SendMessage envía un prompt a Claude y retorna la respuesta.
// Reutiliza el proceso existente si está corriendo.
func (c *PersistentClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	// Asegurar que el proceso está corriendo
	if !c.IsRunning() {
		if err := c.Start(); err != nil {
			return "", err
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Construir el prompt completo
	prompt := userMessage
	if systemPrompt != "" {
		prompt = fmt.Sprintf("Context: %s\n\n%s", systemPrompt, userMessage)
	}

	// Escribir el prompt en stdin
	// Agregamos newline para enviar el comando
	if _, err := fmt.Fprintln(c.stdin, prompt); err != nil {
		return "", fmt.Errorf("error writing to stdin: %w", err)
	}

	// Leer la respuesta hasta ver un delimitador o timeout
	response, err := c.readResponse()
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return response, nil
}

// readResponse lee la respuesta de stdout hasta encontrar el siguiente prompt.
func (c *PersistentClient) readResponse() (string, error) {
	var response bytes.Buffer
	timeout := time.After(c.timeout)

	// Leer línea por línea
	for {
		select {
		case <-timeout:
			return response.String(), nil // Retornar lo que tenemos hasta ahora
		default:
			if !c.scanner.Scan() {
				if err := c.scanner.Err(); err != nil {
					return "", fmt.Errorf("error reading stdout: %w", err)
				}
				return response.String(), nil // EOF
			}

			line := c.scanner.Text()
			response.WriteString(line)
			response.WriteString("\n")

			// Si vemos un prompt de Claude (>), terminamos
			if line == ">" || len(line) > 0 && line[len(line)-1] == '>' {
				// Remover el prompt de la respuesta
				respStr := response.String()
				respStr = respStr[:len(respStr)-len(line)-1]
				return respStr, nil
			}
		}
	}
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *PersistentClient) SendSimpleMessage(message string) (string, error) {
	return c.SendMessage("", message)
}

// Restart reinicia el proceso persistente.
func (c *PersistentClient) Restart() error {
	if err := c.Stop(); err != nil {
		return fmt.Errorf("error stopping: %w", err)
	}
	return c.Start()
}

// CheckInstalled verifica si Claude CLI está instalado (método para implementar interfaz Client).
func (c *PersistentClient) CheckInstalled() error {
	wrapper := NewCLIWrapper()
	return wrapper.CheckInstalled()
}

// GetVersion retorna la versión de Claude CLI (método para implementar interfaz Client).
func (c *PersistentClient) GetVersion() (string, error) {
	wrapper := NewCLIWrapper()
	return wrapper.GetVersion()
}
