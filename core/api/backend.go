package api

import (
	"time"

	"github.com/SomtoJF/odin/core/agents/mainagent"
	"github.com/SomtoJF/odin/core/state"
)

// Backend represents the backend API that the UI communicates with
type Backend struct {
	agent *mainagent.MainAgent

	// Channels for communication with UI
	RequestChan  chan Request       // UI sends requests here
	ResponseChan chan StateSnapshot // Backend sends state snapshots here
	UpdateChan   chan StateUpdate   // Backend sends real-time updates here

	// Control channels
	stopChan chan struct{}
}

// NewBackend creates a new backend instance
func NewBackend() *Backend {
	return &Backend{
		agent:        mainagent.NewMainAgent(),
		RequestChan:  make(chan Request, 10),
		ResponseChan: make(chan StateSnapshot, 10),
		UpdateChan:   make(chan StateUpdate, 100),
		stopChan:     make(chan struct{}),
	}
}

// Start begins the backend goroutines
func (b *Backend) Start() {
	go b.handleRequests()
	go b.publishStateSnapshots()
}

// Stop gracefully stops the backend
func (b *Backend) Stop() {
	close(b.stopChan)
	close(b.RequestChan)
}

// handleRequests processes incoming requests from the UI
func (b *Backend) handleRequests() {
	for req := range b.RequestChan {
		// Forward request to state handler
		state.HandleIncomingMessage(b.agent.State, req.Message, req.Mode)

		// Send update notification
		b.UpdateChan <- StateUpdate{
			Type:    UpdateTypeMessage,
			Payload: req.Message,
		}
	}
}

// publishStateSnapshots periodically sends state snapshots to the UI
func (b *Backend) publishStateSnapshots() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			snapshot := b.GetStateSnapshot()

			// Non-blocking send
			select {
			case b.ResponseChan <- snapshot:
			default:
				// Skip if channel is full
			}

		case <-b.stopChan:
			return
		}
	}
}

// GetStateSnapshot returns the current state snapshot
func (b *Backend) GetStateSnapshot() StateSnapshot {
	s := b.agent.State

	// Lock and read state
	s.StateMx.Lock()
	isExecuting := s.IsExecuting
	currentMode := s.AgentMode
	s.StateMx.Unlock()

	s.MessagesMx.Lock()
	messages := make([]state.Message, len(s.Messages))
	copy(messages, s.Messages)
	s.MessagesMx.Unlock()

	s.MessageQueueMx.Lock()
	queueLength := len(s.MessageQueue)
	s.MessageQueueMx.Unlock()

	return StateSnapshot{
		Messages:    messages,
		IsExecuting: isExecuting,
		CurrentMode: currentMode,
		QueueLength: queueLength,
	}
}

// GetAgent returns the main agent (for testing/debugging)
func (b *Backend) GetAgent() *mainagent.MainAgent {
	return b.agent
}
