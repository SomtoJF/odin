package mainagent

import (
	"github.com/SomtoJF/odin/core/state"
)

// MainAgent represents the primary agent
type MainAgent struct {
	State *state.State
}

// NewMainAgent creates a new MainAgent instance
func NewMainAgent() *MainAgent {
	return &MainAgent{
		State: state.NewState(),
	}
}

// Execute processes a message in the specified mode
func (ma *MainAgent) Execute(body string, mode state.AgentMode) {
	// TODO: Implement ProcessMessage logic
	// ProcessMessage(ma.State, body, mode)
}

// Kill is a no-op for MainAgent (required by interface)
func (ma *MainAgent) Kill() {
	// MainAgent doesn't need cleanup
}
