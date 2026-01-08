package agents

// AgentInterface defines the common interface for all agents
type AgentInterface interface {
	Kill()    // used specifically for subagents
	Execute() // Execute the agent's main logic
}
