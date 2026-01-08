package api

import (
	"github.com/SomtoJF/odin/core/state"
)

// Request represents a message from the UI to the backend
type Request struct {
	Message string
	Mode    state.AgentMode
}

// StateSnapshot represents the current state of the backend
type StateSnapshot struct {
	Messages    []state.Message
	IsExecuting bool
	CurrentMode state.AgentMode
	QueueLength int
}

// StateUpdate represents a real-time update from the backend
type StateUpdate struct {
	Type    UpdateType  // Type of update
	Payload interface{} // Update payload
}

// UpdateType defines the type of state update
type UpdateType string

const (
	UpdateTypeMessage     UpdateType = "message"
	UpdateTypeTodo        UpdateType = "todo"
	UpdateTypeStatus      UpdateType = "status"
	UpdateTypeToolHistory UpdateType = "tool_history"
)
