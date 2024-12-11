package pkg

import (
	"fmt"
	"testing"
)

func TestResultSuccess(t *testing.T) {
	result := Wrap(1, nil)
	result.OnSuccess(func(i int) {
		if i != 1 {
			t.Errorf("Expected 1, got %d", i)
		}
	})
	result.OnFailed(func(err error) {
		t.Errorf("Expected no error, got %v", err)
	})
	if result.Value() != 1 {
		t.Errorf("Expected 1, got %d", result.Value())
	}
	if result.Error() != nil {
		t.Errorf("Expected no error, got %v", result.Error())
	}
	if result.ValueOrDefault(2) != 1 {
		t.Errorf("Expected 1, got %d", result.ValueOrDefault(2))
	}
}

func TestResultFailure(t *testing.T) {
	ErrTest := fmt.Errorf("test error")
	result := Wrap(1, ErrTest)
	result.OnSuccess(func(i int) {
		t.Errorf("Expected error, got %d", i)
	})
	result.OnFailed(func(err error) {
		if err != ErrTest {
			t.Errorf("Expected ErrTest, got %v", err)
		}
	})
	if result.Value() != 1 {
		t.Errorf("Expected 1, got %d", result.Value())
	}
	if result.Error() != ErrTest {
		t.Errorf("Expected ErrTest, got %v", result.Error())
	}
}
