package tools

// ToolName represents a tool identifier
type ToolName string

const (
	ToolNameLS                ToolName = "LS"
	ToolNameGrep              ToolName = "Grep"
	ToolNameGlob              ToolName = "Glob"
	ToolNameWriteFile         ToolName = "WriteFile"
	ToolNameEditTool          ToolName = "EditTool"
	ToolNameMultiEditTool     ToolName = "MultiEditTool"
	ToolNameTodoWrite         ToolName = "TodoWrite"
	ToolNameAgentTool         ToolName = "AgentTool"
	ToolNameWebFetch          ToolName = "WebFetch"
	ToolNameContextSummarizer ToolName = "ContextSummarizer"
	ToolNameInitTool          ToolName = "InitTool"
	ToolNameExecuteCommand    ToolName = "ExecuteCommand"
)

// Tool defines the interface for all tools
type Tool interface {
	PreHook(input map[string]interface{}) (interface{}, error)
	Execute(prehookResponse interface{}, input map[string]interface{}) (interface{}, error)
	PostHook(executionResult interface{}, input map[string]interface{}) (*ToolResponse, error)
}

// ToolResponse represents the result of a tool execution
type ToolResponse struct {
	Data        map[string]interface{} // can be any type. should be unmarshaled into the specific return type of the tool
	Description string                 // Short description of the change that the tool made
}

// ToolRegistry holds all available tools
var ToolRegistry = make(map[ToolName]Tool)

// RegisterTool adds a tool to the registry
func RegisterTool(name ToolName, tool Tool) {
	ToolRegistry[name] = tool
}

// ExecuteTool executes a tool by name with the given input
func ExecuteTool(toolName ToolName, input map[string]interface{}) (*ToolResponse, error) {
	tool, exists := ToolRegistry[toolName]
	if !exists {
		return nil, &ToolNotFoundError{ToolName: string(toolName)}
	}

	// Execute PreHook
	prehookResponse, err := tool.PreHook(input)
	if err != nil {
		return nil, err
	}

	// Execute main logic
	executionResult, err := tool.Execute(prehookResponse, input)
	if err != nil {
		return nil, err
	}

	// Execute PostHook
	toolResponse, err := tool.PostHook(executionResult, input)
	if err != nil {
		return nil, err
	}

	return toolResponse, nil
}

// ToolNotFoundError is returned when a tool is not found in the registry
type ToolNotFoundError struct {
	ToolName string
}

func (e *ToolNotFoundError) Error() string {
	return "tool not found: " + e.ToolName
}
