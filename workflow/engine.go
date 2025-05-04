package workflow

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// WorkflowEngine manages the execution, state, and transitions of a workflow.
// It is responsible for starting workflows, processing conversation updates,
// retrieving the current node's message, and persisting workflow state.
//
// Typical usage:
//
//	engine := NewWorkflowEngine(storage, logger)
//	err := engine.CreateWorkflow(ctx, workflow)
//	state, err := engine.StartWorkflow(ctx, workflowID, userID, callID)
//	state, err := engine.ProcessConversationUpdate(ctx, workflowID, userID, callID, messages)
//	msg, err := engine.GetCurrentNodeMessage(ctx, workflowID, userID, callID)
type WorkflowEngine struct {
	// storage provides persistence for workflow definitions and state.
	storage WorkflowStorage
	// logger is used for logging workflow execution and errors.
	logger *slog.Logger
}

// NewWorkflowEngine returns a new WorkflowEngine using the provided storage and logger.
// If logger is nil, slog.Default() is used.
func NewWorkflowEngine(storage WorkflowStorage, logger *slog.Logger) *WorkflowEngine {
	if logger == nil {
		logger = slog.Default()
	}

	return &WorkflowEngine{
		storage: storage,
		logger:  logger,
	}
}

// CreateWorkflow persists a new workflow definition.
// Returns an error if the workflow is invalid or cannot be saved.
func (e *WorkflowEngine) CreateWorkflow(ctx context.Context, workflow *Workflow) error {
	if workflow.ID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}

	if workflow.StartNodeID == "" {
		return fmt.Errorf("workflow must have a start node")
	}

	if len(workflow.Nodes) == 0 {
		return fmt.Errorf("workflow must have at least one node")
	}

	if _, ok := workflow.Nodes[workflow.StartNodeID]; !ok {
		return fmt.Errorf("start node ID '%s' not found in workflow nodes", workflow.StartNodeID)
	}

	now := time.Now()
	if workflow.CreatedAt.IsZero() {
		workflow.CreatedAt = now
	}
	workflow.UpdatedAt = now

	return e.storage.SaveWorkflow(ctx, workflow)
}

// StartWorkflow initializes or resumes a workflow execution for a given user and call.
// If no state exists, it creates a new state starting at the workflow's start node.
// Returns the current workflow state or an error.
func (e *WorkflowEngine) StartWorkflow(ctx context.Context, workflowID, userID, callID string) (*WorkflowState, error) {
	workflow, found, err := e.storage.GetWorkflow(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}

	// Check if a state already exists for this workflow execution
	state, err := e.storage.GetWorkflowState(ctx, workflowID, userID, callID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow state: %w", err)
	}

	// If no current node is set, set it to the start node
	if state.CurrentNodeID == "" {
		state.CurrentNodeID = workflow.StartNodeID
		// Save the initial state
		if err := e.storage.SaveWorkflowState(ctx, state); err != nil {
			return nil, fmt.Errorf("failed to save initial workflow state: %w", err)
		}
	}

	e.logger.Info("workflow started", "workflowID", workflowID, "userID", userID, "callID", callID)
	return state, nil
}

