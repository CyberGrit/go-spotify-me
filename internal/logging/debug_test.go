package logging

import (
	"testing"
)

func TestDebugLog(t *testing.T) {
	// Test when Debug is false
	Debug = false
	DebugLog("Test %s", "message") // Should not panic or print anything visible

	// Test when Debug is true
	Debug = true
	// We cannot easily assert the stdout output here without mocking os.Stdout
	// or zap, but we can assure it doesn't panic when calling the function.
	DebugLog("Test %s", "message")

	// Reset to false to avoid side effects in other tests if any
	Debug = false
}
