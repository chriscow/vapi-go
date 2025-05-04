// Package workflow provides types and logic for building conversational workflows.
package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// SayNode represents a node that outputs a message to the user.
// It can either output an exact message or generate a message from an LLM prompt.
type SayNode struct {
	BaseNode
	// Message is the exact message to output if MessageType is "exact".
	Message string `json:"message,omitempty"`
	// LLMPrompt is the prompt to use for generating a message if MessageType is "generated".
	LLMPrompt string `json:"llmPrompt,omitempty"`
	// MessageType determines how the message is produced: "exact" or "generated".
	MessageType string `json:"messageType,omitempty"` // "exact" or "generated"
}

// NewSayNode creates a new SayNode with an exact message.
// id is the unique identifier for the node.
// message is the message to output to the user.
func NewSayNode(id string, message string) *SayNode {
	now := time.Now()
	return &SayNode{
		BaseNode: BaseNode{
			NodeID:        id,
			NodeType:      NodeTypeSay,
			CreatedAt:     now,
			LastUpdatedAt: now,
		},
		Message:     message,
		MessageType: "exact",
	}
}

// NewSayNodeWithLLMPrompt creates a new SayNode that generates its message from an LLM prompt.
// id is the unique identifier for the node.
// prompt is the prompt to use for message generation.
func NewSayNodeWithLLMPrompt(id string, prompt string) *SayNode {
	now := time.Now()
	return &SayNode{
		BaseNode: BaseNode{
			NodeID:        id,
			NodeType:      NodeTypeSay,
			CreatedAt:     now,
			LastUpdatedAt: now,
		},
		LLMPrompt:   prompt,
		MessageType: "generated",
	}
}

// Execute runs the SayNode's action, outputting a message to the user or generating one from a prompt.
// It updates the workflow state with the message and marks the node as completed.
// If there is a next node, it updates the current node; otherwise, it marks the workflow as complete.
func (n *SayNode) Execute(ctx context.Context, state *WorkflowState) error {
	logger := slog.Default().With("node", n.NodeID, "type", n.NodeType)

	message := ""
	if n.MessageType == "exact" {
		message = n.Message
	} else if n.MessageType == "generated" {
		// For MVP, we'll use a simple approach to generate a message
		// In a real implementation, we would use an LLM API here

		// Replace variables in the prompt
		prompt := n.LLMPrompt
		for key, value := range state.Variables {
			valueStr, _ := json.Marshal(value)
			prompt = fmt.Sprintf("%s\n%s: %s", prompt, key, string(valueStr))
		}

		// For now, just use the prompt as the message
		// In a real implementation, we would call the LLM here
		message = fmt.Sprintf("Generated message based on: %s", prompt)

		logger.Info("generated message from LLM prompt")
	}

	// In a real implementation, we would send this message to the user
	// For now, just log it
	logger.Info("executing say node", "message", message)

	// Add the message to state variables
	if state.Variables == nil {
		state.Variables = make(map[string]any)
	}
	state.Variables["lastSayMessage"] = message

	// Mark this node as completed
	state.CompletedNodeIDs = append(state.CompletedNodeIDs, n.NodeID)

	// Update the current node to the next node
	if n.NextNodeID != "" {
		state.CurrentNodeID = n.NextNodeID
	} else {
		// If there's no next node, mark the workflow as complete
		state.IsComplete = true
	}

	state.LastUpdatedAt = time.Now()

	return nil
}

// ToMap converts the SayNode to a map[string]any for storage or serialization.
// The returned map contains all relevant fields of the node.
func (n *SayNode) ToMap() map[string]any {
	return map[string]any{
		"id":            n.NodeID,
		"type":          string(n.NodeType),
		"nextNodeId":    n.NextNodeID,
		"message":       n.Message,
		"llmPrompt":     n.LLMPrompt,
		"messageType":   n.MessageType,
		"createdAt":     n.CreatedAt,
		"lastUpdatedAt": n.LastUpdatedAt,
	}
}

// FromMap initializes the SayNode from a map[string]any, typically loaded from storage.
// It sets all relevant fields of the node from the map.
func (n *SayNode) FromMap(data map[string]any) error {
	if id, ok := data["id"].(string); ok {
		n.NodeID = id
	}

	if typeStr, ok := data["type"].(string); ok {
		n.NodeType = NodeType(typeStr)
	}

	if nextNodeID, ok := data["nextNodeId"].(string); ok {
		n.NextNodeID = nextNodeID
	}

	if message, ok := data["message"].(string); ok {
		n.Message = message
	}

	if llmPrompt, ok := data["llmPrompt"].(string); ok {
		n.LLMPrompt = llmPrompt
	}

	if messageType, ok := data["messageType"].(string); ok {
		n.MessageType = messageType
	}

	if createdAt, ok := data["createdAt"].(time.Time); ok {
		n.CreatedAt = createdAt
	}

	if lastUpdatedAt, ok := data["lastUpdatedAt"].(time.Time); ok {
		n.LastUpdatedAt = lastUpdatedAt
	}

	return nil
}
