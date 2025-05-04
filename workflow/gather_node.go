// Package workflow provides types and logic for building conversational workflows.
package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/chriscow/minds"
)

// DataType constants for variable types
type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeInteger DataType = "integer"
	DataTypeBoolean DataType = "boolean"
	DataTypeObject  DataType = "object"
	DataTypeArray   DataType = "array"
	DataTypeEnum    DataType = "enum"
)

// GatherNode represents a node that collects input from the user.
// It specifies which variables to gather, how many attempts to allow, and the prompt to use for LLM extraction.
type GatherNode struct {
	BaseNode
	GatherSchema   *minds.Definition
	MaxAttempts    int
	LLMPrompt      string
	FallbackNodeID string
	FollowUpPrompt string
	ExtractedData  map[string]any // JSON data extracted from conversation
}

// NewGatherNode creates a new GatherNode.
// id is the unique identifier for the node.
// schema is the list of variables to gather from the user.
// maxAttempts is the maximum number of attempts to gather input.
// llmPrompt is the prompt to use for LLM-based extraction.
func NewGatherNode(id string, schema *minds.Definition, maxAttempts int, llmPrompt string) *GatherNode {
	now := time.Now()
	return &GatherNode{
		BaseNode: BaseNode{
			NodeID:        id,
			NodeType:      NodeTypeGather,
			CreatedAt:     now,
			LastUpdatedAt: now,
		},
		GatherSchema:  schema,
		MaxAttempts:   maxAttempts,
		LLMPrompt:     llmPrompt,
		ExtractedData: make(map[string]any),
	}
}

// Execute runs the GatherNode's action, collecting input from the user or simulating extraction for MVP.
// It updates the workflow state with the extracted data and marks the node as completed.
// If there is a next node, it updates the current node; otherwise, it marks the workflow as complete.
func (n *GatherNode) Execute(ctx context.Context, state *WorkflowState) error {
	logger := slog.Default().With("node", n.NodeID, "type", n.NodeType)

	// For MVP, we'll use a simple approach to gather data
	// In a real implementation, we would use an LLM API here

	// Find missing properties from the schema
	missing := make([]string, 0)

	// Initialize the extracted data map if not already present
	if n.ExtractedData == nil {
		n.ExtractedData = make(map[string]any)
	}

	// Check which properties from the schema are missing
	if n.GatherSchema.Type == minds.Object && n.GatherSchema.Properties != nil {
		for propName, _ := range n.GatherSchema.Properties {
			// Check if the property is required and missing
			isRequired := false
			for _, req := range n.GatherSchema.Required {
				if req == propName {
					isRequired = true
					break
				}
			}

			if isRequired {
				if _, ok := n.ExtractedData[propName]; !ok {
					// In a real implementation, we would check if the variable is in state.Variables
					// from previous messages. For MVP, we'll just mark it as missing.
					missing = append(missing, propName)
				}
			}
		}
	}

	// If there are missing properties, generate a prompt to extract them
	if len(missing) > 0 {
		// Generate prompt for the LLM to extract missing data
		prompt := n.LLMPrompt + "\n\nPlease extract the following information:\n"

		for _, propName := range missing {
			propDef := n.GatherSchema.Properties[propName]
			prompt += fmt.Sprintf("- %s: %s (Type: %s, Required: true)\n",
				propName, propDef.Description, propDef.Type)
		}

		// For MVP, we'll simulate data extraction
		// In a real implementation, we would call the LLM API here
		logger.Info("would call LLM to extract data", "prompt", prompt)

		// For now, just add some dummy data
		// In a real implementation, this would come from the LLM's response
		for _, propName := range missing {
			propDef := n.GatherSchema.Properties[propName]

			switch propDef.Type {
			case minds.String:
				n.ExtractedData[propName] = fmt.Sprintf("Sample %s", propName)
			case minds.Number, minds.Integer:
				n.ExtractedData[propName] = 42
			case minds.Boolean:
				n.ExtractedData[propName] = true
			default:
				// Handle any other types
				if len(propDef.Enum) > 0 {
					n.ExtractedData[propName] = propDef.Enum[0]
				} else {
					n.ExtractedData[propName] = "unknown"
				}
			}
		}

		// Add the extracted data to state.Variables
		if state.Variables == nil {
			state.Variables = make(map[string]any)
		}

		for k, v := range n.ExtractedData {
			state.Variables[k] = v
		}

		// Log the extracted data
		dataJSON, _ := json.Marshal(n.ExtractedData)
		logger.Info("extracted data", "data", string(dataJSON))
	}

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

// ToMap converts the GatherNode to a map[string]any for storage or serialization.
// The returned map contains all relevant fields of the node, including schema and extracted data.
func (n *GatherNode) ToMap() map[string]any {
	schemaJSON, _ := json.Marshal(n.GatherSchema)

	return map[string]any{
		"id":             n.NodeID,
		"type":           string(n.NodeType),
		"nextNodeId":     n.NextNodeID,
		"gatherSchema":   json.RawMessage(schemaJSON),
		"maxAttempts":    n.MaxAttempts,
		"llmPrompt":      n.LLMPrompt,
		"fallbackNodeId": n.FallbackNodeID,
		"followUpPrompt": n.FollowUpPrompt,
		"extractedData":  n.ExtractedData,
		"createdAt":      n.CreatedAt,
		"lastUpdatedAt":  n.LastUpdatedAt,
	}
}

// FromMap initializes the GatherNode from a map[string]any, typically loaded from storage.
// It sets all relevant fields of the node from the map, including schema and extracted data.
func (n *GatherNode) FromMap(data map[string]any) error {
	if id, ok := data["id"].(string); ok {
		n.NodeID = id
	}

	if typeStr, ok := data["type"].(string); ok {
		n.NodeType = NodeType(typeStr)
	}

	if nextNodeID, ok := data["nextNodeId"].(string); ok {
		n.NextNodeID = nextNodeID
	}

	if schemaData, ok := data["gatherSchema"].(json.RawMessage); ok {
		var schema minds.Definition
		if err := json.Unmarshal(schemaData, &schema); err == nil {
			n.GatherSchema = &schema
		}
	}

	if maxAttempts, ok := data["maxAttempts"].(int); ok {
		n.MaxAttempts = maxAttempts
	}

	if llmPrompt, ok := data["llmPrompt"].(string); ok {
		n.LLMPrompt = llmPrompt
	}

	if fallbackNodeID, ok := data["fallbackNodeId"].(string); ok {
		n.FallbackNodeID = fallbackNodeID
	}

	if followUpPrompt, ok := data["followUpPrompt"].(string); ok {
		n.FollowUpPrompt = followUpPrompt
	}

	if extractedData, ok := data["extractedData"].(map[string]any); ok {
		n.ExtractedData = extractedData
	}

	if createdAt, ok := data["createdAt"].(time.Time); ok {
		n.CreatedAt = createdAt
	}

	if lastUpdatedAt, ok := data["lastUpdatedAt"].(time.Time); ok {
		n.LastUpdatedAt = lastUpdatedAt
	}

	return nil
}