// ProcessConversationUpdate processes a new conversation update for a workflow execution.
// It loads the workflow and state, executes the current node, advances the workflow as needed,
// and persists the updated state. Returns the updated workflow state or an error.
func (e *WorkflowEngine) ProcessConversationUpdate(ctx context.Context, workflowID, userID, callID string, messages []map[string]any) (*WorkflowState, error) {
	logger := e.logger.With(
		"workflowID", workflowID,
		"userID", userID,
		"callID", callID,
	)

	// Get workflow and current state
	workflow, found, err := e.storage.GetWorkflow(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}

	state, err := e.storage.GetWorkflowState(ctx, workflowID, userID, callID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow state: %w", err)
	}

	if state.IsComplete {
		logger.Info("workflow already complete")
		return state, nil
	}

	// Update the last message time
	state.LastMessageAt = time.Now()

	// Get the current node
	currentNode, ok := workflow.Nodes[state.CurrentNodeID]
	if !ok {
		return nil, ErrNodeNotFound{NodeID: state.CurrentNodeID}
	}

	logger.Info("processing message for node", "nodeID", currentNode.ID(), "nodeType", currentNode.Type())

	// Execute the node
	if err := currentNode.Execute(ctx, state); err != nil {
		logger.Error("node execution failed", "nodeID", currentNode.ID(), "error", err)
		return nil, fmt.Errorf("node execution failed: %w", err)
	}

	// Save the updated state
	if err := e.storage.SaveWorkflowState(ctx, state); err != nil {
		logger.Error("failed to save workflow state", "error", err)
		return nil, fmt.Errorf("failed to save workflow state: %w", err)
	}

	// Check if we need to continue to the next node
	if !state.IsComplete && state.CurrentNodeID != currentNode.ID() {
		// Get the next node
		nextNode, ok := workflow.Nodes[state.CurrentNodeID]
		if !ok {
			logger.Warn("next node not found", "nextNodeID", state.CurrentNodeID)
			return state, nil
		}

		// If the next node is a Say node, execute it immediately
		if nextNode.Type() == NodeTypeSay {
			logger.Info("executing next node automatically", "nodeID", nextNode.ID(), "nodeType", nextNode.Type())

			// Execute the node
			if err := nextNode.Execute(ctx, state); err != nil {
				logger.Error("node execution failed", "nodeID", nextNode.ID(), "error", err)
				return nil, fmt.Errorf("node execution failed: %w", err)
			}

			// Save the updated state
			if err := e.storage.SaveWorkflowState(ctx, state); err != nil {
				logger.Error("failed to save workflow state", "error", err)
				return nil, fmt.Errorf("failed to save workflow state: %w", err)
			}
		}
	}

	return state, nil
}

// GetCurrentNodeMessage returns the message to send for the current node in the workflow execution.
// It loads the workflow and state, determines the current node, and returns the appropriate message
// (e.g., a Say node's message or a Gather node's prompt for missing variables).
func (e *WorkflowEngine) GetCurrentNodeMessage(ctx context.Context, workflowID, userID, callID string) (string, error) {
	logger := e.logger.With(
		"workflowID", workflowID,
		"userID", userID,
		"callID", callID,
	)

	// Get workflow and current state
	workflow, found, err := e.storage.GetWorkflow(ctx, workflowID)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow: %w", err)
	}

	if !found {
		return "", fmt.Errorf("workflow not found: %s", workflowID)
	}

	state, err := e.storage.GetWorkflowState(ctx, workflowID, userID, callID)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow state: %w", err)
	}

	if state.IsComplete {
		logger.Info("workflow already complete")
		return "Workflow complete", nil
	}

	// Get the current node
	currentNode, ok := workflow.Nodes[state.CurrentNodeID]
	if !ok {
		return "", ErrNodeNotFound{NodeID: state.CurrentNodeID}
	}

	logger.Info("getting message for node", "nodeID", currentNode.ID(), "nodeType", currentNode.Type())

	// Return message based on node type
	switch node := currentNode.(type) {
	case *SayNode:
		if node.MessageType == "exact" {
			return node.Message, nil
		} else if node.MessageType == "generated" {
			// For MVP, just return the prompt
			// In a real implementation, we would call the LLM
			return fmt.Sprintf("Generated message based on: %s", node.LLMPrompt), nil
		}
	case *GatherNode:
		// For Gather nodes, check if we need to prompt for specific variables
		missing := make([]GatherVariable, 0)

		for _, variable := range node.Variables {
			if _, ok := state.Variables[variable.Name]; !ok {
				missing = append(missing, variable)
			}
		}

		if len(missing) > 0 {
			// Generate a prompt for the missing variables
			prompt := "I need to gather some information from you:\n"
			for _, v := range missing {
				prompt += fmt.Sprintf("- %s: %s\n", v.Name, v.Description)
			}
			return prompt, nil
		}

		return "Thank you for providing that information.", nil
	}

	return "Please continue with our conversation.", nil
}
