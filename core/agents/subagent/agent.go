package subagent

import (
	"github.com/SomtoJF/odin/core/state"
)

// SubAgent represents a spawned child agent
type SubAgent struct {
	ID          uint
	State       *state.State
	ParentState *state.State
	Mode        state.AgentMode // The mode this subagent was spawned in
}

// NewSubAgent creates a new SubAgent instance
func NewSubAgent(mode state.AgentMode, parentState *state.State) *SubAgent {
	return &SubAgent{
		State:       state.NewState(),
		ParentState: parentState,
		Mode:        mode,
	}
}

// Execute runs the subagent's main logic
func (sa *SubAgent) Execute() {
	// TODO: Implement subagent execution logic
	// Call planner on latest Message
	// Planner calls tools and resumes iteration loop
	// When the problem has been solved, we end the execution and kill the agent
	// Note: Subagents cannot spawn other subagents - AgentTool is not in their tool list
	sa.Kill()
}

// Kill removes the subagent from the parent state
func (sa *SubAgent) Kill() {
	// Find current running subagent in parent state and remove it
	// TODO: Implement removal logic
}
