package telemetry

import (
	"context"
	"testing"
)

func TestSetupOTelSDK(t *testing.T) {
	// Skip actual connection test in unit tests as it requires a running collector
	// Just verify function signature and compilation

	// Create a context
	ctx := context.Background()

	// Check if setup function exists and is callable with bad endpoint (should fail or return shutdown)
	// Note: We use a dummy endpoint that won't connect, but Setup should likely proceed
	// until exporter starts, which might happen in background.
	// Actually, New() for exporters might return nil immediately or try to connect.
	// GRPC usually connects lazily or fails if dial fails.

	// To safely test without external dependencies, we'd need to mock exports, which is hard with this implementation.
	// For now, we will trust the library code and just check if we can call it without panic (smoke test).

	shutdown, err := SetupOTelSDK(ctx, "test-service", "1.0.0", "localhost:4317")

	// Even if it fails to connect to collector,err might be nil if it's async,
	// or it might error if Dial fails.
	// We just ensure it doesn't panic.

	if shutdown != nil {
		_ = shutdown(ctx)
	}

	if err != nil {
		// It's acceptable for it to fail if no collector is running
		t.Logf("Setup failed as expected (no collector): %v", err)
	}
}
