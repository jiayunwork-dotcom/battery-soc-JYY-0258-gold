package profile

import "testing"

func TestCCProfile(t *testing.T) {
	p := CC(10, 0.5, 4)
	if len(p) != 4 {
		t.Fatalf("CC: got %d steps, want 4", len(p))
	}
	if p[0].Current != 10 {
		t.Errorf("CC current = %f, want 10", p[0].Current)
	}
}

func TestCCCVTaper(t *testing.T) {
	p := CCCV(10, 3, 3, 0.1)
	if len(p) != 6 {
		t.Fatalf("CCCV: got %d steps, want 6", len(p))
	}
	// Last CV step should be near zero
	if p[5].Current >= p[0].Current {
		t.Error("CV phase should taper current")
	}
}

func TestDriveCycle(t *testing.T) {
	p := DriveCycle(100, 0.01, 2)
	if len(p) != 20 {
		t.Fatalf("DriveCycle: got %d steps, want 20", len(p))
	}
}

func TestAnalyze(t *testing.T) {
	p := CC(-50, 1.0, 3)
	a := Analyze(p)
	if a.PeakDischarge != -50 {
		t.Errorf("PeakDischarge = %f, want -50", a.PeakDischarge)
	}
	if a.DurationH != 3.0 {
		t.Errorf("DurationH = %f, want 3.0", a.DurationH)
	}
}
