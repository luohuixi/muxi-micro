package errs

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrBasicFunctionality(t *testing.T) {
	t.Run("NewErr creates valid error", func(t *testing.T) {
		e := NewErr("AUTH", "authentication failed")
		if e.Name() != "AUTH" {
			t.Errorf("expected name 'AUTH', got '%s'", e.Name())
		}
		if e.Message() != "authentication failed" {
			t.Errorf("expected message 'authentication failed', got '%s'", e.Message())
		}
		if e.MetaData() != nil {
			t.Error("new error should have nil metadata")
		}
		if e.Unwrap() != nil {
			t.Error("new error should have nil cause")
		}
	})

	t.Run("Error message formatting", func(t *testing.T) {
		e := NewErr("DB", "connection timeout")
		if e.Error() != "[DB] connection timeout" {
			t.Errorf("unexpected error format: %s", e.Error())
		}

		e = e.WithMeta(map[string]interface{}{"host": "db1.example.com", "port": 5432})
		expected := "[DB] connection timeout map[host:db1.example.com port:5432]"
		if e.Error() != expected {
			t.Errorf("unexpected error with metadata: %s", e.Error())
		}

		e = e.WithCause(errors.New("network unreachable"))
		if e.Error() != expected+" => network unreachable" {
			t.Errorf("unexpected error with cause: %s", e.Error())
		}
	})
}

func TestWithMeta(t *testing.T) {
	baseErr := NewErr("VALIDATION", "invalid input")

	t.Run("Add metadata to error", func(t *testing.T) {
		metaErr := baseErr.WithMeta(map[string]interface{}{
			"field": "email",
			"value": "user@example",
		})

		if metaErr.MetaData()["field"] != "email" {
			t.Error("metadata not set correctly")
		}
		if metaErr.Unwrap() != nil {
			t.Error("WithMeta should not alter cause")
		}
	})

	t.Run("Metadata doesn't affect original error", func(t *testing.T) {
		_ = baseErr.WithMeta(map[string]interface{}{"test": "value"})
		if baseErr.MetaData() != nil {
			t.Error("original error metadata should remain nil")
		}
	})

	t.Run("Nil error handling", func(t *testing.T) {
		var nilErr *Err
		if nilErr.WithMeta(nil) != nil {
			t.Error("WithMeta on nil error should return nil")
		}
	})
}

func TestWithCause(t *testing.T) {
	rootErr := errors.New("root error")
	midErr := NewErr("MID", "mid error").WithCause(rootErr)
	topErr := NewErr("TOP", "top error").WithCause(midErr)

	t.Run("Error chain unwrapping", func(t *testing.T) {
		if !errors.Is(topErr.Unwrap(), midErr) {
			t.Error("Unwrap should return direct cause")
		}

		if !errors.Is(errors.Unwrap(topErr.Unwrap()), rootErr) {
			t.Error("Double unwrap should return root error")
		}
	})

	t.Run("Chained error formatting", func(t *testing.T) {
		expected := "[TOP] top error => [MID] mid error => root error"
		if topErr.Error() != expected {
			t.Errorf("chained error format mismatch:\nExpected: %s\nGot:      %s",
				expected, topErr.Error())
		}
	})
}

func TestIsMethod(t *testing.T) {
	dbErr := NewErr("DB", "connection failed")
	authErr := NewErr("AUTH", "unauthorized")
	wrappedErr := dbErr.WithCause(authErr)

	t.Run("Direct comparison", func(t *testing.T) {
		if !errors.Is(dbErr, dbErr) {
			t.Error("errors.Is should match same error")
		}
	})

	t.Run("Match in chain", func(t *testing.T) {
		if !errors.Is(wrappedErr, dbErr) {
			t.Error("errors.Is should match top error")
		}
		if !errors.Is(wrappedErr, authErr) {
			t.Error("errors.Is should match cause error")
		}
	})

	t.Run("Non-match cases", func(t *testing.T) {
		if errors.Is(dbErr, authErr) {
			t.Error("different errors should not match")
		}

		otherDbErr := NewErr("DB", "different message")
		if !errors.Is(dbErr, otherDbErr) {
			t.Error("same name should match regardless of message")
		}

		if errors.Is(wrappedErr, errors.New("random error")) {
			t.Error("should not match unrelated error")
		}
	})
}

func TestAsMethod(t *testing.T) {
	origErr := NewErr("API", "request failed").WithMeta(map[string]interface{}{
		"status": 500,
		"path":   "/users",
	})

	wrapped := fmt.Errorf("wrapper: %w", origErr)

	t.Run("Extract Err type", func(t *testing.T) {
		var target *Err
		if !errors.As(wrapped, &target) {
			t.Fatal("errors.As failed to extract Err type")
		}

		if target.Name() != "API" {
			t.Errorf("expected name 'API', got '%s'", target.Name())
		}
		if target.MetaData()["status"] != 500 {
			t.Error("metadata not preserved")
		}
	})

	t.Run("Extract from chain", func(t *testing.T) {
		chainErr := NewErr("TOP", "top error").WithCause(origErr)

		var target *Err
		if !errors.As(chainErr, &target) {
			t.Fatal("errors.As failed to extract from chain")
		}

		if !errors.Is(target, origErr) {
			t.Error("should extract original error from chain")
		}
	})
}

func TestEdgeCases(t *testing.T) {

	t.Run("Empty metadata", func(t *testing.T) {
		e := NewErr("TEST", "test").WithMeta(map[string]interface{}{})
		if e.MetaData() == nil {
			t.Error("empty map should be preserved")
		}
		if len(e.MetaData()) != 0 {
			t.Error("metadata map should be empty")
		}
	})

	t.Run("Deep error chain", func(t *testing.T) {
		root := errors.New("level 0")
		var current error = root

		// Create 10-level error chain
		for i := 1; i <= 10; i++ {
			current = NewErr("LVL", "level error").WithCause(current)
		}

		// Verify we can find root error
		if !errors.Is(current, root) {
			t.Error("should find root error in deep chain")
		}

		// Verify chain length in error message
		expectedParts := 11 // 10 levels + root
		errorStr := current.Error()
		count := 0
		for i := 0; i < len(errorStr); i++ {
			if errorStr[i] == '=' && i+1 < len(errorStr) && errorStr[i+1] == '>' {
				count++
			}
		}
		if count != expectedParts-1 {
			t.Errorf("expected %d chain separators, got %d", expectedParts-1, count)
		}
	})
}
