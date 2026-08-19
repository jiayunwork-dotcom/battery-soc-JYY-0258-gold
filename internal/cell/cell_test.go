package cell

import (
	"testing"

	"battery-soc/internal/coulomb"
)

func TestNewCellDefaults(t *testing.T) {
	c, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	st := c.State()
	if st.SOC != 50 {
		t.Errorf("SOC = %f, want 50", st.SOC)
	}
	if st.Voltage < 3.5 || st.Voltage > 4.0 {
		t.Errorf("Voltage = %f, expected 3.5-4.0", st.Voltage)
	}
}

func TestCellCharge(t *testing.T) {
	c, _ := New(DefaultConfig())
	c.Step(10, 1.0) // charge 10A for 1h
	st := c.State()
	if st.SOC <= 50 {
		t.Errorf("SOC should increase on charge, got %f", st.SOC)
	}
}

func TestPackMinSOC(t *testing.T) {
	p, err := NewPack(3, DefaultConfig())
	if err != nil {
		t.Fatalf("NewPack: %v", err)
	}
	// Discharge the pack
	p.Step(-10, 1.0)
	if p.MinSOC() >= 50 {
		t.Errorf("MinSOC should decrease on discharge, got %f", p.MinSOC())
	}
}

func TestPackRunProfile(t *testing.T) {
	p, _ := NewPack(2, DefaultConfig())
	samples := []coulomb.CurrentSample{
		{DT: 0.5, Current: 5},
		{DT: 0.5, Current: -5},
	}
	voltages := p.RunProfile(samples)
	if len(voltages) != 2 {
		t.Fatalf("got %d voltages, want 2", len(voltages))
	}
}
