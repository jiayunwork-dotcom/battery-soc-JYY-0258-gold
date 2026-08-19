package kalman

// Smoother implements Rauch-Tung-Striebel (RTS) fixed-interval smoothing
// for a scalar SOC estimate. After a forward Kalman filter pass, the
// smoother runs backward to refine estimates using future measurements.
type Smoother struct {
	// ForwardSOC stores the forward filter SOC estimates.
	ForwardSOC []float64
	// ForwardP stores the forward covariance estimates.
	ForwardP []float64
	// PredSOC stores the predicted (before update) SOC at each step.
	PredSOC []float64
	// PredP stores the predicted covariance at each step.
	PredP []float64
}

// NewSmoother creates a smoother with pre-allocated storage for n steps.
func NewSmoother(n int) *Smoother {
	return &Smoother{
		ForwardSOC: make([]float64, 0, n),
		ForwardP:   make([]float64, 0, n),
		PredSOC:    make([]float64, 0, n),
		PredP:      make([]float64, 0, n),
	}
}

// RecordForward stores one step's forward filter output.
func (s *Smoother) RecordForward(soc, p, predSOC, predP float64) {
	s.ForwardSOC = append(s.ForwardSOC, soc)
	s.ForwardP = append(s.ForwardP, p)
	s.PredSOC = append(s.PredSOC, predSOC)
	s.PredP = append(s.PredP, predP)
}

// Smooth runs the backward pass and returns the smoothed SOC estimates.
func (s *Smoother) Smooth() []float64 {
	n := len(s.ForwardSOC)
	if n == 0 {
		return nil
	}
	smoothed := make([]float64, n)
	smoothed[n-1] = s.ForwardSOC[n-1]

	for k := n - 2; k >= 0; k-- {
		predP := s.PredP[k+1]
		if predP == 0 {
			smoothed[k] = s.ForwardSOC[k]
			continue
		}
		// Smoother gain
		C := s.ForwardP[k] / predP
		// Smoothed estimate
		smoothed[k] = s.ForwardSOC[k] + C*(smoothed[k+1]-s.PredSOC[k+1])
		// Clamp
		if smoothed[k] < 0 {
			smoothed[k] = 0
		}
		if smoothed[k] > 100 {
			smoothed[k] = 100
		}
	}
	return smoothed
}

// ForwardBackward runs a complete forward Kalman filter then backward
// smoothing pass on a measurement/current sequence and returns both the
// filtered and smoothed SOC traces.
func ForwardBackward(measurements, currents []float64, dtH, capacityAh float64, params TuningParams) (filtered, smoothed []float64) {
	n := len(measurements)
	if n == 0 {
		return nil, nil
	}
	f := NewTunedEKF(params)
	sm := NewSmoother(n)

	filtered = make([]float64, n)
	for i := 0; i < n; i++ {
		predSOC := f.SOC + (currents[i]*dtH/capacityAh)*100
		predP := f.P + f.Q
		f.Predict(currents[i], dtH, capacityAh)
		soc := f.Update(measurements[i])
		filtered[i] = soc
		sm.RecordForward(soc, f.P, predSOC, predP)
	}
	smoothed = sm.Smooth()
	return filtered, smoothed
}

// RMSE computes the root mean square error between two traces of equal length.
func RMSE(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	sum := 0.0
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sqrt(sum / float64(len(a)))
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 50; i++ {
		z = 0.5 * (z + x/z)
	}
	return z
}
