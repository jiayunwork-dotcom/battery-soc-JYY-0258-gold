package cell

import "battery-soc/internal/coulomb"

// BalanceStrategy defines how cells in a pack are balanced.
type BalanceStrategy int

const (
	// NoBalance means no active balancing.
	NoBalance BalanceStrategy = iota
	// PassiveBalance bleeds excess charge from high-SOC cells.
	PassiveBalance
	// ActiveBalance transfers charge from high to low SOC cells.
	ActiveBalance
)

// BalanceConfig controls the balancing algorithm.
type BalanceConfig struct {
	Strategy     BalanceStrategy
	Threshold    float64 // SOC difference threshold to trigger balancing
	BalanceCurr  float64 // balancing current in Amps
	BalanceTime  float64 // duration of balance step in hours
}

// DefaultBalanceConfig returns sensible defaults for passive balancing.
func DefaultBalanceConfig() BalanceConfig {
	return BalanceConfig{
		Strategy:    PassiveBalance,
		Threshold:   2.0, // 2% SOC difference triggers balance
		BalanceCurr: 0.1, // 100mA balance current
		BalanceTime: 0.01, // 36 seconds per step
	}
}

// Balance applies one balancing step to a pack. For passive balancing,
// cells above the minimum SOC + threshold are discharged by BalanceCurr
// for BalanceTime. Returns the number of cells balanced.
func Balance(p *Pack, cfg BalanceConfig) int {
	if cfg.Strategy == NoBalance || len(p.Cells) < 2 {
		return 0
	}

	minSOC := p.MinSOC()
	count := 0

	switch cfg.Strategy {
	case PassiveBalance:
		for _, c := range p.Cells {
			if c.State().SOC > minSOC+cfg.Threshold {
				// Bleed current to lower SOC
				c.Step(-cfg.BalanceCurr, cfg.BalanceTime)
				count++
			}
		}
	case ActiveBalance:
		maxSOC := p.MaxSOC()
		if maxSOC-minSOC < cfg.Threshold {
			return 0
		}
		// Find highest and lowest cells
		var highest, lowest *Cell
		for _, c := range p.Cells {
			if c.State().SOC == maxSOC && highest == nil {
				highest = c
			}
			if c.State().SOC == minSOC && lowest == nil {
				lowest = c
			}
		}
		if highest != nil && lowest != nil {
			highest.Step(-cfg.BalanceCurr, cfg.BalanceTime)
			lowest.Step(cfg.BalanceCurr, cfg.BalanceTime)
			count = 2
		}
	}
	return count
}

// BalanceUntilConverged repeatedly applies balancing until the pack
// imbalance is below the threshold or maxIterations is reached.
func BalanceUntilConverged(p *Pack, cfg BalanceConfig, maxIterations int) int {
	total := 0
	for i := 0; i < maxIterations; i++ {
		if p.Imbalance() <= cfg.Threshold {
			break
		}
		n := Balance(p, cfg)
		if n == 0 {
			break
		}
		total += n
	}
	return total
}

// PackSOCAfterBalance runs a profile on a pack with periodic balancing
// and returns final pack SOC.
func PackSOCAfterBalance(p *Pack, samples []coulomb.CurrentSample, cfg BalanceConfig, balanceEvery int) float64 {
	for i, s := range samples {
		p.Step(s.Current, s.DT)
		if balanceEvery > 0 && (i+1)%balanceEvery == 0 {
			Balance(p, cfg)
		}
	}
	return p.MinSOC()
}
