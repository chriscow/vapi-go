package workflow

import (
	"context"
	"fmt"
	"time"
)

// WorkflowState represents the current state of a workflow execution
type WorkflowState struct {
	WorkflowID       string         `json:"workflowId"`
	UserID           string         `json:"userId"`
	CallID           string         `json:"callId"`
	CurrentNodeID    string         `json:"currentNodeId"`
	CompletedNodeIDs []string       `json:"completedNodeIds"`
	Variables        map[string]any `json:"variables"`
	LastMessageAt    time.Time      `json:"lastMessageAt"`
	LastUpdatedAt    time.Time      `json:"lastUpdatedAt"`
	IsComplete       bool           `json:"isComplete"`
}

// Workflow represents a complete workflow definition
type Workflow struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Nodes       map[string]Node `json:"nodes"`
	StartNodeID string          `json:"startNodeId"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ErrNodeNotFound is returned when a node is not found in a workflow
type ErrNodeNotFound struct {
	NodeID string
}

func (e ErrNodeNotFound) Error() string {
	return fmt.Sprintf("node not found: %s", e.NodeID)
}

// WorkflowStorage defines the interface for workflow storage
type WorkflowStorage interface {
	// SaveWorkflow saves a workflow definition
	SaveWorkflow(ctx context.Context, workflow *Workflow) error

	// GetWorkflow retrieves a workflow by ID
	GetWorkflow(ctx context.Context, workflowID string) (*Workflow, bool, error)

	// SaveWorkflowState saves the current state of a workflow execution
	SaveWorkflowState(ctx context.Context, state *WorkflowState) error

	// GetWorkflowState retrieves the current state of a workflow execution
	GetWorkflowState(ctx context.Context, workflowID, userID, callID string) (*WorkflowState, error)
}
