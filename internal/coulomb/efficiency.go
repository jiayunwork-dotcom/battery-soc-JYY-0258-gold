package coulomb

// Efficiency models the coulombic efficiency of a battery: not all charge
// put in is recoverable. Charge efficiency is typically 99%+ for Li-ion
// but discharge efficiency may be lower due to internal losses.
type Efficiency struct {
	ChargeEff    float64 // fraction [0,1], default 0.99
	DischargeEff float64 // fraction [0,1], default 0.995
}

// DefaultEfficiency returns typical Li-ion coulombic efficiencies.
func DefaultEfficiency() Efficiency {
	return Efficiency{ChargeEff: 0.99, DischargeEff: 0.995}
}

// AdjustedCurrent applies coulombic efficiency to the raw current.
// Charge current is reduced (less goes in) and discharge current is
// increased (more comes out of SOC accounting).
func (e *Efficiency) AdjustedCurrent(current float64) float64 {
	if current > 0 {
		return current * e.ChargeEff
	}
	if current < 0 {
		return current / e.DischargeEff
	}
	return 0
}

// CoulombSOCWithEfficiency integrates one step with coulombic efficiency
// correction applied.
func CoulombSOCWithEfficiency(capacityAh, prevSOC, current, dtH float64, eff Efficiency) float64 {
	adjCurrent := eff.AdjustedCurrent(current)
	return CoulombSOC(capacityAh, prevSOC, adjCurrent, dtH)
}

// CoulombFromLogWithEfficiency integrates a series with efficiency correction.
func CoulombFromLogWithEfficiency(capacityAh float64, samples []CurrentSample, eff Efficiency) float64 {
	soc := 50.0
	for _, s := range samples {
		soc = CoulombSOCWithEfficiency(capacityAh, soc, s.Current, s.DT, eff)
	}
	return soc
}

// RoundTripEfficiency computes the overall round-trip efficiency from
// charge and discharge efficiencies.
func RoundTripEfficiency(eff Efficiency) float64 {
	return eff.ChargeEff * eff.DischargeEff
}

// EnergyLoss computes the energy lost in a charge-discharge cycle given
// the charge transferred and nominal voltage.
func EnergyLoss(chargeAh, nominalV float64, eff Efficiency) float64 {
	energyIn := chargeAh * nominalV
	energyOut := energyIn * RoundTripEfficiency(eff)
	return energyIn - energyOut
}
