package parse

import (
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	input := "dt,current\n0.5,10\n1.0,-5\n"
	samples, err := ReadCSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadCSV: %v", err)
	}
	if len(samples) != 2 {
		t.Fatalf("got %d samples, want 2", len(samples))
	}
	if samples[0].DT != 0.5 || samples[0].Current != 10 {
		t.Errorf("sample[0] = %+v", samples[0])
	}
}

func TestReadJSON(t *testing.T) {
	input := `[{"dt":0.5,"current":10},{"dt":1.0,"current":-5}]`
	samples, err := ReadJSON(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}
	if len(samples) != 2 {
		t.Fatalf("got %d samples, want 2", len(samples))
	}
}

func TestComputeStats(t *testing.T) {
	input := "dt,current\n1.0,10\n1.0,-5\n"
	samples, _ := ReadCSV(strings.NewReader(input))
	stats := ComputeStats(samples)
	if stats.Count != 2 {
		t.Errorf("Count = %d, want 2", stats.Count)
	}
	if stats.TotalChargeAh != 5 {
		t.Errorf("TotalChargeAh = %f, want 5", stats.TotalChargeAh)
	}
}
