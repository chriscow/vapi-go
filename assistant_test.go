package vapi

import (
	"context"
	"os"
	"testing"
)

func TestGetAssistant_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Ensure VAPI_API_KEY is set
	if os.Getenv("VAPI_API_KEY") == "" {
		t.Skip("VAPI_API_KEY not set")
	}

	ctx := context.Background()
	assistantID := "0f7534fa-deee-4feb-a48f-6e6e64eb38e7"

	// Make the API call
	got, err := GetAssistant(ctx, assistantID)
	if err != nil {
		t.Fatalf("GetAssistant() error = %v", err)
	}

	// Basic validation of the response
	if got == nil {
		t.Fatal("GetAssistant() returned nil result")
	}

	// Verify specific fields
	if got.Name == nil || *got.Name == "" {
		t.Error("Assistant name is missing or empty")
	} else {
		t.Logf("Assistant name: %s", *got.Name)
	}

	// Verify model configuration
	if got.Model == nil {
		t.Error("Model configuration is nil")
	} else {
		if got.Model.Provider == "" {
			t.Error("Model provider is empty")
		} else {
			t.Logf("Model provider: %s", got.Model.Provider)
		}
		if got.Model.Model == "" {
			t.Error("Model name is empty")
		} else {
			t.Logf("Model name: %s", got.Model.Model)
		}
	}

	// Verify voice configuration
	if got.Voice == nil {
		t.Error("Voice configuration is nil")
	} else {
		if got.Voice.Provider == "" {
			t.Error("Voice provider is empty")
		} else {
			t.Logf("Voice provider: %s", got.Voice.Provider)
		}
		if got.Voice.VoiceID == "" {
			t.Error("VoiceID is empty")
		} else {
			t.Logf("VoiceID: %s", got.Voice.VoiceID)
		}
	}

	// Verify first message
	if got.FirstMessage != nil {
		t.Logf("First message: %s", *got.FirstMessage)
	}

	// Verify end call message
	if got.EndCallMessage != nil {
		t.Logf("End call message: %s", *got.EndCallMessage)
	}

	// Log that the JSON was saved to a file for examination
	t.Logf("Raw JSON response saved to assistant-%s-response.json", assistantID)
}
