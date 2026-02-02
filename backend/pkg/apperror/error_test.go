package apperror

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorContext(t *testing.T) {
	baseErr := errors.New("base error")

	t.Run("New", func(t *testing.T) {
		err := New(CodeInternal, "message")
		if err.Code != CodeInternal {
			t.Errorf("Expected code %s, got %s", CodeInternal, err.Code)
		}
		if err.Message != "message" {
			t.Errorf("Expected message %s, got %s", "message", err.Message)
		}
	})

	t.Run("Wrap", func(t *testing.T) {
		err := Wrap(CodeDatabaseError, "db failed", baseErr)
		if err.Cause != baseErr {
			t.Error("Expected cause to be baseErr")
		}
		if !errors.Is(err, baseErr) {
			t.Error("Wrapper error should match cause via Unwrap")
		}
	})

	t.Run("WithDetails", func(t *testing.T) {
		err := New(CodeInternal, "").WithDetails("key", "value")
		if err.Details["key"] != "value" {
			t.Errorf("Expected detail value 'value', got %v", err.Details["key"])
		}
	})

	t.Run("WithField", func(t *testing.T) {
		err := New(CodeInvalidInput, "").WithField("email")
		if err.Field != "email" {
			t.Errorf("Expected field 'email', got %s", err.Field)
		}
		if err.Error() != fmt.Sprintf("[%s] email: ", CodeInvalidInput) {
			t.Errorf("Unexpected error string: %s", err.Error())
		}
	})
}

func TestHelpers(t *testing.T) {
	t.Run("IsCode", func(t *testing.T) {
		err := NotFound("user")
		if !IsCode(err, CodeNotFound) {
			t.Error("Expected IsCode to return true")
		}
		if IsCode(errors.New("other"), CodeNotFound) {
			t.Error("Expected IsCode to return false for non-app error")
		}
	})

	t.Run("IsNotFound", func(t *testing.T) {
		if !IsNotFound(NotFound("x")) {
			t.Error("IsNotFound failed")
		}
	})

	t.Run("GetCode", func(t *testing.T) {
		if GetCode(AuthorizedError()) != CodeUnauthorized {
			t.Error("GetCode failed")
		}
		if GetCode(errors.New("other")) != CodeInternal {
			t.Error("GetCode should default to Internal")
		}
	})
}

func AuthorizedError() error {
	return Unauthorized("test")
}
