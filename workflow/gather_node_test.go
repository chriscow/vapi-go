package workflow

import (
	"context"
	"testing"
	"time"

	"github.com/chriscow/minds"
)

func TestGatherNode_Execute(t *testing.T) {
	// Create a simple schema for a user profile
	schema := &minds.Definition{
		Type:        minds.Object,
		Description: "User Profile",
		Properties: map[string]minds.Definition{
			"name": {
				Type:        minds.String,
				Description: "User's full name",
			},
			"age": {
				Type:        minds.Integer,
				Description: "User's age in years",
			},
			"email": {
				Type:        minds.String,
				Description: "User's email address",
			},
		},
		Required: []string{"name", "email"},
	}

	// Create a GatherNode with the schema
	node := NewGatherNode("profile_collector", schema, 3, "Extract the user profile information from the conversation.")

	// Initialize a workflow state
	state := &WorkflowState{
		WorkflowID:       "test-workflow",
		UserID:           "test-user",
		CallID:           "test-call",
		CurrentNodeID:    node.NodeID,
		CompletedNodeIDs: []string{},
		Variables:        make(map[string]any),
		LastMessageAt:    time.Now(),
		LastUpdatedAt:    time.Now(),
	}

	// Execute the node with the state
	err := node.Execute(context.Background(), state)
	if err != nil {
		t.Fatalf("Error executing GatherNode: %v", err)
	}

	// Verify that the required fields were "extracted"
	if _, ok := state.Variables["name"]; !ok {
		t.Errorf("Expected 'name' to be extracted")
	}
	if _, ok := state.Variables["email"]; !ok {
		t.Errorf("Expected 'email' to be extracted")
	}

	// Verify that the node was marked as completed
	if len(state.CompletedNodeIDs) != 1 || state.CompletedNodeIDs[0] != node.NodeID {
		t.Errorf("Node not properly marked as completed")
	}

	// Test ToMap and FromMap
	nodeMap := node.ToMap()

	// Create a new node and load from map
	newNode := &GatherNode{}
	err = newNode.FromMap(nodeMap)
	if err != nil {
		t.Fatalf("Error loading node from map: %v", err)
	}

	// Verify that the new node has the same properties
	if newNode.NodeID != node.NodeID {
		t.Errorf("Expected NodeID to be %s, got %s", node.NodeID, newNode.NodeID)
	}
	if newNode.MaxAttempts != node.MaxAttempts {
		t.Errorf("Expected MaxAttempts to be %d, got %d", node.MaxAttempts, newNode.MaxAttempts)
	}
	if newNode.LLMPrompt != node.LLMPrompt {
		t.Errorf("Expected LLMPrompt to be %s, got %s", node.LLMPrompt, newNode.LLMPrompt)
	}
}
