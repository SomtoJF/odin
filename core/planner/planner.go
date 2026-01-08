package planner

import (
	"github.com/SomtoJF/odin/core/state"
)

// PlannerInput represents input to the planner
type PlannerInput struct {
	LatestMessage      state.Message            // latest message sent from the user
	AvailableTools     []ToolDescription        // list of tools the planner is allowed to call
	Context            []state.ContextItem      // context retrieved from state
	CustomInstructions string                   // Data from the ODIN.md
	Config             state.Config             // Configuration
	CachedFiles        []state.CachedFileInfo   // List of files currently in cache
}

// ToolDescription describes a tool available to the planner
type ToolDescription struct {
	ToolName        string
	ToolDescription string
	ToolInput       map[string]interface{}
}

// ToolExecutionCallInput represents a tool call request
type ToolExecutionCallInput struct {
	ToolName  string
	ToolInput map[string]interface{}
}

// PlannerOutput represents the planner's decision
type PlannerOutput struct {
	Explanation   string                 `json:"explanation"`
	TaskCompleted bool                   `json:"taskCompleted"`
	ExecuteTool   ToolExecutionCallInput `json:"executeTool"` // Tool to execute next
}

// CallPlanner invokes the planner with the given input
func CallPlanner(input PlannerInput) PlannerOutput {
	// TODO: Implement planner logic
	// This will call the LLM with the appropriate system prompt
	// based on the agent mode and return the next action to take
	return PlannerOutput{}
}
