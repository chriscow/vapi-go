// Package workflow provides types and logic for building conversational workflows and VAPI integration.
package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/chriscow/vapi-go"
)

// ProcessVAPIUpdate processes a VAPI conversation update and advances the workflow if needed.
//
// It takes the context, workflow engine, workflow ID, user ID, call ID, a slice of VAPI messages,
// the VAPI control URL, and a logger. It converts VAPI messages to a format the workflow engine can process,
// advances the workflow, and if the current node is a Say node, sends its message to the user via VAPI.
//
// Returns an error if processing or message sending fails.
func ProcessVAPIUpdate(
	ctx context.Context,
	engine *WorkflowEngine,
	workflowID string,
	userID string,
	callID string,
	messages []vapi.Message,
	controlURL string,
	logger *slog.Logger,
) error {
	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With(
		"workflowID", workflowID,
		"userID", userID,
		"callID", callID,
	)

	// Convert VAPI messages to a format the workflow engine can process
	// processedMessages := make([]map[string]any, len(messages))
	// for i, msg := range messages {
	// 	processedMessages[i] = map[string]any{
	// 		"content":          msg.Message,
	// 		"role":             msg.Role,
	// 		"secondsFromStart": msg.SecondsFromStart,
	// 	}
	// }

	// Process the message and advance the workflow
	_, err := engine.ProcessConversationUpdate(ctx, workflowID, userID, callID, messages)
	if err != nil {
		return fmt.Errorf("failed to process conversation update: %w", err)
	}

	// If the current node is a Say node, send its message to the user
	if len(messages) > 0 && messages[len(messages)-1].Role == "user" {
		// Only respond to user messages
		message, err := engine.GetCurrentNodeMessage(ctx, workflowID, userID, callID)
		if err != nil {
			logger.Error("failed to get current node message", "error", err)
			return err
		}

		if message != "" {
			logger.Info("sending message to user", "message", message)
			if err := sendVAPIMessage(ctx, controlURL, message, logger); err != nil {
				logger.Error("failed to send message to user", "error", err)
				return err
			}
		}
	}

	return nil
}

// sendVAPIMessage sends a message to the user via the VAPI control API.
//
// It takes the context, VAPI control URL, message to send, and a logger.
// It constructs the appropriate control message, sends it as a POST request to the control URL,
// and logs the result. Returns an error if the request fails or the response is not HTTP 200 OK.
func sendVAPIMessage(ctx context.Context, controlURL string, message string, logger *slog.Logger) error {
	// Prepare the control message
	controlMsg := map[string]any{
		"type": "add-message",
		"message": map[string]any{
			"role":    "system",
			"content": message,
		},
		"triggerResponseEnabled": true,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(controlMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal control message: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", controlURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	logger.Info("message sent successfully", "response", string(body))
	return nil
}
