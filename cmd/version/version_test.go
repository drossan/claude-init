package version

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestNewVersionCommand_CreatesCommand verifica que se crea el comando correctamente.
func TestNewVersionCommand_CreatesCommand(t *testing.T) {
	cmd := NewVersionCommand()

	if cmd == nil {
		t.Fatal("NewVersionCommand() returned nil")
	}

	if cmd.Use != "version" {
		t.Errorf("expected Use 'version', got '%s'", cmd.Use)
	}

	if cmd.Short != "Show version information" {
		t.Errorf("expected Short 'Show version information', got '%s'", cmd.Short)
	}
}

// TestVersionCommand_DefaultOutput_ShowsVersion verifica que muestra la versión por defecto.
func TestVersionCommand_DefaultOutput_ShowsVersion(t *testing.T) {
	cmd := NewVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, Version) {
		t.Errorf("expected output to contain version '%s', got: %s", Version, output)
	}
	if !strings.Contains(output, "claude-init") {
		t.Errorf("expected output to contain 'claude-init', got: %s", output)
	}
}

// TestVersionCommand_WithShortFlag_ShowsOnlyNumber verifica que solo muestra el número de versión.
func TestVersionCommand_WithShortFlag_ShowsOnlyNumber(t *testing.T) {
	cmd := NewVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	cmd.SetArgs([]string{"--short"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != Version {
		t.Errorf("expected output to be exactly '%s', got: '%s'", Version, output)
	}

	// Verificar que no contiene texto adicional
	if strings.Contains(output, "claude-init") {
		t.Error("expected --short to only show version number")
	}
}

// TestVersionCommand_WithJSONFlag_ReturnsJSON verifica que la salida en JSON es correcta.
func TestVersionCommand_WithJSONFlag_ReturnsJSON(t *testing.T) {
	cmd := NewVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	cmd.SetArgs([]string{"--json"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()

	// Verificar que es JSON válido
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("expected valid JSON, got error: %v\nOutput: %s", err, output)
	}

	// Verificar campos esperados
	if _, ok := result["version"]; !ok {
		t.Error("expected JSON to contain 'version' field")
	}
	if result["version"] != Version {
		t.Errorf("expected version '%s', got '%v'", Version, result["version"])
	}
}

// TestVersionCommand_WithVerboseFlag_ShowsDetails verifica que muestra información detallada.
func TestVersionCommand_WithVerboseFlag_ShowsDetails(t *testing.T) {
	// Establecer valores de prueba
	oldCommit := Commit
	oldDate := BuildDate
	Commit = "abc123"
	BuildDate = "2024-01-15"
	defer func() {
		Commit = oldCommit
		BuildDate = oldDate
	}()

	cmd := NewVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	cmd.SetArgs([]string{"--verbose"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()

	// Verificar que muestra la versión
	if !strings.Contains(output, Version) {
		t.Errorf("expected output to contain version '%s', got: %s", Version, output)
	}

	// Verificar que muestra información adicional
	if !strings.Contains(output, Commit) {
		t.Errorf("expected verbose output to contain commit '%s', got: %s", Commit, output)
	}
	if !strings.Contains(output, BuildDate) {
		t.Errorf("expected verbose output to contain build date '%s', got: %s", BuildDate, output)
	}
}

// TestVersionCommand_JSONWithVerbose_IncludesAllFields verifica que JSON verbose incluye todos los campos.
func TestVersionCommand_JSONWithVerbose_IncludesAllFields(t *testing.T) {
	// Establecer valores de prueba
	oldCommit := Commit
	oldDate := BuildDate
	Commit = "def456"
	BuildDate = "2024-02-20"
	defer func() {
		Commit = oldCommit
		BuildDate = oldDate
	}()

	cmd := NewVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	cmd.SetArgs([]string{"--json", "--verbose"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	// Verificar todos los campos
	if _, ok := result["version"]; !ok {
		t.Error("expected JSON to contain 'version' field")
	}
	if _, ok := result["commit"]; !ok {
		t.Error("expected JSON to contain 'commit' field")
	}
	if _, ok := result["build_date"]; !ok {
		t.Error("expected JSON to contain 'build_date' field")
	}

	// Verificar valores
	if result["commit"] != Commit {
		t.Errorf("expected commit '%s', got '%v'", Commit, result["commit"])
	}
	if result["build_date"] != BuildDate {
		t.Errorf("expected build_date '%s', got '%v'", BuildDate, result["build_date"])
	}
}

// TestVersionCommand_WithShortAndJSON_ReturnsJSONVersionOnly verifica que short + json solo devuelve versión.
func TestVersionCommand_WithShortAndJSON_ReturnsJSONVersionOnly(t *testing.T) {
	cmd := NewVersionCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	cmd.SetArgs([]string{"--short", "--json"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	// Verificar que solo tiene el campo version
	if len(result) != 1 {
		t.Errorf("expected JSON with --short to only have 'version' field, got %d fields", len(result))
	}
	if _, ok := result["version"]; !ok {
		t.Error("expected JSON to contain 'version' field")
	}
}
