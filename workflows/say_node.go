package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// SayNode represents a node that outputs a message to the user
type SayNode struct {
	BaseNode
	Message     string `json:"message,omitempty"`
	LLMPrompt   string `json:"llmPrompt,omitempty"`
	MessageType string `json:"messageType,omitempty"` // "exact" or "generated"
}

// NewSayNode creates a new SayNode with an exact message
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

// NewSayNodeWithLLMPrompt creates a new SayNode with an LLM prompt
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

// Execute runs the Say node's action
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

// ToMap converts the SayNode to a map for storage
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

// FromMap initializes the SayNode from a map
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
