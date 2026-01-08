package main

import (
	"log"

	"github.com/SomtoJF/odin/core/api"
	components "github.com/SomtoJF/odin/ui"
)

func main() {
	// Create and start backend
	backend := api.NewBackend()
	backend.Start()

	// Create UI with injected channels
	ui := components.NewUI(
		backend.RequestChan,
		backend.ResponseChan,
		backend.UpdateChan,
	)

	// Run UI (blocks until quit)
	if err := ui.Run(); err != nil {
		log.Fatalf("Failed to run UI: %v", err)
	}

	// Graceful shutdown
	backend.Stop()
}
