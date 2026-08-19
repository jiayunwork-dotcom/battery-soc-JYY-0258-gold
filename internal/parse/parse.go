// Package parse provides multi-format reading of battery current/voltage
// log data. It supports CSV, TSV, and JSON array formats, decoupling I/O
// from the estimation algorithms.
package parse

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"battery-soc/internal/coulomb"
)

// ReadCSV reads current samples from a CSV with columns dt,current.
// The first row is skipped if it cannot be parsed as numbers (header).
func ReadCSV(r io.Reader) ([]coulomb.CurrentSample, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	var samples []coulomb.CurrentSample
	lineNo := 0
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse: line %d: %w", lineNo+1, err)
		}
		lineNo++
		if len(row) < 2 {
			continue
		}
		dt, err1 := strconv.ParseFloat(strings.TrimSpace(row[0]), 64)
		cur, err2 := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		if err1 != nil || err2 != nil {
			if lineNo == 1 {
				continue // skip header
			}
			return nil, fmt.Errorf("parse: line %d: invalid number", lineNo)
		}
		samples = append(samples, coulomb.CurrentSample{DT: dt, Current: cur})
	}
	if len(samples) == 0 {
		return nil, fmt.Errorf("parse: no valid samples found")
	}
	return samples, nil
}

// VoltageCurrentSample extends CurrentSample with voltage and temperature.
type VoltageCurrentSample struct {
	DT          float64 `json:"dt"`
	Current     float64 `json:"current"`
	Voltage     float64 `json:"voltage"`
	Temperature float64 `json:"temperature"`
}

// ReadVoltageCSV reads extended samples with columns dt,current,voltage,temp.
func ReadVoltageCSV(r io.Reader) ([]VoltageCurrentSample, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	var samples []VoltageCurrentSample
	lineNo := 0
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse: line %d: %w", lineNo+1, err)
		}
		lineNo++
		if len(row) < 4 {
			continue
		}
		dt, e1 := strconv.ParseFloat(strings.TrimSpace(row[0]), 64)
		cur, e2 := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		vol, e3 := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		temp, e4 := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			if lineNo == 1 {
				continue
			}
			return nil, fmt.Errorf("parse: line %d: invalid number", lineNo)
		}
		samples = append(samples, VoltageCurrentSample{
			DT: dt, Current: cur, Voltage: vol, Temperature: temp,
		})
	}
	if len(samples) == 0 {
		return nil, fmt.Errorf("parse: no valid samples found")
	}
	return samples, nil
}

// ReadJSON reads current samples from a JSON array:
// [{"dt":0.5,"current":10},{"dt":1.0,"current":-5}]
func ReadJSON(r io.Reader) ([]coulomb.CurrentSample, error) {
	var raw []struct {
		DT      float64 `json:"dt"`
		Current float64 `json:"current"`
	}
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, fmt.Errorf("parse: json: %w", err)
	}
	samples := make([]coulomb.CurrentSample, len(raw))
	for i, r := range raw {
		samples[i] = coulomb.CurrentSample{DT: r.DT, Current: r.Current}
	}
	if len(samples) == 0 {
		return nil, fmt.Errorf("parse: no samples in JSON")
	}
	return samples, nil
}

// WriteCSV writes samples back to CSV format.
func WriteCSV(w io.Writer, samples []coulomb.CurrentSample) error {
	if _, err := io.WriteString(w, "dt,current\n"); err != nil {
		return err
	}
	for _, s := range samples {
		line := fmt.Sprintf("%.4f,%.4f\n", s.DT, s.Current)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}
