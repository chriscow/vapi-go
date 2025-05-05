package workflow

import (
	"context"
	"testing"
	"time"

	"github.com/chriscow/minds"
	"github.com/chriscow/vapi-go"
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
	messages := []vapi.Message{}
	err := node.Execute(context.Background(), state, messages)
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
}
