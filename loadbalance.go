package utils

type Backend struct {
	Name   string
	Weight int
}

// LoadBalancedByWeight selects a backend based on their weight.
func LoadBalancedByWeight(backends []Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	totalWeight := 0
	for _, b := range backends {
		totalWeight += b.Weight
	}
	r := int(RandomInt64(1, int64(totalWeight)))
	for _, b := range backends {
		if r <= b.Weight {
			return &b
		}
		r -= b.Weight
	}
	i := RandomInt64(0, int64(len(backends)-1))
	return &backends[i]
}
