package device

import "math/rand"

type Device struct {
	WorkingTime          int
	ReplacementCost      int
	WasFailed            bool
	ProbabilityOfFailure func(time int) float64
}

// Should be called every tick to determine if the device has failed

func (d *Device) Failed() bool {
	result := rand.Float64() < d.ProbabilityOfFailure(d.WorkingTime)
	if result {
		d.WasFailed = true
	}
	return result
}
