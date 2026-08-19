package thermal

import "testing"

func TestCapacityFactorNominal(t *testing.T) {
	p := DefaultParams()
	f := p.CapacityFactor(25)
	if f != 1.0 {
		t.Errorf("CapacityFactor(25) = %f, want 1.0", f)
	}
}

func TestCapacityFactorCold(t *testing.T) {
	p := DefaultParams()
	f := p.CapacityFactor(0)
	if f >= 1.0 || f <= 0 {
		t.Errorf("CapacityFactor(0) = %f, expected in (0, 1)", f)
	}
}

func TestResistanceFactorIncreasesCold(t *testing.T) {
	p := DefaultParams()
	r25 := p.ResistanceFactor(25)
	r0 := p.ResistanceFactor(0)
	if r0 <= r25 {
		t.Errorf("resistance at 0C (%f) should be > at 25C (%f)", r0, r25)
	}
}

func TestSimulationHeatsUp(t *testing.T) {
	sim := NewSimulation(25, 25, 1000, 1.0)
	sim.Step(100, 0.02, 0.1)
	if sim.CellTemp <= 25 {
		t.Errorf("cell should heat up under load, got %f", sim.CellTemp)
	}
}
