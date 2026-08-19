package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"battery-soc/internal/coulomb"
)

// usageError signals a usage problem (exit code 2) as opposed to a runtime
// failure (exit code 1).
type usageError struct{ msg string }

func (e *usageError) Error() string { return e.msg }

func main() {
	err := run()
	if err != nil {
		if _, ok := err.(*usageError); ok {
			fmt.Fprintln(os.Stderr, "usage error:", err)
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// run parses flags, loads a current log, and prints the final coulomb-counted
// SOC. It never panics; all failures are returned as errors.
func run() error {
	inputPath := flag.String("input", "", "path to current-log CSV (default: read stdin)")
	capacity := flag.Float64("capacity", 50.0, "battery capacity in Ah")
	flag.Parse()

	if *capacity <= 0 {
		return &usageError{"capacity must be positive"}
	}

	data, err := loadInput(*inputPath)
	if err != nil {
		return err
	}

	samples, err := parseCurrentLog(data)
	if err != nil {
		return fmt.Errorf("parse current log: %w", err)
	}

	finalSOC := coulomb.CoulombFromLog(*capacity, samples)
	fmt.Printf("final_soc: %.4f\n", finalSOC)
	return nil
}

// loadInput reads from the given path, or from stdin when path is empty.
// A missing file is a runtime error; an absent piped stdin is a usage error.
func loadInput(path string) ([]byte, error) {
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("cannot read input file: %w", err)
		}
		return b, nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("cannot stat stdin: %w", err)
	}
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		return nil, &usageError{"no input provided: pipe data via stdin or use -input <path>"}
	}
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("cannot read stdin: %w", err)
	}
	return b, nil
}

// parseCurrentLog decodes a CSV with columns dt,current (dt in hours, current
// in Amps; a textual header row is skipped). Non-numeric, non-header rows fail.
func parseCurrentLog(data []byte) ([]coulomb.CurrentSample, error) {
	r := csv.NewReader(strings.NewReader(string(data)))
	r.FieldsPerRecord = 2
	recs, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv: %w", err)
	}

	var samples []coulomb.CurrentSample
	for i, rec := range recs {
		if len(rec) < 2 {
			return nil, fmt.Errorf("line %d: need 2 fields (dt,current)", i+1)
		}
		dt, err1 := strconv.ParseFloat(strings.TrimSpace(rec[0]), 64)
		cur, err2 := strconv.ParseFloat(strings.TrimSpace(rec[1]), 64)
		if err1 != nil || err2 != nil {
			if i == 0 {
				continue // tolerate a header row
			}
			return nil, fmt.Errorf("line %d: invalid number", i+1)
		}
		samples = append(samples, coulomb.CurrentSample{DT: dt, Current: cur})
	}
	if len(samples) == 0 {
		return nil, fmt.Errorf("no valid samples found")
	}
	return samples, nil
}
