package claude

import (
	"time"
)

// HybridClient usa el CLIWrapper directo.
//
// NOTA: El cliente persistente está deshabilitado porque Claude CLI
// en modo interactivo no es determinista y puede colgarse esperando input.
// El CLIWrapper con -p es más lento pero 100% fiable.
type HybridClient struct {
	wrapper *CLIWrapper
	// Persistent client deshabilitado por ahora
	// persistent *PersistentClient
}

// NewHybridClient crea un nuevo cliente híbrido (usa solo CLIWrapper).
func NewHybridClient() *HybridClient {
	return &HybridClient{
		wrapper: NewCLIWrapper(),
	}
}

// SendMessage envía un mensaje usando CLIWrapper.
func (h *HybridClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	return h.wrapper.SendMessage(systemPrompt, userMessage)
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (h *HybridClient) SendSimpleMessage(message string) (string, error) {
	return h.wrapper.SendSimpleMessage(message)
}

// CheckInstalled verifica si Claude CLI está instalado.
func (h *HybridClient) CheckInstalled() error {
	return h.wrapper.CheckInstalled()
}

// GetVersion retorna la versión de Claude CLI.
func (h *HybridClient) GetVersion() (string, error) {
	return h.wrapper.GetVersion()
}

// Stop no hace nada (solo para compatibilidad con interfaz).
func (h *HybridClient) Stop() error {
	return nil
}

// IsRunning siempre retorna false (no hay proceso persistente).
func (h *HybridClient) IsRunning() bool {
	return false
}

// SetIdleTimeout no hace nada (solo para compatibilidad).
func (h *HybridClient) SetIdleTimeout(timeout time.Duration) {
	// No-op: persistent client deshabilitado
}
