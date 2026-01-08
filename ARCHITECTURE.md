# Odin Code Architecture

## Overview

Odin Code implements a **frontend/backend separation** using Go channels for communication. The UI and backend run in the same process but are completely decoupled, communicating only via typed channels.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                         main.go                                   │
│                                                                   │
│  backend := api.NewBackend()                                      │
│  backend.Start()                                                  │
│                                                                   │
│  ui := components.NewUI(                                          │
│      backend.RequestChan,  ← Channels injected via constructor  │
│      backend.ResponseChan,                                        │
│      backend.UpdateChan,                                          │
│  )                                                                │
│                                                                   │
│  ui.Run()  ← Blocks on main thread                              │
└──────────────────────────────────────────────────────────────────┘

                            ↓         ↓         ↓

┌─────────────────────────┐         Channels         ┌──────────────────────────┐
│      UI (Frontend)      │◄────────────────────────►│   Backend (API Layer)    │
│                         │                           │                          │
│  ui/ui.go               │  Request (msg, mode) →   │  core/api/backend.go     │
│  ┌─────────────────┐   │                           │  ┌──────────────────┐    │
│  │ tview.App       │   │  ← StateSnapshot (100ms)  │  │ MainAgent        │    │
│  │ messageView     │   │                           │  │ State            │    │
│  │ inputField      │   │  ← UpdateChan (events)    │  │ handleRequests() │    │
│  │ todoList        │   │                           │  │ publishSnapshots()│    │
│  │ statusBar       │   │                           │  └──────────────────┘    │
│  └─────────────────┘   │                           │                          │
│                         │                           │  Goroutines:             │
│  Methods:               │                           │  - Request handler       │
│  - listenToBackend()    │                           │  - Snapshot publisher    │
│  - updateUIFromSnapshot()│                          │                          │
│  - handleUserInput()    │                           │  State Management:       │
│  - updateMessageView()  │                           │  - Messages              │
│  - updateTodoList()     │                           │  - MessageQueue          │
│  - updateStatusBar()    │                           │  - IsExecuting           │
│                         │                           │  - SubAgents             │
└─────────────────────────┘                           └──────────────────────────┘
         No direct access                                     No direct access
         to State/Agent! ✓                                    to UI components! ✓
```

## Key Design Principles

### 1. Separation of Concerns

**UI Responsibilities:**
- Render tview components
- Handle user input (keyboard, text entry)
- Send requests to backend via channels
- Receive and display state snapshots
- Update display based on backend state

**Backend Responsibilities:**
- Process user messages
- Manage agent state
- Execute planner iteration loop
- Publish state snapshots
- Handle message queue

### 2. Dependency Injection via Channels

The UI receives **only channels** as dependencies:

```go
// ui/ui.go
type UI struct {
    requestChan  chan<- api.Request       // Write-only
    responseChan <-chan api.StateSnapshot // Read-only
    updateChan   <-chan api.StateUpdate   // Read-only
    // ... tview components
}

func NewUI(
    requestChan chan<- api.Request,
    responseChan <-chan api.StateSnapshot,
    updateChan <-chan api.StateUpdate,
) *UI {
    // Constructor injection - idiomatic Go
}
```

**Benefits:**
- UI cannot access State or MainAgent directly
- Clear contract between layers
- Easy to test with mock channels
- Type-safe communication
- Future-proof (can swap to network later)

### 3. Unidirectional Data Flow

```
User Input → UI → RequestChan → Backend → State → ResponseChan → UI → Display
                                   ↓
                              UpdateChan
                                   ↓
                                  UI
```

- **User → Backend**: Via `RequestChan` (commands flow down)
- **Backend → UI**: Via `ResponseChan` and `UpdateChan` (state flows up)
- **No circular dependencies**

## Channel Types

### 1. RequestChan (UI → Backend)

```go
type Request struct {
    Message string          // User's message
    Mode    state.AgentMode // ask_mode, edit_mode, plan_mode
}
```

**Usage:**
```go
// UI sends request
ui.requestChan <- api.Request{
    Message: "Fix the bug in auth.go",
    Mode:    state.AgentModeEdit,
}
```

### 2. ResponseChan (Backend → UI)

```go
type StateSnapshot struct {
    Messages    []state.Message // Full conversation history
    IsExecuting bool            // Agent busy?
    CurrentMode state.AgentMode // Current mode
    QueueLength int             // Messages in queue
}
```

**Usage:**
```go
// UI receives snapshot every 100ms
snapshot := <-ui.responseChan
ui.updateMessageView(snapshot.Messages)
ui.updateTodoList(snapshot.Messages)
```

### 3. UpdateChan (Backend → UI)

```go
type StateUpdate struct {
    Type    UpdateType  // "message", "todo", "status", "tool_history"
    Payload interface{} // Event-specific data
}
```

**Usage:**
```go
// Backend sends real-time events
backend.updateChan <- StateUpdate{
    Type:    api.UpdateTypeMessage,
    Payload: "Processing your request...",
}
```

## Message Flow Example

### User sends "Fix bug in auth.go"

```
1. User types in inputField, presses Enter

2. UI.handleUserInput() called:
   ui.requestChan <- Request{
       Message: "Fix bug in auth.go",
       Mode:    state.AgentModeEdit,
   }

