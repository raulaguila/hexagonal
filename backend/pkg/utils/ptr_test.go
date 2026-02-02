package utils

import (
	"testing"
)

func TestDeref(t *testing.T) {
	val := "test"
	ptr := &val

	result := Deref(ptr, "default")
	if result != "test" {
		t.Errorf("Expected 'test', got %v", result)
	}

	var nilPtr *string
	result = Deref(nilPtr, "default")
	if result != "default" {
		t.Errorf("Expected 'default', got %v", result)
	}

	intVal := 123
	intPtr := &intVal
	intResult := Deref(intPtr, 0)
	if intResult != 123 {
		t.Errorf("Expected 123, got %v", intResult)
	}

	var nilIntPtr *int
	intResult = Deref(nilIntPtr, 456)
	if intResult != 456 {
		t.Errorf("Expected 456, got %v", intResult)
	}
}
