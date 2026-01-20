package ai

import (
	"testing"
)

// TestCLIClientIsAvailable verifies that CLIClient.IsAvailable works correctly
// when created via the factory.
// Regression test for nil pointer dereference bug.
func TestCLIClientIsAvailable(t *testing.T) {
	factory := NewClientFactory()

	client, err := factory.CreateClient(ProviderCLI)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	cliClient, ok := client.(*CLIClient)
	if !ok {
		t.Fatalf("Expected *CLIClient, got %T", client)
	}

	// This should not panic even if Claude CLI is not installed
	available, err := cliClient.IsAvailable()
	if err != nil {
		// It's OK if Claude CLI is not installed
		t.Logf("Claude CLI not available (expected in test env): %v", err)
	}
	t.Logf("CLIClient available: %v", available)
}

// TestCLIClientHasWrapper verifies that CLIClient created via factory
// has a non-nil wrapper field.
func TestCLIClientHasWrapper(t *testing.T) {
	factory := NewClientFactory()

	client, err := factory.CreateClient(ProviderCLI)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	cliClient, ok := client.(*CLIClient)
	if !ok {
		t.Fatalf("Expected *CLIClient, got %T", client)
	}

	if cliClient.wrapper == nil {
		t.Error("CLIClient.wrapper is nil, this will cause panic when calling methods")
	}
}
