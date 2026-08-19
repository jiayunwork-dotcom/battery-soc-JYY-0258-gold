package coulomb

import "testing"

// TestCoulombSOC_ChargeDischarge exercises single-step integration and clamping
// (other category). Charge/discharge move SOC; over/under-charge clamps to [0,100].
func TestCoulombSOC_ChargeDischarge(t *testing.T) {
	// 50 Ah battery, charge 5 Ah in 1h from SOC 50 -> 60.
	if got := CoulombSOC(50, 50, 5, 1); got != 60 {
		t.Fatalf("expected 60, got %v", got)
	}
	// Discharge 50 Ah from 50 in 1h -> 0 (clamped, not negative).
	if got := CoulombSOC(50, 50, -50, 1); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
	// Overcharge clamps to 100.
	if got := CoulombSOC(50, 95, 50, 1); got != 100 {
		t.Fatalf("expected 100, got %v", got)
	}
}

// TestCoulombFromLog_Slice integrates a sample series and verifies the empty/nil
// slice returns the starting SOC of 50 (slice category).
func TestCoulombFromLog_Slice(t *testing.T) {
	samples := []CurrentSample{
		{DT: 1, Current: 5},   // +10 -> 60
		{DT: 1, Current: -5},  // -10 -> 50
		{DT: 1, Current: -50}, // -100 -> 0 (clamped)
	}
	if got := CoulombFromLog(50, samples); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
	// nil/empty slice returns the starting SOC of 50.
	if empty := CoulombFromLog(50, nil); empty != 50 {
		t.Fatalf("expected 50 for empty slice, got %v", empty)
	}
}
