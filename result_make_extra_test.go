package opera

import (
	"errors"
	"testing"
)

func TestTryYield(t *testing.T) {
	if got := TryYield(42, nil); got != 42 {
		t.Fatalf("TryYield() = %d, want 42", got)
	}

	wantErr := errors.New("try yield failed")
	result := Do(func() int {
		return TryYield(42, wantErr)
	})
	if !errors.Is(result.Err(), wantErr) {
		t.Fatalf("TryYield() error = %v, want %v", result.Err(), wantErr)
	}
}

func TestTryYieldUnit(t *testing.T) {
	if got := TryYieldUnit(nil); got != U {
		t.Fatalf("TryYieldUnit() = %v, want Unit", got)
	}

	wantErr := errors.New("unit yield failed")
	result := Do(func() Unit {
		return TryYieldUnit(wantErr)
	})
	if !errors.Is(result.Err(), wantErr) {
		t.Fatalf("TryYieldUnit() error = %v, want %v", result.Err(), wantErr)
	}
}

func TestYieldErr(t *testing.T) {
	wantErr := errors.New("yield err failed")
	result := Do(func() int {
		if wantErr != nil {
			return YieldErr[int](wantErr)
		}
		return 42
	})
	if !errors.Is(result.Err(), wantErr) {
		t.Fatalf("YieldErr() error = %v, want %v", result.Err(), wantErr)
	}
}

func TestYieldErrWithNilError(t *testing.T) {
	defer func() {
		if recovered := recover(); !errors.Is(recovered.(error), ErrNilYieldError) {
			t.Fatalf("YieldErr(nil) panic = %v, want %v", recovered, ErrNilYieldError)
		}
	}()

	YieldErr[int](nil)
}
