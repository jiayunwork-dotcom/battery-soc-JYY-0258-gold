package kalman

import "testing"

// TestNewEKF_Init verifies the filter initializes SOC and a positive covariance
// (other category).
func TestNewEKF_Init(t *testing.T) {
	k := NewEKF(42)
	if k == nil {
		t.Fatal("NewEKF returned nil")
	}
	if k.SOC != 42 {
		t.Fatalf("expected SOC 42, got %v", k.SOC)
	}
	if k.P <= 0 {
		t.Fatalf("expected positive P, got %v", k.P)
	}
}

// TestEKFStep_Correction verifies a step pulls SOC toward the measurement z and
// updates k.SOC (other category). With current 0 the prediction stays at 50, so
// the corrected value must lie strictly between prediction and z.
func TestEKFStep_Correction(t *testing.T) {
	k := NewEKF(50)
	got := k.Step(70, 0, 1, 50)
	if got <= 50 || got >= 70 {
		t.Fatalf("expected correction strictly between 50 and 70, got %v", got)
	}
	if k.SOC != got {
		t.Fatalf("k.SOC not updated: %v != %v", k.SOC, got)
	}
	// Clamping: a huge negative measurement must not push SOC below 0.
	k2 := NewEKF(10)
	if out := k2.Step(-1000, 0, 1, 50); out != 0 {
		t.Fatalf("expected clamp to 0, got %v", out)
	}
}
