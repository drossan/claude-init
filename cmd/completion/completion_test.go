package completion

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestNewCompletionCommand_CreatesCommand verifica que se crea el comando correctamente.
func TestNewCompletionCommand_CreatesCommand(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	cmd := NewCompletionCommand(rootCmd)

	if cmd == nil {
		t.Fatal("NewCompletionCommand() returned nil")
	}

	if cmd.Use != "completion [bash|zsh|fish|powershell]" {
		t.Errorf("expected Use 'completion [bash|zsh|fish|powershell]', got '%s'", cmd.Use)
	}

	if cmd.Short != "Generate completion script" {
		t.Errorf("expected Short 'Generate completion script', got '%s'", cmd.Short)
	}
}

// TestCompletionCommand_WithBash_GeneratesBashScript verifica que genera script de bash.
func TestCompletionCommand_WithBash_GeneratesBashScript(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd := &cobra.Command{Use: "ia-start"}
	rootCmd.SetOut(buf)
	cmd := NewCompletionCommand(rootCmd)

	cmd.SetArgs([]string{"bash"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()

	// Verificar que contiene elementos típicos de un script de completación de bash
	if !strings.Contains(output, "complete") {
		t.Error("expected bash completion script to contain 'complete'")
	}
	// Los scripts de bash completión suelen tener funciones con __start_
	if !strings.Contains(output, "__start_") {
		t.Error("expected bash completion script to contain __start_ functions")
	}
}

// TestCompletionCommand_WithZsh_GeneratesZshScript verifica que genera script de zsh.
func TestCompletionCommand_WithZsh_GeneratesZshScript(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd := &cobra.Command{Use: "ia-start"}
	rootCmd.SetOut(buf)
	cmd := NewCompletionCommand(rootCmd)

	cmd.SetArgs([]string{"zsh"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()

	// Verificar que contiene elementos típicos de un script de completación de zsh
	if !strings.Contains(output, "compdef") {
		t.Error("expected zsh completion script to contain 'compdef'")
	}
}

// TestCompletionCommand_WithFish_GeneratesFishScript verifica que genera script de fish.
func TestCompletionCommand_WithFish_GeneratesFishScript(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd := &cobra.Command{Use: "ia-start"}
	rootCmd.SetOut(buf)
	cmd := NewCompletionCommand(rootCmd)

	cmd.SetArgs([]string{"fish"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()

	// Verificar que contiene elementos típicos de un script de completación de fish
	if !strings.Contains(output, "complete") {
		t.Error("expected fish completion script to contain 'complete'")
	}
}

// TestCompletionCommand_WithPowerShell_GeneratesPowerShellScript verifica que genera script de PowerShell.
func TestCompletionCommand_WithPowerShell_GeneratesPowerShellScript(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd := &cobra.Command{Use: "ia-start"}
	rootCmd.SetOut(buf)
	cmd := NewCompletionCommand(rootCmd)

	cmd.SetArgs([]string{"powershell"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()

	// Verificar que contiene elementos típicos de un script de PowerShell
	if !strings.Contains(output, "ia-start") {
		t.Error("expected powershell completion script to contain command name")
	}
}

// TestCompletionCommand_WithInvalidShell_ReturnsError verifica que retorna error con shell inválido.
func TestCompletionCommand_WithInvalidShell_ReturnsError(t *testing.T) {
	rootCmd := &cobra.Command{Use: "ia-start"}
	cmd := NewCompletionCommand(rootCmd)
	buf := new(bytes.Buffer)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"invalid_shell"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error with invalid shell, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported shell") && !strings.Contains(err.Error(), "invalid") {
		t.Errorf("expected error message to mention unsupported shell, got: %v", err)
	}
}

// TestCompletionCommand_NoArgs_ReturnsError verifica que retorna error sin argumentos.
func TestCompletionCommand_NoArgs_ReturnsError(t *testing.T) {
	rootCmd := &cobra.Command{Use: "ia-start"}
	cmd := NewCompletionCommand(rootCmd)
	buf := new(bytes.Buffer)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without arguments, got nil")
	}
}

// TestCompletionCommand_CaseInsensitive verifica que acepta shells en mayúsculas/minúsculas.
func TestCompletionCommand_CaseInsensitive(t *testing.T) {
	// Probar con mayúsculas - la implementación usa ToLower, así que debería funcionar
	testCases := []string{"BASH", "Bash", "BaSh"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd := &cobra.Command{Use: "ia-start"}
			rootCmd.SetOut(buf)

			// Necesitamos crear un nuevo comando para cada test porque cobra no resetea los args
			testCmd := NewCompletionCommand(rootCmd)

			testCmd.SetArgs([]string{tc})
			err := testCmd.Execute()
			if err != nil {
				t.Fatalf("expected shell name '%s' to work (case-insensitive), got: %v", tc, err)
			}

			output := buf.String()
			if !strings.Contains(output, "completion") && !strings.Contains(output, "complete") && !strings.Contains(output, "compdef") {
				t.Error("expected completion script to be generated")
			}
		})
	}
}
