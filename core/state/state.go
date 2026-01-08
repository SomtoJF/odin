package state

import (
	"fmt"
	"time"
)

// NewState creates a new State with default configuration
func NewState() *State {
	return &State{
		Messages:        []Message{},
		MessageQueue:    []QueuedMessage{},
		IsExecuting:     false,
		SubAgents:       []SubAgent{},
		Context:         []ContextItem{},
		FileCache:       make(map[string]*CachedFile),
		FileCacheConfig: NewDefaultFileCacheConfig(),
	}
}

// NewDefaultFileCacheConfig returns default cache configuration
func NewDefaultFileCacheConfig() FileCacheConfig {
	return FileCacheConfig{
		MaxCacheSize:      100 * 1024 * 1024,    // 100MB
		MaxFileSize:       10 * 1024 * 1024,     // 10MB per file
		MaxEntries:        500,                  // 500 files max
		EvictionPolicy:    EvictionPolicyHybrid, // Prefer evicting unmodified
		TTL:               0,                    // No expiration by default
		EnableAutoRefresh: true,                 // Auto-detect file changes
		MaxSegments:       50,                   // Max 50 segments per file
	}
}

// WriteMessageToState adds a new message to the state
func (s *State) WriteMessageToState(messageBody string) {
	s.MessagesMx.Lock()
	defer s.MessagesMx.Unlock()

	newMessage := Message{Body: messageBody}
	s.Messages = append(s.Messages, newMessage)
}

// GetCachedFile retrieves a file from cache (thread-safe)
func (s *State) GetCachedFile(filePath string) (*CachedFile, bool) {
	s.FileCacheMx.RLock()
	defer s.FileCacheMx.RUnlock()

	cached, exists := s.FileCache[filePath]
	return cached, exists
}

// UpdateCacheAccess updates access tracking for a cache entry
func (s *State) UpdateCacheAccess(filePath string) {
	s.FileCacheMx.Lock()
	defer s.FileCacheMx.Unlock()

	if cached, exists := s.FileCache[filePath]; exists {
		cached.LastAccessed = time.Now()
		cached.AccessCount++
		s.FileCache[filePath] = cached
	}
}

// GetCachedSegment retrieves a specific segment from partial cache
func (s *State) GetCachedSegment(filePath string, startLine, endLine int) (*CachedSegment, bool) {
	s.FileCacheMx.RLock()
	defer s.FileCacheMx.RUnlock()

	cached, exists := s.FileCache[filePath]
	if !exists || !cached.IsPartial {
		return nil, false
	}

	key := fmt.Sprintf("%d-%d", startLine, endLine)
	segment, found := cached.PartialCache[key]
	return segment, found
}

// AddCachedSegment adds a segment to partial cache
func (s *State) AddCachedSegment(filePath string, segment *CachedSegment) {
	s.FileCacheMx.Lock()
	defer s.FileCacheMx.Unlock()

	cached, exists := s.FileCache[filePath]
	if !exists {
		// Create new partial cache entry
		cached = &CachedFile{
			FilePath:     filePath,
			IsPartial:    true,
			PartialCache: make(map[string]*CachedSegment),
			CachedAt:     time.Now(),
			LastAccessed: time.Now(),
			AccessCount:  1,
		}
		s.FileCache[filePath] = cached
	}

	if cached.PartialCache == nil {
		cached.PartialCache = make(map[string]*CachedSegment)
	}

	// Add segment with range key
	key := fmt.Sprintf("%d-%d", segment.StartLine, segment.EndLine)
	cached.PartialCache[key] = segment

	// Evict oldest segments if over limit
	if len(cached.PartialCache) > s.FileCacheConfig.MaxSegments {
		s.evictOldestSegment(cached)
	}
}

// evictOldestSegment removes oldest segment from partial cache
func (s *State) evictOldestSegment(cached *CachedFile) {
	var oldestKey string
	oldestTime := time.Now()

	for key, segment := range cached.PartialCache {
		if segment.CachedAt.Before(oldestTime) {
			oldestTime = segment.CachedAt
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(cached.PartialCache, oldestKey)
	}
}

// NeedsEviction checks if cache needs eviction before adding new file
func (s *State) NeedsEviction(newFileSize int64) bool {
	config := s.FileCacheConfig

	// Check entry count
	if len(s.FileCache) >= config.MaxEntries {
		return true
	}

	// Check total size
	currentSize := int64(0)
	for _, cached := range s.FileCache {
		currentSize += cached.Size
	}

	return currentSize+newFileSize > config.MaxCacheSize
}

// EvictLRU evicts least recently used file
func (s *State) EvictLRU() {
	s.FileCacheMx.Lock()
	defer s.FileCacheMx.Unlock()

	var oldestPath string
	oldestTime := time.Now()

	for path, cached := range s.FileCache {
		if cached.LastAccessed.Before(oldestTime) {
			oldestTime = cached.LastAccessed
			oldestPath = path
		}
	}

	if oldestPath != "" {
		delete(s.FileCache, oldestPath)
	}
}

// EvictHybrid prefers evicting unmodified files, then uses LRU
func (s *State) EvictHybrid() {
	s.FileCacheMx.Lock()
	defer s.FileCacheMx.Unlock()

	var targetPath string
	oldestUnmodified := time.Now()
	oldestModified := time.Now()
	var hasUnmodified bool
	var oldestModifiedPath string

	// Find oldest unmodified and modified files
	for path, cached := range s.FileCache {
		if !cached.IsModified {
			hasUnmodified = true
			if cached.LastAccessed.Before(oldestUnmodified) {
				oldestUnmodified = cached.LastAccessed
				targetPath = path
			}
		} else {
			if cached.LastAccessed.Before(oldestModified) {
				oldestModified = cached.LastAccessed
				oldestModifiedPath = path
			}
		}
	}

	// Prefer evicting unmodified files
	if !hasUnmodified && oldestModifiedPath != "" {
		targetPath = oldestModifiedPath
	}

	if targetPath != "" {
		delete(s.FileCache, targetPath)
	}
}
