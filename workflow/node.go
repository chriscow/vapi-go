package workflow

import (
	"context"
	"time"

	"github.com/chriscow/vapi-go"
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
	Execute(ctx context.Context, state *WorkflowState, messages []vapi.Message) error

	// ToMap converts the node to a map for storage
	// ToMap() map[string]any

	// FromMap initializes the node from a map
	// FromMap(data map[string]any) error
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