3. Backend.handleRequests() receives:
   req := <-b.RequestChan
   state.HandleIncomingMessage(b.agent.State, req.Message, req.Mode)

4. State processing:
   - Check IsExecuting
   - If busy: Add to MessageQueue
   - If idle: Start ProcessMessage()

5. Backend publishes snapshots (every 100ms):
   snapshot := StateSnapshot{
       Messages:    [...messages with new request...],
       IsExecuting: true,
       CurrentMode: state.AgentModeEdit,
   }
   b.ResponseChan <- snapshot

6. UI.listenToBackend() receives:
   snapshot := <-ui.responseChan
   ui.updateUIFromSnapshot(snapshot)

7. UI updates display:
   - Message view shows new message
   - Status bar shows "Executing"
   - Todos appear as agent works

8. Agent completes, backend sends final snapshot:
   snapshot := StateSnapshot{
       Messages: [...with AnswerSummary...],
       IsExecuting: false,
   }

9. UI displays result and returns to idle
```

## Testing Strategy

### Testing UI (with Mock Backend)

```go
func TestUI_HandleUserInput(t *testing.T) {
    // Create mock channels
    reqChan := make(chan api.Request, 1)
    respChan := make(chan api.StateSnapshot, 1)
    updateChan := make(chan api.StateUpdate, 1)

    // Create UI with mock channels
    ui := components.NewUI(reqChan, respChan, updateChan)

    // Simulate user input
    ui.handleUserInput("test message")

    // Verify request sent
    req := <-reqChan
    assert.Equal(t, "test message", req.Message)
}
```

### Testing Backend (without UI)

```go
func TestBackend_HandleRequest(t *testing.T) {
    // Create backend
    backend := api.NewBackend()
    backend.Start()

    // Send request
    backend.RequestChan <- api.Request{
        Message: "test",
        Mode:    state.AgentModeAsk,
    }

    // Verify state snapshot
    snapshot := <-backend.ResponseChan
    assert.NotEmpty(t, snapshot.Messages)
}
```

## Why This Architecture?

### ✅ Advantages

1. **True Separation**: UI cannot access State/Agent internals
2. **Testability**: Each layer can be tested independently with mocks
3. **Maintainability**: Clear boundaries, easy to modify either side
4. **Type Safety**: Channels provide compile-time type checking
5. **Concurrency**: Channels handle goroutine communication safely
6. **Future-Proof**: Easy to swap channels for network calls (HTTP, gRPC)
7. **Idiomatic Go**: Uses channels as intended, constructor injection

### ❌ What We Avoid

1. **Tight Coupling**: UI doesn't directly call state methods
2. **Shared Mutable State**: No direct access to State fields
3. **Circular Dependencies**: Unidirectional data flow only
4. **Global State**: Everything injected via constructors
5. **Hidden Dependencies**: All dependencies explicit in constructors

## Comparison to Your Original Design

### Before (No DI)
```go
// main.go
ui := components.NewUI()
ui.Run()

// ui/ui.go
type UI struct{} // Empty!
func NewUI() *UI { return &UI{} }
```

**Problems:**
- UI had no access to backend
- Couldn't display agent state
- Impossible to test
- Just showed static ASCII art

### After (Channel-Based DI)
```go
// main.go
backend := api.NewBackend()
backend.Start()
ui := components.NewUI(backend.RequestChan, backend.ResponseChan, backend.UpdateChan)
ui.Run()

// ui/ui.go
type UI struct {
    requestChan  chan<- api.Request
    responseChan <-chan api.StateSnapshot
    updateChan   <-chan api.StateUpdate
    // ... components
}
```

**Benefits:**
- ✅ UI can communicate with backend
- ✅ Proper dependency injection
- ✅ Testable with mock channels
- ✅ Full-featured TUI with real-time updates
- ✅ Frontend/backend separation maintained

## Future Enhancements

### Easy to Add:
1. **Network API**: Replace channels with HTTP handlers
2. **Multiple UIs**: CLI, Web UI, VS Code extension all use same backend
3. **Remote Backend**: Backend runs on server, UI connects via network
4. **Persistent State**: Add Redis/database for state publishing
5. **Metrics**: Monitor channel buffer sizes, message latency

### Example: Adding HTTP API

```go
// Could easily add HTTP handlers alongside channels
type Backend struct {
    agent        *mainagent.MainAgent
    requestChan  chan Request  // For TUI
    responseChan chan StateSnapshot
    httpServer   *http.Server  // For web clients
}

func (b *Backend) HandleHTTPRequest(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var req Request
    json.NewDecoder(r.Body).Decode(&req)

    // Use same state management
    state.HandleIncomingMessage(b.agent.State, req.Message, req.Mode)

    // Return snapshot
    snapshot := b.GetStateSnapshot()
    json.NewEncoder(w).Encode(snapshot)
}
```

The channel-based architecture makes this trivial to add without touching UI code!

## Summary

The Odin Code architecture achieves **true frontend/backend separation** using:

1. **Channel-based DI**: UI receives only channels, not concrete types
2. **Unidirectional data flow**: User → Backend → State → UI
3. **Type safety**: All communication via typed channels
4. **Testability**: Each layer tested independently with mocks
5. **Idiomatic Go**: Constructor injection, interfaces, channels

This is the **correct way** to implement frontend/backend separation in Go while maintaining the simplicity of running in a single process.
