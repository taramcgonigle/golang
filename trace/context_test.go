package trace

import (
	"context"
	"testing"
)

// TestNewContextWithTrace verifies that a new context with a trace ID is created correctly.
func TestNewContextWithTrace(t *testing.T) {
	ctx := NewContextWithTrace(context.Background()) // Start with a base context
	traceID := GetTraceID(ctx)                       // Extract the trace ID from context

	// The trace ID should not be empty
	if traceID == "" {
		t.Fatal("Expected a trace ID, got empty string")
	}
}

// TestGetTraceID verifies that calling this multiple times on the same context retrieves the same  trace ID
func TestTraceIDIsConsistent(t *testing.T) {
	ctx := NewContextWithTrace(context.Background()) // Create a traceable context

	// Retrieve the trace ID twice
	id1 := GetTraceID(ctx)
	id2 := GetTraceID(ctx)

	if id1 != id2 { // They should be identical — no new ID should be generated
		t.Errorf("Expected consistent TraceID, got %s and %s", id1, id2)
	}
}
