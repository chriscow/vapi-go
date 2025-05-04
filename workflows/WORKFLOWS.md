# Workflow System for Health Coach

## Overview

The workflow system provides a flexible way to define conversational flows for the health coach application. It enables structured interactions with users while maintaining natural conversation flow.

## Design

The workflow system follows a node-based approach where each node represents a specific interaction step. The system is designed to be:

1. **Simple**: Easy to understand and implement
2. **Flexible**: Supports different types of interactions
3. **Stateful**: Maintains conversation state across sessions
4. **Persistent**: Stores workflow definitions and execution state in Firestore

## Core Components

### Nodes

Nodes are the building blocks of workflows. Each node represents a specific interaction step:

- **Say Node**: Outputs a message to the user
  - Can use a static message or generate one dynamically using an LLM
  - Transitions automatically to the next node

- **Gather Node**: Collects information from the user
  - Defines variables to extract from the conversation
  - Uses an LLM to extract information from user responses
  - Supports different data types (string, number, boolean, enum)
  - Can mark variables as required or optional

### Workflows

A workflow is a collection of nodes with defined transitions between them. Workflows:

- Have a unique identifier
- Include a name and description
- Define a starting node
- Store the nodes and their relationships

### Workflow State

The system maintains the state of each workflow execution:

- Which workflow is being executed
- Current position in the workflow
- Collected variables
- Execution history
- Completion status

### Storage

All workflow data is stored in Firestore:

- Workflow definitions are stored in the `workflows` collection
- Workflow execution states are stored in `users/{userID}/workflowStates` 

## Implementation

The workflow system is implemented as a package of Go files:

- `node.go`: Core interfaces and base types
- `say_node.go`: Say node implementation
- `gather_node.go`: Gather node implementation
- `storage.go`: Firestore storage implementation
- `engine.go`: Workflow execution engine
- `intake_workflow.go`: Example intake workflow

## Integration with VAPI

The workflow system integrates with VAPI through:

- `vapi_integration.go`: Helper functions for VAPI integration

This allows the workflow to:
1. Process incoming VAPI conversation updates
2. Extract information using LLMs
3. Send system messages to change conversation direction

## Usage Example

Here's how to use the workflow system:

```go
// Create a workflow storage
storage := NewFirestoreWorkflowStorage(firestoreClient)

// Create a workflow engine
engine := NewWorkflowEngine(storage, logger)

// Create a health coaching intake workflow
workflow := NewIntakeWorkflow()

// Save the workflow
err := engine.CreateWorkflow(ctx, workflow)
if err != nil {
    log.Fatal(err)
}

// Start the workflow for a user
state, err := engine.StartWorkflow(ctx, "intake_workflow", userID, callID)
if err != nil {
    log.Fatal(err)
}

// Process VAPI updates
err = ProcessVAPIUpdate(ctx, engine, "intake_workflow", userID, callID, messages, controlURL, logger)
if err != nil {
    log.Fatal(err)
}
```

## Implementation Plan

To implement this workflow system, follow these steps:

1. **Phase 1: Core Components**
   - Implement node interfaces and base types
   - Implement Say and Gather node types
   - Implement Firestore storage

2. **Phase 2: Workflow Engine**
   - Implement workflow execution engine
   - Add support for workflow state management
   - Implement node execution logic

3. **Phase 3: VAPI Integration**
   - Create helpers for VAPI conversation processing
   - Implement message sending functionality
   - Add workflow state transitions based on conversation updates

4. **Phase 4: Example Workflows**
   - Create the intake workflow implementation
   - Test with sample conversations
   - Refine node implementations based on testing

5. **Phase 5: Documentation and Testing**
   - Document the workflow system
   - Add unit tests
   - Create examples of workflow usage

## Future Enhancements

Future enhancements could include:

1. **Branching Logic**: Add conditional nodes for branching based on user responses
2. **Timeout Handling**: Add support for timeouts and handling user inactivity
3. **Workflow Versioning**: Track and migrate between workflow versions
4. **UI Builder**: Create a UI for defining and managing workflows
5. **Analytics**: Add tracking of workflow execution metrics 