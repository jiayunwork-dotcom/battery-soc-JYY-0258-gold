package report

import (
	"fmt"
	"io"

	"battery-soc/internal/coulomb"
)

// WriteSOCTraceCSV writes a SOC trace as CSV with columns time_h,soc_pct.
func WriteSOCTraceCSV(w io.Writer, trace []float64, samples []coulomb.CurrentSample) error {
	if _, err := io.WriteString(w, "time_h,soc_pct\n"); err != nil {
		return err
	}
	t := 0.0
	for i, soc := range trace {
		if i < len(samples) {
			t += samples[i].DT
		}
		line := fmt.Sprintf("%.4f,%.4f\n", t, soc)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

// WriteDetailedCSV writes a detailed per-step CSV with time, SOC,
// current, voltage, and temperature columns.
func WriteDetailedCSV(w io.Writer, times, socs, currents, voltages, temps []float64) error {
	header := "time_h,soc_pct,current_a,voltage_v,temperature_c\n"
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	n := len(times)
	for i := 0; i < n; i++ {
		var soc, cur, vol, temp float64
		if i < len(socs) {
			soc = socs[i]
		}
		if i < len(currents) {
			cur = currents[i]
		}
		if i < len(voltages) {
			vol = voltages[i]
		}
		if i < len(temps) {
			temp = temps[i]
		}
		line := fmt.Sprintf("%.4f,%.4f,%.4f,%.4f,%.2f\n", times[i], soc, cur, vol, temp)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

// WriteComparisonCSV writes a comparison of two SOC estimation methods.
func WriteComparisonCSV(w io.Writer, times, socA, socB []float64, labelA, labelB string) error {
	header := fmt.Sprintf("time_h,%s,%s,diff\n", labelA, labelB)
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	n := len(times)
	for i := 0; i < n; i++ {
		a, b := 0.0, 0.0
		if i < len(socA) {
			a = socA[i]
		}
		if i < len(socB) {
			b = socB[i]
		}
		diff := a - b
		line := fmt.Sprintf("%.4f,%.4f,%.4f,%.4f\n", times[i], a, b, diff)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

// WriteCycleReport writes cycle-by-cycle degradation data as CSV.
func WriteCycleReport(w io.Writer, cycles []int, capacities, sohs []float64) error {
	if _, err := io.WriteString(w, "cycle,capacity_ah,soh\n"); err != nil {
		return err
	}
	n := len(cycles)
	for i := 0; i < n; i++ {
		var cap, soh float64
		if i < len(capacities) {
			cap = capacities[i]
		}
		if i < len(sohs) {
			soh = sohs[i]
		}
		line := fmt.Sprintf("%d,%.4f,%.6f\n", cycles[i], cap, soh)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}
