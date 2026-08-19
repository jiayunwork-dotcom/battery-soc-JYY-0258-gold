package soh

import (
	"math"
	"testing"
)

// TestSOH_Error verifies an error is returned for non-positive initial capacity
// (error category).
func TestSOH_Error(t *testing.T) {
	if _, err := SOH(0, 40); err == nil {
		t.Fatal("expected error for initialCap=0")
	}
	if _, err := SOH(-5, 40); err == nil {
		t.Fatal("expected error for negative initialCap")
	}
}

// TestSOH_Ratio verifies the health ratio computation (other category).
func TestSOH_Ratio(t *testing.T) {
	got, err := SOH(100, 80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(got-0.8) > 1e-9 {
		t.Fatalf("expected 0.8, got %v", got)
	}
	// Ratio is clamped to 1 when currentCap exceeds initialCap.
	if c, _ := SOH(100, 120); c != 1 {
		t.Fatalf("expected clamp to 1, got %v", c)
	}
}

// TestCycleLife_Error verifies an error for non-positive rated throughput
// (error category).
func TestCycleLife_Error(t *testing.T) {
	if _, err := CycleLife(10, 0); err == nil {
		t.Fatal("expected error for ratedThroughputAh=0")
	}
	if _, err := CycleLife(10, -1); err == nil {
		t.Fatal("expected error for negative ratedThroughputAh")
	}
}

// TestDegradationModel_Retention verifies monotonic capacity loss and the 0 floor
// (other category).
func TestDegradationModel_Retention(t *testing.T) {
	if DegradationModel(0) != 1 {
		t.Fatal("expected 1.0 at 0 cycles")
	}
	if math.Abs(DegradationModel(1000)-0.8) > 1e-9 {
		t.Fatalf("expected 0.8 at 1000 cycles, got %v", DegradationModel(1000))
	}
	if DegradationModel(100000) != 0 {
		t.Fatal("expected 0 floor at very high cycles")
	}
}
