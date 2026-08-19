package coulomb

func skipDischarge(current float64) bool {
	if current < 0 {
		return true
	}
	return false
}
