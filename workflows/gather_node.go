package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// DataType represents the type of data to gather
type DataType string

const (
	// DataTypeString represents string data
	DataTypeString DataType = "string"
	// DataTypeNumber represents number data
	DataTypeNumber DataType = "number"
	// DataTypeBoolean represents boolean data
	DataTypeBoolean DataType = "boolean"
	// DataTypeEnum represents enum data
	DataTypeEnum DataType = "enum"
)

// GatherVariable represents a variable to gather from the user
type GatherVariable struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DataType    DataType `json:"dataType"`
	Required    bool     `json:"required"`
	EnumValues  []string `json:"enumValues,omitempty"` // Only for DataTypeEnum
}

// GatherNode represents a node that collects input from the user
type GatherNode struct {
	BaseNode
	Variables      []GatherVariable `json:"variables"`
	MaxAttempts    int              `json:"maxAttempts"`
	LLMPrompt      string           `json:"llmPrompt"`
	FallbackNodeID string           `json:"fallbackNodeId,omitempty"`
	FollowUpPrompt string           `json:"followUpPrompt,omitempty"`
	ExtractedData  map[string]any   `json:"extractedData,omitempty"`
}

// NewGatherNode creates a new GatherNode
func NewGatherNode(id string, variables []GatherVariable, maxAttempts int, llmPrompt string) *GatherNode {
	now := time.Now()
	return &GatherNode{
		BaseNode: BaseNode{
			NodeID:        id,
			NodeType:      NodeTypeGather,
			CreatedAt:     now,
			LastUpdatedAt: now,
		},
		Variables:   variables,
		MaxAttempts: maxAttempts,
		LLMPrompt:   llmPrompt,
	}
}

// Execute runs the Gather node's action
func (n *GatherNode) Execute(ctx context.Context, state *WorkflowState) error {
	logger := slog.Default().With("node", n.NodeID, "type", n.NodeType)

	// For MVP, we'll use a simple approach to gather data
	// In a real implementation, we would use an LLM API here

	// Check if we need to extract data from the conversation
	missing := make([]GatherVariable, 0)

	// Initialize the extracted data map if not already present
	if n.ExtractedData == nil {
		n.ExtractedData = make(map[string]any)
	}

	// Check which variables are missing
	for _, variable := range n.Variables {
		if _, ok := n.ExtractedData[variable.Name]; !ok {
			// In a real implementation, we would check if the variable is in state.Variables
			// from previous messages. For MVP, we'll just mark it as missing.
			missing = append(missing, variable)
		}
	}

	// If there are missing variables, generate a prompt to extract them
	if len(missing) > 0 {
		// Generate prompt for the LLM to extract missing data
		prompt := n.LLMPrompt + "\n\nPlease extract the following information:\n"
		for _, v := range missing {
			prompt += fmt.Sprintf("- %s: %s (Type: %s, Required: %t)\n",
				v.Name, v.Description, v.DataType, v.Required)
		}

		// For MVP, we'll simulate data extraction
		// In a real implementation, we would call the LLM API here
		logger.Info("would call LLM to extract data", "prompt", prompt)

		// For now, just add some dummy data
		// In a real implementation, this would come from the LLM's response
		for _, variable := range missing {
			switch variable.DataType {
			case DataTypeString:
				n.ExtractedData[variable.Name] = fmt.Sprintf("Sample %s", variable.Name)
			case DataTypeNumber:
				n.ExtractedData[variable.Name] = 42
			case DataTypeBoolean:
				n.ExtractedData[variable.Name] = true
			case DataTypeEnum:
				if len(variable.EnumValues) > 0 {
					n.ExtractedData[variable.Name] = variable.EnumValues[0]
				} else {
					n.ExtractedData[variable.Name] = "unknown"
				}
			}
		}

		// In a real implementation, some variables might still be missing after extraction
		// For MVP, we'll assume all variables were extracted successfully

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

// ToMap converts the GatherNode to a map for storage
func (n *GatherNode) ToMap() map[string]any {
	variablesMaps := make([]map[string]any, len(n.Variables))
	for i, v := range n.Variables {
		variablesMaps[i] = map[string]any{
			"name":        v.Name,
			"description": v.Description,
			"dataType":    string(v.DataType),
			"required":    v.Required,
			"enumValues":  v.EnumValues,
		}
	}

	return map[string]any{
		"id":             n.NodeID,
		"type":           string(n.NodeType),
		"nextNodeId":     n.NextNodeID,
		"variables":      variablesMaps,
		"maxAttempts":    n.MaxAttempts,
		"llmPrompt":      n.LLMPrompt,
		"fallbackNodeId": n.FallbackNodeID,
		"followUpPrompt": n.FollowUpPrompt,
		"extractedData":  n.ExtractedData,
		"createdAt":      n.CreatedAt,
		"lastUpdatedAt":  n.LastUpdatedAt,
	}
}

// FromMap initializes the GatherNode from a map
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

	if variablesData, ok := data["variables"].([]map[string]any); ok {
		n.Variables = make([]GatherVariable, len(variablesData))
		for i, varData := range variablesData {
			var variable GatherVariable

			if name, ok := varData["name"].(string); ok {
				variable.Name = name
			}

			if desc, ok := varData["description"].(string); ok {
				variable.Description = desc
			}

			if dataType, ok := varData["dataType"].(string); ok {
				variable.DataType = DataType(dataType)
			}

			if required, ok := varData["required"].(bool); ok {
				variable.Required = required
			}

			if enumValues, ok := varData["enumValues"].([]string); ok {
				variable.EnumValues = enumValues
			}

			n.Variables[i] = variable
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
