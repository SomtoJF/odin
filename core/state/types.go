package state

import (
	"sync"
	"time"
)

// AgentMode represents the operational mode of the agent
type AgentMode string

const (
	AgentModeAsk  AgentMode = "ask_mode"
	AgentModePlan AgentMode = "plan_mode"
	AgentModeEdit AgentMode = "edit_mode"
)

// TodoStatus represents the status of a todo item
type TodoStatus string

const (
	TodoStatusCompleted  TodoStatus = "completed"
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
)

// EvictionPolicy defines cache eviction strategies
type EvictionPolicy string

const (
	EvictionPolicyLRU    EvictionPolicy = "lru"    // Least Recently Used
	EvictionPolicyLFU    EvictionPolicy = "lfu"    // Least Frequently Used
	EvictionPolicyHybrid EvictionPolicy = "hybrid" // Prefer evicting unmodified files
)

// State represents the global state of Odin Code
type State struct {
	AgentMode    AgentMode
	Messages     []Message
	MessageQueue []QueuedMessage // Queued messages when agent is busy
	IsExecuting  bool            // Tracks if agent is actively processing a message
	SubAgents    []SubAgent      // Currently running subagent
	Context      []ContextItem
	CustomInstructions string // Data from the ODIN.md
	Config             Config // includes configuration data specifying permissions for the agent

	// File Cache System - treats file data as key-value cache
	FileCache       map[string]*CachedFile // filepath -> cached file data
	FileCacheConfig FileCacheConfig        // Configuration for cache behavior

	// Mutexes for thread safety
	MessagesMx     sync.Mutex   // since state is a shared resource we will need to use a mutex to prevent RW race conditions
	MessageQueueMx sync.Mutex   // mutex for message queue
	SubAgentsMx    sync.Mutex   //
	StateMx        sync.Mutex   // since state is a shared resource we will need to use a mutex to prevent RW race conditions
	StdinMx        sync.Mutex   // Different concurrent processed might want to use the standard input at the same time
	FileCacheMx    sync.RWMutex // RWMutex for read-heavy file cache operations
}

// Message represents a single message exchange
type Message struct {
	Body          string        // The actual message sent by the user
	AnswerSummary string        // Final answer/result returned to user after iteration loop completes
	Todos         []Todo        // Can be empty especially if the message was sent in ask mode where it doesnt have access to the TodoWrite tool
	ToolHistory   []ToolHistory // A list of tools which were called to answer the user's message
	Updates       []string      // String array of realtime updates made to the CLI to inform the user what is going on
}

// QueuedMessage represents a message in the queue
type QueuedMessage struct {
	Body      string
	Mode      AgentMode // Mode to use when processing this message
	Timestamp time.Time
}

// Todo represents a task item
type Todo struct {
	ID      uint
	Status  TodoStatus // "pending", "in_progress" or "completed"
	Content string
}

// ToolHistory tracks tool execution history
type ToolHistory struct {
	ToolName           string
	AffectedFiles      []string
	ToolUseDescription string // Short description of the change that the tool made
}

// SubAgent represents a spawned subagent
type SubAgent struct {
	ID    uint
	Todos []Todo
}

// Config holds agent configuration
type Config struct {
	AllowedCommands   []string
	ForbiddenCommands []string
}

// ContextItem represents a piece of context
type ContextItem struct {
	Content       string
	FilePath      *string // context might come from a command result
	SourceCommand *string // source command for the data if the data is a result from a cli command
}

// CachedFile represents a file in the cache
type CachedFile struct {
	FilePath string // Absolute file path (serves as the key)

	// Content storage - structure determined by actual reads, not file size
	FullContent  *string                     // Full file contents (when tool reads entire file)
	PartialCache map[string]*CachedSegment   // Line range segments (when tool reads specific ranges)
	IsPartial    bool                        // True if PartialCache is used, false if FullContent

	// File metadata
	ContentHash string // SHA-256 hash for integrity checking (full file)
	Size        int64  // Total file size in bytes
	TotalLines  int    // Total number of lines in file

	// Timestamps for cache management
	CachedAt     time.Time // When file was first cached
	LastAccessed time.Time // For LRU eviction
	AccessCount  int       // Usage tracking
	ModTime      time.Time // File modification time from filesystem

	// State tracking
	IsModified   bool   // Track if file was written in this session
	OriginalHash string // Hash at first read (for dirty detection)
	IsTruncated  bool   // Whether content was truncated for size
}

// CachedSegment represents a cached portion of a file
type CachedSegment struct {
	StartLine int       // Starting line number (1-based)
	EndLine   int       // Ending line number (inclusive)
	Content   string    // The actual content of this segment
	Hash      string    // Hash of this segment
	CachedAt  time.Time // When this segment was cached
}

// FileCacheConfig holds cache configuration
type FileCacheConfig struct {
	MaxCacheSize      int64          // Total cache size in bytes (e.g., 100MB)
	MaxFileSize       int64          // Max size per file (e.g., 10MB)
	MaxEntries        int            // Max number of cached files (e.g., 500)
	EvictionPolicy    EvictionPolicy // LRU, LFU, or Hybrid
	TTL               time.Duration  // Optional: expire after duration
	EnableAutoRefresh bool           // Re-read if mod time changed

	// Partial caching settings - cache structure mirrors actual reads
	MaxSegments int // Max segments per file (e.g., 50)
}

// CachedFileInfo provides cache metadata to planner
type CachedFileInfo struct {
	FilePath    string    // Absolute file path (the cache key)
	Size        int64     // File size in bytes
	IsModified  bool      // Whether file was modified in this session
	CachedAt    time.Time // When file was cached
	IsTruncated bool      // Whether content was truncated for size
}
