package workflows

import (
	"context"
	"fmt"
	"time"
)

// NodeType represents the type of workflow node
type NodeType string

const (
	// NodeTypeSay represents a node that outputs a message
	NodeTypeSay NodeType = "say"
	// NodeTypeGather represents a node that collects input
	NodeTypeGather NodeType = "gather"
)

// Node represents a single node in a workflow
type Node interface {
	// ID returns the unique identifier for this node
	ID() string

	// Type returns the type of node
	Type() NodeType

	// Execute runs the node's action
	Execute(ctx context.Context, state *WorkflowState) error

	// ToMap converts the node to a map for storage
	ToMap() map[string]any

	// FromMap initializes the node from a map
	FromMap(data map[string]any) error
}

// BaseNode contains common fields for all node types
type BaseNode struct {
	NodeID        string    `json:"id"`
	NodeType      NodeType  `json:"type"`
	NextNodeID    string    `json:"nextNodeId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

// ID returns the node's ID
func (n *BaseNode) ID() string {
	return n.NodeID
}

// Type returns the node's type
func (n *BaseNode) Type() NodeType {
	return n.NodeType
}

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
	GetWorkflow(ctx context.Context, workflowID string) (*Workflow, error)

	// SaveWorkflowState saves the current state of a workflow execution
	SaveWorkflowState(ctx context.Context, state *WorkflowState) error

	// GetWorkflowState retrieves the current state of a workflow execution
	GetWorkflowState(ctx context.Context, workflowID, userID, callID string) (*WorkflowState, error)
}
