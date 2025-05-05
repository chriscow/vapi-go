// Package workflow provides types and logic for building conversational workflows.
package workflow

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chriscow/vapi-go"
)

type MessageType string

const (
	// MessageTypeExact indicates the message is an exact string.
	MessageTypeExact MessageType = "exact"
	// MessageTypeGenerated indicates the message is generated from a prompt.
	MessageTypeGenerated MessageType = "generated"
)

// SayNode represents a node that outputs a message to the user.
// It can either output an exact message or generate a message from an LLM prompt.
type SayNode struct {
	BaseNode
	// Message is the exact message to output if MessageType is "exact".
	Message string
	// LLMPrompt is the prompt to use for generating a message if MessageType is "generated".
	LLMPrompt string
	// MessageType determines how the message is produced: "exact" or "generated".
	MessageType MessageType
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
		MessageType: MessageTypeExact,
	}
}

// NewSayNodeWithLLMPrompt creates a new SayNode that generates its message from an LLM prompt.
// id is the unique identifier for the node.
// prompt is the prompt to use for message generation.
func NewGeneratedSayNode(id string, prompt string) *SayNode {
	now := time.Now()
	return &SayNode{
		BaseNode: BaseNode{
			NodeID:        id,
			NodeType:      NodeTypeSay,
			CreatedAt:     now,
			LastUpdatedAt: now,
		},
		LLMPrompt:   prompt,
		MessageType: MessageTypeGenerated,
	}
}

// Execute runs the SayNode's action, outputting a message to the user or generating one from a prompt.
// It updates the workflow state with the message and marks the node as completed.
// If there is a next node, it updates the current node; otherwise, it marks the workflow as complete.
func (n *SayNode) Execute(ctx context.Context, state *WorkflowState, messages []vapi.Message) error {
	logger := slog.Default().With("node", n.NodeID, "type", n.NodeType)

	message := ""
	if n.MessageType == MessageTypeExact {
		message = n.Message
	} else if n.MessageType == MessageTypeGenerated {
		// Simulate LLM message generation
		// In a real implementation, we would call an LLM API here
		message = fmt.Sprintf("Generated message from prompt: %s", n.LLMPrompt)
	} else {
		return fmt.Errorf("invalid message type: %s", n.MessageType)
	}

	// In a real implementation, we would send this message to the user
	// For now, just log it
	logger.Info("executing say node", "message", message)

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

// // ToMap converts the SayNode to a map[string]any for storage or serialization.
// // The returned map contains all relevant fields of the node.
// func (n *SayNode) ToMap() map[string]any {
// 	return map[string]any{
// 		"id":            n.NodeID,
// 		"type":          string(n.NodeType),
// 		"nextNodeId":    n.NextNodeID,
// 		"message":       n.Message,
// 		"llmPrompt":     n.LLMPrompt,
// 		"messageType":   n.MessageType,
// 		"createdAt":     n.CreatedAt,
// 		"lastUpdatedAt": n.LastUpdatedAt,
// 	}
// }

// // FromMap initializes the SayNode from a map[string]any, typically loaded from storage.
// // It sets all relevant fields of the node from the map.
// func (n *SayNode) FromMap(data map[string]any) error {
// 	if id, ok := data["id"].(string); ok {
// 		n.NodeID = id
// 	}

// 	if typeStr, ok := data["type"].(string); ok {
// 		n.NodeType = NodeType(typeStr)
// 	}

// 	if nextNodeID, ok := data["nextNodeId"].(string); ok {
// 		n.NextNodeID = nextNodeID
// 	}

// 	if message, ok := data["message"].(string); ok {
// 		n.Message = message
// 	}

// 	if llmPrompt, ok := data["llmPrompt"].(string); ok {
// 		n.LLMPrompt = llmPrompt
// 	}

// 	if messageType, ok := data["messageType"].(string); ok {
// 		n.MessageType = messageType
// 	}

// 	if createdAt, ok := data["createdAt"].(time.Time); ok {
// 		n.CreatedAt = createdAt
// 	}

// 	if lastUpdatedAt, ok := data["lastUpdatedAt"].(time.Time); ok {
// 		n.LastUpdatedAt = lastUpdatedAt
// 	}

// 	return nil
// }
