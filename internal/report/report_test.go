package report

import (
	"bytes"
	"testing"
)

func TestBuildSummary(t *testing.T) {
	trace := []float64{50, 60, 55, 40}
	times := []float64{0, 1, 2, 3}
	s := BuildSummary(trace, times, 10.0, "coulomb")
	if s.InitialSOC != 50 {
		t.Errorf("InitialSOC = %f, want 50", s.InitialSOC)
	}
	if s.FinalSOC != 40 {
		t.Errorf("FinalSOC = %f, want 40", s.FinalSOC)
	}
	if s.MinSOC != 40 {
		t.Errorf("MinSOC = %f, want 40", s.MinSOC)
	}
	if s.MaxSOC != 60 {
		t.Errorf("MaxSOC = %f, want 60", s.MaxSOC)
	}
}

func TestWriteText(t *testing.T) {
	r := &FullReport{
		Summary: Summary{Method: "coulomb", FinalSOC: 45, Steps: 10},
	}
	var buf bytes.Buffer
	err := WriteText(&buf, r)
	if err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
