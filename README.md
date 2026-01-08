# Odin Code

An AI-powered coding assistant with three operational modes: Ask, Edit, and Plan.

## Architecture Overview

Odin Code implements a sophisticated agent-based architecture with the following key components:

### Three Operational Modes

1. **Ask Mode**: Research and question-answering with read-only tools
2. **Edit Mode**: Full access to all tools including file writing capabilities
3. **Plan Mode**: Strategic planning with read-only access

### Core Components

```
core/
├── agents/          # Agent implementations
│   ├── agent.go     # Agent interface
│   ├── mainagent/   # Main agent implementation
│   └── subagent/    # Sub-agent implementation
├── api/             # Backend API layer
│   ├── backend.go   # Backend with channel-based communication
│   └── types.go     # Request/Response types
├── cache/           # File caching system (future)
├── planner/         # Planning and decision logic
│   └── planner.go   # Planner types and functions
├── state/           # State management
│   ├── types.go     # Core type definitions
│   ├── state.go     # State operations
│   └── handler.go   # Message handling
└── tools/           # Tool system
    └── tool.go      # Tool interface and registry

ui/
├── ui.go            # Main UI with channel-based backend communication
└── components/      # UI components
    └── introascii/  # ASCII art component
```

## Key Features

### File Cache System

- **Intelligent Caching**: Files are cached on first read with O(1) lookups
- **Partial Caching**: Supports caching specific line ranges for large files
- **LLM-Based Validation**: Uses cheap models to validate cache sufficiency before edits
- **Smart Eviction**: Hybrid strategy preferring unmodified files
- **Memory Bounded**: Configurable limits (100MB default, 500 files max)

### Message Queue System

- **Automatic Queuing**: Messages are queued when agent is busy
- **Sequential Processing**: Each message processed completely before the next
- **Dynamic Mode Switching**: Each message can specify its own mode

### Agent System

- **Main Agent**: Primary agent with full capabilities
- **Sub-Agents**: Spawned for specific tasks, cannot spawn other sub-agents
- **Thread-Safe**: All state mutations protected by mutexes

### Frontend/Backend Separation

- **Channel-Based Communication**: UI and backend communicate via Go channels
- **Decoupled Architecture**: UI has no direct access to State or MainAgent
- **Real-Time Updates**: Backend publishes state snapshots every 100ms
- **Request/Response Model**: UI sends requests, backend processes them asynchronously

## Building and Running

### Prerequisites

- Go 1.25.3 or later
- CompileDaemon (for development)

### Build

```bash
go build -o odin main.go
```

### Run

```bash
./odin
```

### Development Mode (with hot reload)

```bash
make run
```

## Project Status

**Current Status**: Initial structure setup complete

### Completed

- ✅ Core directory structure
- ✅ State management types and operations
- ✅ Agent interfaces (MainAgent, SubAgent)
- ✅ Tool interface and registry system
- ✅ Planner types and interface
- ✅ Message handling and queuing logic
- ✅ File cache types and operations
- ✅ **Backend API layer with channel-based communication**
- ✅ **Complete TUI implementation with tview**
- ✅ **Message view, input field, todo list, status bar**
- ✅ **Frontend/Backend separation via Go channels**
- ✅ **Mode switching (Ctrl+A: Ask, Ctrl+E: Edit, Ctrl+P: Plan)**

### TODO

- [ ] Implement planner LLM integration
- [ ] Implement individual tools (LS, Grep, Glob, etc.)
- [ ] Add LLM-based cache validation for EditTool
- [ ] Implement message processing iteration loop
- [ ] Add configuration file support (odinconfig.json)
- [ ] Implement ODIN.md custom instructions
- [ ] Add Redis state publishing (optional)
- [ ] Implement ContextSummarizer tool
- [ ] Add WebFetch tool with HTML-to-markdown conversion

## Architecture Details

For detailed architecture information, see [design/DESIGN.md](design/DESIGN.md).

### Channel-Based Communication Flow

```
┌─────────────────┐           Channels            ┌─────────────────┐
│                 │  ─────────────────────────>   │                 │
│   UI (Client)   │    Request (string, mode)     │  Backend (API)  │
│   - tview       │                                │  - MainAgent    │
│   - components  │  <─────────────────────────   │  - State        │
│                 │   StateSnapshot (100ms)        │  - Planner      │
└─────────────────┘                                └─────────────────┘
        ↑                                                   │
        │                                                   │
        └───────────── UpdateChan (real-time) ────────────┘
```

### Message Processing Flow

1. **UI sends request** → User types message, hits Enter
2. **Request sent via channel** → `backend.RequestChan <- Request{Message, Mode}`
3. **Backend receives** → HandleIncomingMessage() called
4. **Check if busy** → If busy, add to queue; if idle, process immediately
5. **Processing** → Run iteration loop: `Planner → Tool Call → Planner`
6. **State updates** → Backend publishes snapshots every 100ms
7. **UI updates** → Receives snapshots, updates message view/todos/status
8. **Complete** → Process next queued message

### Tool System

Tools implement a three-phase lifecycle:

1. **PreHook**: Validation and preparation
2. **Execute**: Main tool logic
3. **PostHook**: Response formatting

## License

[Add your license here]

## Contributing

[Add contributing guidelines here]
