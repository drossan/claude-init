package claude

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CLIWrapper es un wrapper para ejecutar Claude CLI.
type CLIWrapper struct {
	timeout time.Duration
}

// NewCLIWrapper crea un nuevo wrapper para Claude CLI.
func NewCLIWrapper() *CLIWrapper {
	return &CLIWrapper{
		timeout: 120 * time.Second, // 2 minutos por defecto
	}
}

// SetTimeout cambia el timeout para las ejecuciones.
func (c *CLIWrapper) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// CheckInstalled verifica si Claude CLI está instalado.
func (c *CLIWrapper) CheckInstalled() error {
	cmd := exec.Command("claude", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Claude CLI not found. Please install it from: https://claude.com/claude-code\n\nError: %w", err)
	}

	version := string(output)
	if !strings.Contains(strings.ToLower(version), "claude") {
		return fmt.Errorf("Claude CLI found but version check failed: %s", version)
	}

	return nil
}

// GetVersion retorna la versión de Claude CLI instalada.
func (c *CLIWrapper) GetVersion() (string, error) {
	cmd := exec.Command("claude", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get Claude CLI version: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// SendMessage envía un mensaje a Claude CLI y retorna la respuesta.
// systemPrompt se usa como contexto adicional y se pasa con --system-prompt.
func (c *CLIWrapper) SendMessage(systemPrompt, userMessage string) (string, error) {
	// Construir argumentos: -p para print mode
	args := []string{"-p"}

	// Agregar system prompt si existe
	if systemPrompt != "" {
		args = append(args, "--system-prompt", systemPrompt)
	}

	// Agregar el mensaje del usuario
	args = append(args, userMessage)

	cmd := exec.Command("claude", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Ejecutar
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("claude CLI error: %w\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *CLIWrapper) SendSimpleMessage(message string) (string, error) {
	return c.SendMessage("", message)
}

// SendMessageWithJSONSchema envía un mensaje esperando respuesta JSON validada.
func (c *CLIWrapper) SendMessageWithJSONSchema(systemPrompt, userMessage, jsonSchema string) (string, error) {
	args := []string{"-p", "--json-schema", jsonSchema}

	if systemPrompt != "" {
		args = append(args, "--system-prompt", systemPrompt)
	}

	args = append(args, userMessage)

	cmd := exec.Command("claude", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("claude CLI error: %w\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
