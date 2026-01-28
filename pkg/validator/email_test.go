package validator

import (
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.com", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"user@", false},
		{"", false},
		{"a@b", false},
		{"a@b.c", false}, // TLD must be at least 2 chars
		{"a@b.co", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v, expected %v", tt.email, got, tt.expected)
			}
		})
	}
}
