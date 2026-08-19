// Package report generates structured analysis reports from battery
// simulation or estimation results. Reports can be written as plain text
// or JSON for downstream consumption.
package report

import (
	"encoding/json"
	"fmt"
	"io"
)

// Summary holds the key metrics from a SOC estimation run.
type Summary struct {
	InitialSOC    float64 `json:"initial_soc"`
	FinalSOC      float64 `json:"final_soc"`
	MinSOC        float64 `json:"min_soc"`
	MaxSOC        float64 `json:"max_soc"`
	TotalChargeAh float64 `json:"total_charge_ah"`
	TotalTimeH    float64 `json:"total_time_h"`
	Steps         int     `json:"steps"`
	Method        string  `json:"method"`
}

// SOCTrace records SOC at each time step for plotting or analysis.
type SOCTrace struct {
	Time []float64 `json:"time_h"`
	SOC  []float64 `json:"soc_pct"`
}

// FullReport bundles the summary and optional trace.
type FullReport struct {
	Summary Summary   `json:"summary"`
	Trace   *SOCTrace `json:"trace,omitempty"`
}

// WriteText writes a human-readable report to w.
func WriteText(w io.Writer, r *FullReport) error {
	lines := []string{
		"Battery SOC Estimation Report",
		fmt.Sprintf("  Method:         %s", r.Summary.Method),
		fmt.Sprintf("  Steps:          %d", r.Summary.Steps),
		fmt.Sprintf("  Total Time:     %.2f h", r.Summary.TotalTimeH),
		fmt.Sprintf("  Initial SOC:    %.2f %%", r.Summary.InitialSOC),
		fmt.Sprintf("  Final SOC:      %.2f %%", r.Summary.FinalSOC),
		fmt.Sprintf("  Min SOC:        %.2f %%", r.Summary.MinSOC),
		fmt.Sprintf("  Max SOC:        %.2f %%", r.Summary.MaxSOC),
		fmt.Sprintf("  Total Charge:   %.2f Ah", r.Summary.TotalChargeAh),
		"",
	}
	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes the report as indented JSON.
func WriteJSON(w io.Writer, r *FullReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// BuildSummary constructs a Summary from a SOC trace.
func BuildSummary(trace []float64, times []float64, chargeAh float64, method string) Summary {
	if len(trace) == 0 {
		return Summary{Method: method}
	}
	s := Summary{
		InitialSOC:    trace[0],
		FinalSOC:      trace[len(trace)-1],
		MinSOC:        trace[0],
		MaxSOC:        trace[0],
		TotalChargeAh: chargeAh,
		Steps:         len(trace),
		Method:        method,
	}
	for _, v := range trace {
		if v < s.MinSOC {
			s.MinSOC = v
		}
		if v > s.MaxSOC {
			s.MaxSOC = v
		}
	}
	if len(times) > 0 {
		s.TotalTimeH = times[len(times)-1]
	}
	return s
}
