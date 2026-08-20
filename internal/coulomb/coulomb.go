// Package coulomb implements classical coulomb-counting SOC estimation.
// Current is in Amperes (+ charge, - discharge) and dt is in hours.
package coulomb

// CurrentSample is one integration step: DT hours at Current Amps.
type CurrentSample struct {
	DT      float64
	Current float64
}

func clampSOC(s float64) float64 {
	if s < 0 {
		return 0
	}
	if s > 100 {
		return 100
	}
	return s
}

// CoulombSOC integrates one step:
//
//	soc = prevSOC + current*dtH/capacityAh*100
//
// The result is clamped to [0,100].
func CoulombSOC(capacityAh, prevSOC, current, dtH float64) float64 {
	soc := prevSOC + (current*dtH/capacityAh)*100
	return clampSOC(soc)
}

// CoulombFromLog integrates a whole series of current samples starting at SOC=50,
// clamping each step, and returns the final SOC.
func CoulombFromLog(capacityAh float64, samples []CurrentSample) float64 {
	soc := 50.0
	for _, s := range samples {
		soc = CoulombSOC(capacityAh, soc, s.Current, s.DT)
	}
	return soc
}
