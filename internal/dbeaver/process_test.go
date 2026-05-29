package dbeaver

import "testing"

// TestIsRunning verifies that IsRunning() executes without panicking.
// The actual process detection depends on system state and is verified manually.
// This test ensures the function is callable and returns a boolean.
func TestIsRunning(t *testing.T) {
	// IsRunning must not panic on any platform.
	// The result is non-deterministic (depends on whether DBeaver is open),
	// but must be a boolean — no crash is the contract being tested here.
	result := IsRunning()
	// result is either true or false — both are valid.
	// We verify the function completes without error.
	_ = result
	t.Logf("IsRunning() = %v (system-dependent)", result)
}
