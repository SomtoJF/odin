package state

import (
	"time"
)

// HandleIncomingMessage routes messages to queue or immediate processing
func HandleIncomingMessage(s *State, body string, mode AgentMode) {
	s.StateMx.Lock()
	isExecuting := s.IsExecuting
	s.StateMx.Unlock()

	if isExecuting {
		// Agent is busy - add to queue
		s.MessageQueueMx.Lock()
		s.MessageQueue = append(s.MessageQueue, QueuedMessage{
			Body:      body,
			Mode:      mode,
			Timestamp: time.Now(),
		})
		s.MessageQueueMx.Unlock()
	} else {
		// Agent is idle - process immediately
		go ProcessMessage(s, body, mode)
	}
}

// ProcessMessage processes a single message through the iteration loop
func ProcessMessage(s *State, body string, mode AgentMode) {
	// Mark as executing
	s.StateMx.Lock()
	s.IsExecuting = true
	s.AgentMode = mode
	s.StateMx.Unlock()

	// Add message to history
	message := Message{Body: body}
	s.MessagesMx.Lock()
	s.Messages = append(s.Messages, message)
	_ = len(s.Messages) - 1 // messageIndex for future use
	s.MessagesMx.Unlock()

	// TODO: Get tools for this mode
	// tools := GetModeTools(mode, false)

	// TODO: Run iteration loop until task completed
	// taskCompleted := false
	// for !taskCompleted {
	//     Build planner input with cached files info
	//     Call planner
	//     Execute tool or mark as completed
	// }

	// TODO: Return result to user
	// ReturnAnswerToUser(s.Messages[messageIndex].AnswerSummary)

	// Mark as idle
	s.StateMx.Lock()
	s.IsExecuting = false
	s.StateMx.Unlock()

	// Process next message in queue if any
	ProcessNextMessageInQueue(s)
}

// ProcessNextMessageInQueue dequeues and processes the next message
func ProcessNextMessageInQueue(s *State) {
	s.MessageQueueMx.Lock()

	if len(s.MessageQueue) == 0 {
		s.MessageQueueMx.Unlock()
		return // No more messages
	}

	// Dequeue first message
	nextMessage := s.MessageQueue[0]
	s.MessageQueue = s.MessageQueue[1:]
	s.MessageQueueMx.Unlock()

	// Process it
	go ProcessMessage(s, nextMessage.Body, nextMessage.Mode)
}

// GetModeTools returns tools available for the given mode
func GetModeTools(mode AgentMode, isSubAgent bool) []interface{} {
	// TODO: Implement tool selection based on mode
	// Base tools available to all modes (read-only tools)
	// baseTools := []Tool{
	//     LS, Grep, Glob, WebFetch, ContextSummarizer, ExecuteCommand, TodoWrite
	// }

	// Add mode-specific tools based on mode
	// switch mode {
	// case AgentModeAskMode, AgentModePlanMode:
	//     Only base tools
	// case AgentModeEdit:
	//     Add WriteFile, EditTool, MultiEditTool, InitTool
	// }

	// Add AgentTool only for main agents (not subagents)
	return []interface{}{}
}
