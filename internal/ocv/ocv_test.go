package ocv

import "testing"

func TestVoltageAtSOC(t *testing.T) {
	tbl := DefaultNMC()
	v := tbl.VoltageAtSOC(50)
	if v < 3.7 || v > 3.8 {
		t.Errorf("VoltageAtSOC(50) = %f, expected ~3.73", v)
	}
}

func TestSOCAtVoltage(t *testing.T) {
	tbl := DefaultNMC()
	soc := tbl.SOCAtVoltage(3.73)
	if soc < 45 || soc > 55 {
		t.Errorf("SOCAtVoltage(3.73) = %f, expected ~50", soc)
	}
}

func TestMonotonic(t *testing.T) {
	tbl := DefaultNMC()
	if !Monotonic(tbl) {
		t.Error("DefaultNMC should be monotonic")
	}
}

func TestCalibrateTable(t *testing.T) {
	pts := []CalibrationPoint{
		{SOC: 0, Voltage: 3.0},
		{SOC: 50, Voltage: 3.7},
		{SOC: 100, Voltage: 4.2},
	}
	tbl, err := CalibrateTable(pts)
	if err != nil {
		t.Fatalf("CalibrateTable: %v", err)
	}
	v := tbl.VoltageAtSOC(25)
	if v < 3.3 || v > 3.5 {
		t.Errorf("interpolated V(25%%) = %f", v)
	}
}
