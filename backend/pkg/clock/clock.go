// Package clock provides a time abstraction for testability.
// Use Clock interface instead of calling time.Now() directly to enable mocking in tests.
package clock

import "time"

// Clock is an interface for getting the current time.
// This allows for easy mocking in tests.
type Clock interface {
	// Now returns the current time
	Now() time.Time
}

// RealClock is the production implementation that uses time.Now()
type RealClock struct{}

// Now returns the current time
func (RealClock) Now() time.Time {
	return time.Now()
}

// New returns a new RealClock instance
func New() Clock {
	return RealClock{}
}

// MockClock is a mock implementation for testing
type MockClock struct {
	fixedTime time.Time
}

// NewMock creates a MockClock with a fixed time
func NewMock(t time.Time) *MockClock {
	return &MockClock{fixedTime: t}
}

// Now returns the fixed time
func (m *MockClock) Now() time.Time {
	return m.fixedTime
}

// SetTime updates the fixed time
func (m *MockClock) SetTime(t time.Time) {
	m.fixedTime = t
}

// Advance moves the clock forward by the given duration
func (m *MockClock) Advance(d time.Duration) {
	m.fixedTime = m.fixedTime.Add(d)
}
